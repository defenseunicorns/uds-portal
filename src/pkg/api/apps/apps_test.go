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
