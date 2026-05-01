// Copyright 2025 Defense Unicorns
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
		url         string
		want        string
	}{
		{
			name:        "formats dashed name and title case",
			packageName: "mission-control",
			url:         "mission.uds.dev",
			want:        "Mission Control",
		},
		{
			name:        "normalizes uds acronym",
			packageName: "uds-registry",
			url:         "registry.uds.dev",
			want:        "UDS Registry",
		},
		{
			name:        "my account special case",
			packageName: "anything",
			url:         myAccountURL,
			want:        myAccountName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := displayNameForApp(tt.packageName, tt.url)
			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestToAPIApps(t *testing.T) {
	tests := []struct {
		name    string
		pkgs    []Package
		wantLen int
	}{
		{
			name:    "returns My Account when no packages",
			pkgs:    nil,
			wantLen: 1,
		},
		{
			name: "prepends My Account when packages exist",
			pkgs: []Package{
				{
					Metadata: Metadata{Name: "grafana"},
					Status:   Status{Endpoints: []string{"grafana.uds.dev"}},
				},
			},
			wantLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toAPIApps(nil, tt.pkgs)
			if got == nil {
				t.Fatal("expected empty slice, got nil")
			}
			if len(got) != tt.wantLen {
				t.Fatalf("expected %d apps, got %d", tt.wantLen, len(got))
			}
			if got[0].URL != myAccountURL {
				t.Fatalf("expected %q first, got %q", myAccountURL, got[0].URL)
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
