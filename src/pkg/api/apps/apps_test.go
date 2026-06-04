// Copyright 2025-2026 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

package apps

import (
	"context"
	"net/http"
	"testing"

	"github.com/defenseunicorns/uds-portal/src/pkg/api/auth/incluster"
)

func TestFormatPackageName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		packageName string
		want        string
	}{
		{
			name:        "formats dashed name and title case",
			packageName: "mission-control",
			want:        "Mission Control",
		},
		{
			name:        "normalizes uds acronym",
			packageName: "uds-registry",
			want:        "UDS Registry",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := formatPackageName(tt.packageName)
			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestDisplayNameForApp(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		title       string
		packageName string
		want        string
	}{
		{
			name:        "returns annotation title when present",
			title:       "My Custom Title",
			packageName: "my-package",
			want:        "My Custom Title",
		},
		{
			name:        "falls back to formatPackageName when title empty",
			title:       "",
			packageName: "uds-registry",
			want:        "UDS Registry",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := displayNameForApp(tt.title, tt.packageName)
			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestEndpointURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		host         string
		gateway      string
		tenantDomain string
		adminDomain  string
		want         string
	}{
		// tenant / empty gateway
		{"tenant gateway + tenant domain", "podinfo", "tenant", "uds.dev", "", "podinfo.uds.dev"},
		{"empty gateway treated as tenant", "app", "", "uds.dev", "", "app.uds.dev"},
		{"tenant + empty tenant domain returns empty", "app", "tenant", "", "", ""},
		{"empty gateway + empty tenant domain returns empty", "app", "", "", "", ""},

		// admin gateway
		{
			"admin + explicit adminDomain",
			"grafana",
			"admin",
			"uds.dev",
			"admin.example.com",
			"grafana.admin.example.com",
		},
		{
			"admin + empty adminDomain falls back to admin.tenantDomain",
			"grafana",
			"admin",
			"uds.dev",
			"",
			"grafana.admin.uds.dev",
		},
		{"admin + both domains empty returns empty", "grafana", "admin", "", "", ""},

		// custom gateway
		{"custom gateway + tenant domain", "app", "passthrough", "uds.dev", "", "app.passthrough.uds.dev"},
		{"custom gateway + empty tenant domain returns empty", "app", "passthrough", "", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := endpointURL(tt.host, tt.gateway, tt.tenantDomain, tt.adminDomain)
			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestToAPIApps(t *testing.T) {
	t.Parallel()
	const tenantDomain = "uds.dev"
	const adminDomain = ""
	tests := []struct {
		name         string
		pkgs         []Package
		accountURL   string
		tenantDomain string
		adminDomain  string
		wantLen      int
		wantFirstURL string // empty means "no My Account entry expected"
	}{
		{
			name:         "returns My Account when no packages",
			pkgs:         nil,
			accountURL:   "sso.uds.dev",
			tenantDomain: tenantDomain,
			adminDomain:  adminDomain,
			wantLen:      1,
			wantFirstURL: "sso.uds.dev",
		},
		{
			name: "prepends My Account when packages exist",
			pkgs: []Package{
				{
					Metadata: Metadata{Name: "grafana"},
					Spec: Spec{Network: Network{Expose: []Expose{
						{Host: "grafana", Gateway: "tenant"},
					}}},
				},
			},
			accountURL:   "sso.uds.dev",
			tenantDomain: tenantDomain,
			adminDomain:  adminDomain,
			wantLen:      2,
			wantFirstURL: "sso.uds.dev",
		},
		{
			name: "uses custom My Account URL",
			pkgs: []Package{
				{
					Metadata: Metadata{Name: "grafana"},
					Spec: Spec{Network: Network{Expose: []Expose{
						{Host: "grafana", Gateway: "tenant"},
					}}},
				},
			},
			accountURL:   "sso.example.com",
			tenantDomain: tenantDomain,
			adminDomain:  adminDomain,
			wantLen:      2,
			wantFirstURL: "sso.example.com",
		},
		{
			name: "omits My Account when URL is empty",
			pkgs: []Package{
				{
					Metadata: Metadata{Name: "grafana"},
					Spec: Spec{Network: Network{Expose: []Expose{
						{Host: "grafana", Gateway: "tenant"},
					}}},
				},
			},
			accountURL:   "",
			tenantDomain: tenantDomain,
			adminDomain:  adminDomain,
			wantLen:      1,
			wantFirstURL: "",
		},
		{
			name:         "omits My Account when URL is empty and no packages",
			pkgs:         nil,
			accountURL:   "",
			tenantDomain: tenantDomain,
			adminDomain:  adminDomain,
			wantLen:      0,
			wantFirstURL: "",
		},
		{
			name: "skips tile when tenant domain empty",
			pkgs: []Package{
				{
					Metadata: Metadata{Name: "grafana"},
					Spec: Spec{Network: Network{Expose: []Expose{
						{Host: "grafana", Gateway: "tenant"},
					}}},
				},
			},
			accountURL:   "",
			tenantDomain: "",
			adminDomain:  "",
			wantLen:      0,
			wantFirstURL: "",
		},
		{
			name: "uses explicit adminDomain for admin gateway tile",
			pkgs: []Package{
				{
					Metadata: Metadata{Name: "grafana"},
					Spec: Spec{Network: Network{Expose: []Expose{
						{Host: "grafana", Gateway: "admin"},
					}}},
				},
			},
			accountURL:   "",
			tenantDomain: "uds.dev",
			adminDomain:  "admin.example.com",
			wantLen:      1,
			wantFirstURL: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := toAPIApps(nil, tt.pkgs, tt.accountURL, tt.tenantDomain, tt.adminDomain)
			if got == nil {
				t.Fatal("expected non-nil slice, got nil")
			}
			if len(got) != tt.wantLen {
				t.Fatalf("expected %d apps, got %d", tt.wantLen, len(got))
			}
			if tt.wantFirstURL == "" {
				for _, app := range got {
					if app.Name == myAccountName {
						t.Fatalf("did not expect My Account entry, got %+v", app)
					}
				}
				// For the "uses explicit adminDomain" case, verify the URL
				if tt.name == "uses explicit adminDomain for admin gateway tile" {
					if len(got) > 0 && got[0].URL != "grafana.admin.example.com" {
						t.Fatalf("expected URL %q, got %q", "grafana.admin.example.com", got[0].URL)
					}
				}
				return
			}
			if got[0].URL != tt.wantFirstURL {
				t.Fatalf("expected first URL %q, got %q", tt.wantFirstURL, got[0].URL)
			}
		})
	}
}

func TestFilterByUserGroup(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		userGroups []string
		pkgs       []Package
		wantLen    int
	}{
		{
			name:       "requires SSO section",
			userGroups: []string{"/UDS Core/Admin"},
			pkgs: []Package{
				{
					Metadata: Metadata{Name: "no-sso"},
					Spec: Spec{
						Network: Network{Expose: []Expose{{Host: "no-sso", Gateway: "tenant"}}},
						SSO:     nil,
					},
				},
			},
			wantLen: 0,
		},
		{
			name:       "allows when groups not specified",
			userGroups: []string{"/UDS Core/Auditor"},
			pkgs: []Package{
				{
					Metadata: Metadata{Name: "sso-without-groups"},
					Spec: Spec{
						Network: Network{Expose: []Expose{{Host: "auditor", Gateway: "tenant"}}},
						SSO:     []SSO{{Groups: Groups{AnyOf: nil}}},
					},
				},
			},
			wantLen: 1,
		},
		{
			name:       "allows matching group",
			userGroups: []string{"/UDS Core/Auditor", "/UDS Core/Admin"},
			pkgs: []Package{
				{
					Metadata: Metadata{Name: "group-matched"},
					Spec: Spec{
						Network: Network{Expose: []Expose{{Host: "group-matched", Gateway: "tenant"}}},
						SSO:     []SSO{{Groups: Groups{AnyOf: []string{"/UDS Core/Developer", "/UDS Core/Admin"}}}},
					},
				},
			},
			wantLen: 1,
		},
		{
			name:       "disallows non-matching group",
			userGroups: []string{"/UDS Core/Auditor"},
			pkgs: []Package{
				{
					Metadata: Metadata{Name: "group-mismatched"},
					Spec: Spec{
						Network: Network{Expose: []Expose{{Host: "group-mismatched", Gateway: "tenant"}}},
						SSO:     []SSO{{Groups: Groups{AnyOf: []string{"/UDS Core/Developer", "/UDS Core/Admin"}}}},
					},
				},
			},
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatalf("failed creating request: %v", err)
			}
			req = req.WithContext(context.WithValue(req.Context(), incluster.GroupKey, tt.userGroups))

			got := filterByUserGroup(req, tt.pkgs, true)
			if len(got) != tt.wantLen {
				t.Fatalf("expected %d packages, got %d", tt.wantLen, len(got))
			}
		})
	}
}

func TestToAPIApps_GatewayTagging(t *testing.T) {
	t.Parallel()
	const tenantDomain = "uds.dev"
	tests := []struct {
		name        string
		pkgs        []Package
		adminDomain string
		wantApps    []APIApp
	}{
		{
			name: "admin gateway tagged (fallback to admin.tenantDomain)",
			pkgs: []Package{{
				Metadata: Metadata{Name: "grafana"},
				Spec: Spec{Network: Network{Expose: []Expose{
					{Host: "grafana", Gateway: "admin"},
				}}},
			}},
			adminDomain: "",
			wantApps: []APIApp{
				{Name: "Grafana", URL: "grafana.admin.uds.dev", Gateway: "admin", Group: groupUDSCore},
			},
		},
		{
			name: "admin gateway tagged with explicit adminDomain",
			pkgs: []Package{{
				Metadata: Metadata{Name: "grafana"},
				Spec: Spec{Network: Network{Expose: []Expose{
					{Host: "grafana", Gateway: "admin"},
				}}},
			}},
			adminDomain: "admin.example.com",
			wantApps: []APIApp{
				{Name: "Grafana", URL: "grafana.admin.example.com", Gateway: "admin", Group: groupUDSCore},
			},
		},
		{
			name: "tenant gateway tagged",
			pkgs: []Package{{
				Metadata: Metadata{Name: "podinfo"},
				Spec: Spec{Network: Network{Expose: []Expose{
					{Host: "podinfo", Gateway: "tenant"},
				}}},
			}},
			adminDomain: "",
			wantApps:    []APIApp{{Name: "Podinfo", URL: "podinfo.uds.dev", Gateway: "tenant", Group: groupOther}},
		},
		{
			name: "mixed expose yields one tile per expose entry",
			pkgs: []Package{{
				Metadata: Metadata{Name: "mixed"},
				Spec: Spec{Network: Network{Expose: []Expose{
					{Host: "front", Gateway: "tenant"},
					{Host: "back", Gateway: "admin"},
				}}},
			}},
			adminDomain: "",
			wantApps: []APIApp{
				{Name: "Mixed", URL: "back.admin.uds.dev", Gateway: "admin", Group: groupOther},
				{Name: "Mixed", URL: "front.uds.dev", Gateway: "tenant", Group: groupOther},
			},
		},
		{
			name: "custom gateway",
			pkgs: []Package{{
				Metadata: Metadata{Name: "custom-app"},
				Spec: Spec{Network: Network{Expose: []Expose{
					{Host: "app", Gateway: "custom"},
				}}},
			}},
			adminDomain: "",
			wantApps: []APIApp{
				{Name: "Custom App", URL: "app.custom.uds.dev", Gateway: "custom", Group: groupOther},
			},
		},
		{
			name: "empty gateway treated as tenant",
			pkgs: []Package{{
				Metadata: Metadata{Name: "legacy"},
				Spec: Spec{Network: Network{Expose: []Expose{
					{Host: "legacy", Gateway: ""},
				}}},
			}},
			adminDomain: "",
			wantApps:    []APIApp{{Name: "Legacy", URL: "legacy.uds.dev", Gateway: "", Group: groupOther}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := toAPIApps(nil, tt.pkgs, "", tenantDomain, tt.adminDomain)
			if len(got) != len(tt.wantApps) {
				t.Fatalf("expected %d apps, got %d: %+v", len(tt.wantApps), len(got), got)
			}
			for i, want := range tt.wantApps {
				if got[i].Name != want.Name || got[i].URL != want.URL || got[i].Gateway != want.Gateway ||
					got[i].Group != want.Group {
					t.Errorf("app %d: expected %+v, got %+v", i, want, got[i])
				}
			}
		})
	}
}

func TestToAPIApps_MyAccountHasNoGateway(t *testing.T) {
	t.Parallel()
	got := toAPIApps(nil, nil, "sso.uds.dev", "uds.dev", "")
	if len(got) != 1 {
		t.Fatalf("expected 1 app, got %d", len(got))
	}
	if got[0].Gateway != "" {
		t.Fatalf("expected empty gateway for My Account, got %q", got[0].Gateway)
	}
}

func TestFilterHiddenEndpoints(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		pkgs      []Package
		wantPkgs  int
		wantHosts []string
	}{
		{
			name: "annotation-listed host excluded",
			pkgs: []Package{{
				Metadata: Metadata{
					Name:        "gitlab",
					Annotations: map[string]string{"uds.dev/portal-hide-apps": "gitlab-registry"},
				},
				Spec: Spec{Network: Network{Expose: []Expose{
					{Host: "gitlab", Gateway: "tenant"},
					{Host: "gitlab-registry", Gateway: "tenant"},
				}}},
			}},
			wantPkgs:  1,
			wantHosts: []string{"gitlab"},
		},
		{
			name: "wildcard host auto-excluded",
			pkgs: []Package{{
				Metadata: Metadata{Name: "gitlab"},
				Spec: Spec{Network: Network{Expose: []Expose{
					{Host: "gitlab", Gateway: "tenant"},
					{Host: "*.pages", Gateway: "tenant"},
				}}},
			}},
			wantPkgs:  1,
			wantHosts: []string{"gitlab"},
		},
		{
			name: "non-listed FQDN host passes through",
			pkgs: []Package{{
				Metadata: Metadata{
					Name:        "gitlab",
					Annotations: map[string]string{"uds.dev/portal-hide-apps": "gitlab-registry"},
				},
				Spec: Spec{Network: Network{Expose: []Expose{
					{Host: "gitlab", Gateway: "tenant"},
				}}},
			}},
			wantPkgs:  1,
			wantHosts: []string{"gitlab"},
		},
		{
			name: "package dropped when all endpoints hidden",
			pkgs: []Package{{
				Metadata: Metadata{
					Name:        "all-hidden",
					Annotations: map[string]string{"uds.dev/portal-hide-apps": "registry"},
				},
				Spec: Spec{Network: Network{Expose: []Expose{
					{Host: "registry", Gateway: "tenant"},
				}}},
			}},
			wantPkgs:  0,
			wantHosts: nil,
		},
		{
			name: "annotation with extra whitespace trimmed",
			pkgs: []Package{{
				Metadata: Metadata{
					Name:        "gitlab",
					Annotations: map[string]string{"uds.dev/portal-hide-apps": " gitlab-registry , pages "},
				},
				Spec: Spec{Network: Network{Expose: []Expose{
					{Host: "gitlab", Gateway: "tenant"},
					{Host: "gitlab-registry", Gateway: "tenant"},
					{Host: "pages", Gateway: "tenant"},
				}}},
			}},
			wantPkgs:  1,
			wantHosts: []string{"gitlab"},
		},
		{
			name: "no annotation passes all through",
			pkgs: []Package{{
				Metadata: Metadata{Name: "podinfo"},
				Spec: Spec{Network: Network{Expose: []Expose{
					{Host: "podinfo", Gateway: "tenant"},
				}}},
			}},
			wantPkgs:  1,
			wantHosts: []string{"podinfo"},
		},
		{
			name: "mixed packages: one filtered one retained",
			pkgs: []Package{
				{
					Metadata: Metadata{
						Name:        "all-hidden",
						Annotations: map[string]string{"uds.dev/portal-hide-apps": "registry"},
					},
					Spec: Spec{Network: Network{Expose: []Expose{
						{Host: "registry", Gateway: "tenant"},
					}}},
				},
				{
					Metadata: Metadata{Name: "podinfo"},
					Spec: Spec{Network: Network{Expose: []Expose{
						{Host: "podinfo", Gateway: "tenant"},
					}}},
				},
			},
			wantPkgs:  1,
			wantHosts: []string{"podinfo"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := filterHiddenEndpoints(tt.pkgs)
			if len(got) != tt.wantPkgs {
				t.Fatalf("expected %d packages, got %d", tt.wantPkgs, len(got))
			}
			var gotHosts []string
			for _, pkg := range got {
				for _, e := range pkg.Spec.Network.Expose {
					gotHosts = append(gotHosts, e.Host)
				}
			}
			if len(gotHosts) != len(tt.wantHosts) {
				t.Fatalf(
					"expected %d host(s) %v, got %d host(s) %v",
					len(tt.wantHosts),
					tt.wantHosts,
					len(gotHosts),
					gotHosts,
				)
			}
			for i, h := range tt.wantHosts {
				if gotHosts[i] != h {
					t.Errorf("host[%d]: expected %q, got %q", i, h, gotHosts[i])
				}
			}
		})
	}
}

func TestGroupForPackage(t *testing.T) {
	t.Parallel()
	tests := []struct {
		pkgName   string
		wantGroup int
	}{
		{"uds-registry", groupUDSPlatform},
		{"uds-ui", groupUDSPlatform},
		{"cat", groupUDSPlatform},
		{"grafana", groupUDSCore},
		{"podinfo", groupOther},
		{"mission-app", groupOther},
	}
	for _, tt := range tests {
		t.Run(tt.pkgName, func(t *testing.T) {
			t.Parallel()
			got := groupForPackage(Package{Metadata: Metadata{Name: tt.pkgName}})
			if got != tt.wantGroup {
				t.Fatalf("expected group %d, got %d", tt.wantGroup, got)
			}
		})
	}
}

func TestToAPIApps_GroupAndSort(t *testing.T) {
	t.Parallel()
	pkgs := []Package{
		{
			Metadata: Metadata{Name: "podinfo"},
			Spec:     Spec{Network: Network{Expose: []Expose{{Host: "podinfo", Gateway: "tenant"}}}},
		},
		{
			Metadata: Metadata{Name: "grafana"},
			Spec:     Spec{Network: Network{Expose: []Expose{{Host: "grafana", Gateway: "admin"}}}},
		},
		{
			Metadata: Metadata{Name: "uds-registry"},
			Spec:     Spec{Network: Network{Expose: []Expose{{Host: "registry", Gateway: "tenant"}}}},
		},
		{
			Metadata: Metadata{Name: "another-app"},
			Spec:     Spec{Network: Network{Expose: []Expose{{Host: "another-app", Gateway: "tenant"}}}},
		},
		{
			Metadata: Metadata{Name: "cat"},
			Spec:     Spec{Network: Network{Expose: []Expose{{Host: "cat", Gateway: "tenant"}}}},
		},
	}

	got := toAPIApps(nil, pkgs, "sso.uds.dev", "uds.dev", "admin.uds.dev")

	// expected order: My Account(0), Cat(10), UDS Registry(10), Grafana(20), Another App(30), Podinfo(30)
	want := []struct {
		name  string
		group int
	}{
		{myAccountName, groupMyAccount},
		{"Cat", groupUDSPlatform},
		{"UDS Registry", groupUDSPlatform},
		{"Grafana", groupUDSCore},
		{"Another App", groupOther},
		{"Podinfo", groupOther},
	}

	if len(got) != len(want) {
		t.Fatalf("expected %d apps, got %d", len(want), len(got))
	}
	for i, w := range want {
		if got[i].Name != w.name {
			t.Errorf("app[%d]: expected name %q, got %q", i, w.name, got[i].Name)
		}
		if got[i].Group != w.group {
			t.Errorf("app[%d] %q: expected group %d, got %d", i, got[i].Name, w.group, got[i].Group)
		}
	}
}
