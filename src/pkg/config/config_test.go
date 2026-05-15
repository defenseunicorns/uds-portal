// Copyright 2025-2026 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

package config

import (
	"strings"
	"testing"
)

func TestParseAdminAppsEnabled(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"", true},
		{"true", true},
		{"false", false},
		{"False", true},
		{"no", true},
		{"random garbage", true},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := parseAdminAppsEnabled(tt.input)
			if got != tt.want {
				t.Errorf("parseAdminAppsEnabled(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestGenerateBootstrapConfigScript_AdminAppsEnabled_True(t *testing.T) {
	AdminAppsEnabled = true
	script := GenerateBootstrapConfigScript()
	want := "window.__APP__.ADMIN_APPS_ENABLED = true;"
	if !strings.Contains(script, want) {
		t.Errorf("expected script to contain %q, got:\n%s", want, script)
	}
}

func TestGenerateBootstrapConfigScript_AdminAppsEnabled_False(t *testing.T) {
	AdminAppsEnabled = false
	script := GenerateBootstrapConfigScript()
	want := "window.__APP__.ADMIN_APPS_ENABLED = false;"
	if !strings.Contains(script, want) {
		t.Errorf("expected script to contain %q, got:\n%s", want, script)
	}
	// Restore default so other tests are unaffected
	AdminAppsEnabled = true
}
