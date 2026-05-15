// Copyright 2025-2026 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

package config

import (
	"strings"
	"testing"
)

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
