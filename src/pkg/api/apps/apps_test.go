// Copyright 2025-2026 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

package apps

import (
	"context"
	"net/http"
	"testing"

	"github.com/defenseunicorns/uds-portal/src/pkg/api/auth/incluster"
)

func TestDisplayNameForApp(t *testing.T) {
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
			got := displayNameForApp(tt.packageName)
			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestEndpointURL(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		gateway string
		domain  string
		want    string
	}{
		{"tenant gateway", "podinfo", "tenant", "uds.dev", "podinfo.uds.dev"},
		{"admin gateway", "grafana", "admin", "uds.dev", "grafana.admin.uds.dev"},
		{"empty gateway treated as tenant", "app", "", "uds.dev", "app.uds.dev"},
		{"custom gateway", "app", "custom", "uds.dev", "app.custom.uds.dev"},
		{"empty domain returns host only", "app", "tenant", "", "app"},
		{"empty domain with admin returns host only", "app", "admin", "", "app"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := endpointURL(tt.host, tt.gateway, tt.domain)
			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestToAPIApps(t *testing.T) {
	const domain = "uds.dev"
	tests := []struct {
		name         string
		pkgs         []Package
		accountURL   string
		wantLen      int
		wantFirstURL string // empty means "no My Account entry expected"
	}{
		{
			name:         "returns My Account when no packages",
			pkgs:         nil,
			accountURL:   "sso.uds.dev",
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
			wantLen:      1,
			wantFirstURL: "",
		},
		{
			name:         "omits My Account when URL is empty and no packages",
			pkgs:         nil,
			accountURL:   "",
			wantLen:      0,
			wantFirstURL: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toAPIApps(nil, tt.pkgs, tt.accountURL, domain)
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
				return
			}
			if got[0].URL != tt.wantFirstURL {
				t.Fatalf("expected first URL %q, got %q", tt.wantFirstURL, got[0].URL)
			}
		})
	}
}

func TestFilterByUserGroup(t *testing.T) {
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
	const domain = "uds.dev"
	tests := []struct {
		name     string
		pkgs     []Package
		wantApps []APIApp
	}{
		{
			name: "admin gateway tagged",
			pkgs: []Package{{
				Metadata: Metadata{Name: "grafana"},
				Spec: Spec{Network: Network{Expose: []Expose{
					{Host: "grafana", Gateway: "admin"},
				}}},
			}},
			wantApps: []APIApp{{Name: "Grafana", URL: "grafana.admin.uds.dev", Gateway: "admin"}},
		},
		{
			name: "tenant gateway tagged",
			pkgs: []Package{{
				Metadata: Metadata{Name: "podinfo"},
				Spec: Spec{Network: Network{Expose: []Expose{
					{Host: "podinfo", Gateway: "tenant"},
				}}},
			}},
			wantApps: []APIApp{{Name: "Podinfo", URL: "podinfo.uds.dev", Gateway: "tenant"}},
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
			wantApps: []APIApp{
				{Name: "Mixed", URL: "front.uds.dev", Gateway: "tenant"},
				{Name: "Mixed", URL: "back.admin.uds.dev", Gateway: "admin"},
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
			wantApps: []APIApp{{Name: "Custom App", URL: "app.custom.uds.dev", Gateway: "custom"}},
		},
		{
			name: "empty gateway treated as tenant",
			pkgs: []Package{{
				Metadata: Metadata{Name: "legacy"},
				Spec: Spec{Network: Network{Expose: []Expose{
					{Host: "legacy", Gateway: ""},
				}}},
			}},
			wantApps: []APIApp{{Name: "Legacy", URL: "legacy.uds.dev", Gateway: ""}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toAPIApps(nil, tt.pkgs, "", domain)
			if len(got) != len(tt.wantApps) {
				t.Fatalf("expected %d apps, got %d: %+v", len(tt.wantApps), len(got), got)
			}
			for i, want := range tt.wantApps {
				if got[i].Name != want.Name || got[i].URL != want.URL || got[i].Gateway != want.Gateway {
					t.Errorf("app %d: expected %+v, got %+v", i, want, got[i])
				}
			}
		})
	}
}

func TestToAPIApps_MyAccountHasNoGateway(t *testing.T) {
	got := toAPIApps(nil, nil, "sso.uds.dev", "uds.dev")
	if len(got) != 1 {
		t.Fatalf("expected 1 app, got %d", len(got))
	}
	if got[0].Gateway != "" {
		t.Fatalf("expected empty gateway for My Account, got %q", got[0].Gateway)
	}
}
