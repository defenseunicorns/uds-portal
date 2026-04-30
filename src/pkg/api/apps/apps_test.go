// Copyright 2025 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

package apps

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/defenseunicorns/uds-portal/src/pkg/api/auth/incluster"
	"github.com/defenseunicorns/uds-portal/src/pkg/config"
)

func TestEndpointsForExposedHostsMatchesAndDedupes(t *testing.T) {
	got := dedupeEndpoints([]string{
		"keycloak.admin.uds.dev",
		"keycloak.admin.uds.dev",
		"keycloak.uds.dev",
		"grafana.admin.uds.dev",
	})
	want := []string{"keycloak.admin.uds.dev", "keycloak.uds.dev", "grafana.admin.uds.dev"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestToAPIAppsReturnsEmptySlice(t *testing.T) {
	got := toAPIApps(nil, nil)
	if got == nil {
		t.Fatal("expected empty slice, got nil")
	}
	if len(got) != 1 {
		t.Fatalf("expected one app (My Account), got %d", len(got))
	}
	if got[0].URL != myAccountURL {
		t.Fatalf("expected %q first, got %q", myAccountURL, got[0].URL)
	}
}

func TestDedupeEndpointsEmptyInput(t *testing.T) {
	got := dedupeEndpoints(nil)
	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestDedupeEndpointsReturnsUniqueInOrder(t *testing.T) {
	got := dedupeEndpoints([]string{
		"grafana.admin.uds.dev",
		"grafana.admin.uds.dev",
		"grafana.uds.dev",
	})
	want := []string{"grafana.admin.uds.dev", "grafana.uds.dev"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

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

func TestToAPIAppsPrependsMyAccount(t *testing.T) {
	pkg := Package{
		Metadata: Metadata{Name: "grafana"},
		Status:   Status{Endpoints: []string{"grafana.uds.dev"}},
	}

	got := toAPIApps(nil, []Package{pkg})

	if len(got) != 2 {
		t.Fatalf("expected 2 apps, got %d", len(got))
	}
	if got[0].URL != myAccountURL {
		t.Fatalf("expected %q first, got %q", myAccountURL, got[0].URL)
	}
}

func TestFilterByUserGroupRequiresSSOSection(t *testing.T) {
	originalLocalMode := config.LocalMode
	config.LocalMode = false
	t.Cleanup(func() { config.LocalMode = originalLocalMode })

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("failed creating request: %v", err)
	}
	req = req.WithContext(context.WithValue(req.Context(), incluster.GroupKey, []string{"/UDS Core/Admin"}))

	pkgs := []Package{
		{
			Metadata: Metadata{Name: "no-sso"},
			Status:   Status{Endpoints: []string{"no-sso.uds.dev"}},
			Spec:     Spec{SSO: nil},
		},
	}

	got := filterByUserGroup(req, pkgs)
	if len(got) != 0 {
		t.Fatalf("expected package without SSO to be excluded, got %d", len(got))
	}
}

func TestFilterByUserGroupAllowsWhenGroupsNotSpecified(t *testing.T) {
	originalLocalMode := config.LocalMode
	config.LocalMode = false
	t.Cleanup(func() { config.LocalMode = originalLocalMode })

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("failed creating request: %v", err)
	}
	req = req.WithContext(context.WithValue(req.Context(), incluster.GroupKey, []string{"/UDS Core/Auditor"}))

	pkgs := []Package{
		{
			Metadata: Metadata{Name: "sso-without-groups"},
			Status:   Status{Endpoints: []string{"auditor.uds.dev"}},
			Spec: Spec{SSO: []SSO{
				{Groups: Groups{AnyOf: nil}},
			}},
		},
	}

	got := filterByUserGroup(req, pkgs)
	if len(got) != 1 {
		t.Fatalf("expected package with SSO but no group constraints to be included, got %d", len(got))
	}
}

func TestFilterByUserGroupAllowsMatchingGroup(t *testing.T) {
	originalLocalMode := config.LocalMode
	config.LocalMode = false
	t.Cleanup(func() { config.LocalMode = originalLocalMode })

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("failed creating request: %v", err)
	}
	req = req.WithContext(context.WithValue(req.Context(), incluster.GroupKey, []string{"/UDS Core/Auditor", "/UDS Core/Admin"}))

	pkgs := []Package{
		{
			Metadata: Metadata{Name: "group-matched"},
			Status:   Status{Endpoints: []string{"group-matched.uds.dev"}},
			Spec: Spec{SSO: []SSO{
				{Groups: Groups{AnyOf: []string{"/UDS Core/Developer", "/UDS Core/Admin"}}},
			}},
		},
	}

	got := filterByUserGroup(req, pkgs)
	if len(got) != 1 {
		t.Fatalf("expected package with matching group to be included, got %d", len(got))
	}
}
