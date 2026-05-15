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

func TestToAPIApps(t *testing.T) {
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
					Status:   Status{Endpoints: []string{"grafana.uds.dev"}},
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
					Status:   Status{Endpoints: []string{"grafana.example.com"}},
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
					Status:   Status{Endpoints: []string{"grafana.uds.dev"}},
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
			got := toAPIApps(nil, tt.pkgs, tt.accountURL)
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
					Status:   Status{Endpoints: []string{"no-sso.uds.dev"}},
					Spec:     Spec{SSO: nil},
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
					Status:   Status{Endpoints: []string{"auditor.uds.dev"}},
					Spec: Spec{SSO: []SSO{
						{Groups: Groups{AnyOf: nil}},
					}},
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
					Status:   Status{Endpoints: []string{"group-matched.uds.dev"}},
					Spec: Spec{SSO: []SSO{
						{Groups: Groups{AnyOf: []string{"/UDS Core/Developer", "/UDS Core/Admin"}}},
					}},
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
					Status:   Status{Endpoints: []string{"group-mismatched.uds.dev"}},
					Spec: Spec{SSO: []SSO{
						{Groups: Groups{AnyOf: []string{"/UDS Core/Developer", "/UDS Core/Admin"}}},
					}},
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

func TestGatewayForEndpoint(t *testing.T) {
	pkg := Package{
		Spec: Spec{Network: Network{Expose: []Expose{
			{Host: "grafana", Gateway: "admin"},
			{Host: "podinfo", Gateway: "tenant"},
		}}},
	}
	tests := []struct {
		name     string
		endpoint string
		want     string
	}{
		{"exact host match", "grafana", "admin"},
		{"prefix host match with domain", "grafana.uds.dev", "admin"},
		{"tenant match", "podinfo.uds.dev", "tenant"},
		{"unknown endpoint", "unknown.uds.dev", ""},
		{"empty endpoint", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := gatewayForEndpoint(pkg, tt.endpoint)
			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestGatewayForEndpoint_NoExpose(t *testing.T) {
	pkg := Package{}
	if got := gatewayForEndpoint(pkg, "anything.uds.dev"); got != "" {
		t.Fatalf("expected empty, got %q", got)
	}
}

func TestGatewayForEndpoint_EmptyHostSkipped(t *testing.T) {
	pkg := Package{
		Spec: Spec{Network: Network{Expose: []Expose{
			{Host: "", Gateway: "admin"},
			{Host: "real", Gateway: "tenant"},
		}}},
	}
	if got := gatewayForEndpoint(pkg, "real.uds.dev"); got != "tenant" {
		t.Fatalf("expected tenant, got %q", got)
	}
}

func TestGatewayForEndpoint_PrefersMostSpecificHost(t *testing.T) {
	pkg := Package{
		Spec: Spec{Network: Network{Expose: []Expose{
			{Host: "app", Gateway: "tenant"},
			{Host: "app.admin", Gateway: "admin"},
		}}},
	}
	if got := gatewayForEndpoint(pkg, "app.admin.uds.dev"); got != "admin" {
		t.Fatalf("expected admin, got %q", got)
	}
}

func TestGatewayForEndpoint_UsesGatewaySegmentForSharedHost(t *testing.T) {
	pkg := Package{
		Spec: Spec{Network: Network{Expose: []Expose{
			{Host: "grafana", Gateway: "tenant"},
			{Host: "grafana", Gateway: "admin"},
		}}},
	}
	if got := gatewayForEndpoint(pkg, "grafana.admin.uds.dev"); got != "admin" {
		t.Fatalf("expected admin, got %q", got)
	}
}

func TestGatewayForEndpoint_PrefersTenantForSharedHostWithoutGatewaySegment(t *testing.T) {
	pkg := Package{
		Spec: Spec{Network: Network{Expose: []Expose{
			{Host: "grafana", Gateway: "admin"},
			{Host: "grafana", Gateway: "tenant"},
		}}},
	}
	if got := gatewayForEndpoint(pkg, "grafana.uds.dev"); got != "tenant" {
		t.Fatalf("expected tenant, got %q", got)
	}
}

func TestToAPIApps_GatewayTagging(t *testing.T) {
	tests := []struct {
		name     string
		pkgs     []Package
		wantApps []APIApp
	}{
		{
			name: "admin gateway tagged",
			pkgs: []Package{{
				Metadata: Metadata{Name: "grafana"},
				Status:   Status{Endpoints: []string{"grafana.admin.uds.dev"}},
				Spec: Spec{Network: Network{Expose: []Expose{
					{Host: "grafana.admin", Gateway: "admin"},
				}}},
			}},
			wantApps: []APIApp{{Name: "Grafana", URL: "grafana.admin.uds.dev", Gateway: "admin"}},
		},
		{
			name: "tenant gateway tagged",
			pkgs: []Package{{
				Metadata: Metadata{Name: "podinfo"},
				Status:   Status{Endpoints: []string{"podinfo.uds.dev"}},
				Spec: Spec{Network: Network{Expose: []Expose{
					{Host: "podinfo", Gateway: "tenant"},
				}}},
			}},
			wantApps: []APIApp{{Name: "Podinfo", URL: "podinfo.uds.dev", Gateway: "tenant"}},
		},
		{
			name: "mixed expose yields per-endpoint placement",
			pkgs: []Package{{
				Metadata: Metadata{Name: "mixed"},
				Status:   Status{Endpoints: []string{"front.uds.dev", "back.admin.uds.dev"}},
				Spec: Spec{Network: Network{Expose: []Expose{
					{Host: "front", Gateway: "tenant"},
					{Host: "back.admin", Gateway: "admin"},
				}}},
			}},
			wantApps: []APIApp{
				{Name: "Mixed", URL: "front.uds.dev", Gateway: "tenant"},
				{Name: "Mixed", URL: "back.admin.uds.dev", Gateway: "admin"},
			},
		},
		{
			name: "missing expose defaults to empty gateway",
			pkgs: []Package{{
				Metadata: Metadata{Name: "legacy"},
				Status:   Status{Endpoints: []string{"legacy.uds.dev"}},
			}},
			wantApps: []APIApp{{Name: "Legacy", URL: "legacy.uds.dev"}},
		},
		{
			name: "unmatched host defaults to empty gateway",
			pkgs: []Package{{
				Metadata: Metadata{Name: "drift"},
				Status:   Status{Endpoints: []string{"drift.uds.dev"}},
				Spec: Spec{Network: Network{Expose: []Expose{
					{Host: "other", Gateway: "admin"},
				}}},
			}},
			wantApps: []APIApp{{Name: "Drift", URL: "drift.uds.dev"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toAPIApps(nil, tt.pkgs, "")
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
	got := toAPIApps(nil, nil, "sso.uds.dev")
	if len(got) != 1 {
		t.Fatalf("expected 1 app, got %d", len(got))
	}
	if got[0].Gateway != "" {
		t.Fatalf("expected empty gateway for My Account, got %q", got[0].Gateway)
	}
}
