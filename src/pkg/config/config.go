// Copyright 2025-2026 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

// Package config contains configuration for the application.
package config

import (
	"fmt"
	"html/template"
	"os"
	"strings"
)

type ClassificationBanners struct {
	Enabled bool   `json:"enabled"`
	Text    string `json:"text"`
	Footer  bool   `json:"footer"`
}

var (
	UDSDomain        = ""
	UDSAdminDomain   = ""
	ClassBannerCfg   = ClassificationBanners{Enabled: false, Text: "", Footer: false}
	AdminAppsEnabled = true
)

const (
	bootstrapConfigNamespace = "window.__APP__"
)

func init() {
	UDSDomain = strings.TrimSpace(os.Getenv("UDS_DOMAIN"))
	UDSAdminDomain = strings.TrimSpace(os.Getenv("UDS_ADMIN_DOMAIN"))

	// Class Banner ENV vars must match the names in chart/templates/deployment.yaml
	bannersEnabled := os.Getenv("CLASSIFICATION_BANNER_ENABLED")
	bannerText := os.Getenv("CLASSIFICATION_BANNER_TEXT")
	bannerFooter := os.Getenv("CLASSIFICATION_BANNER_FOOTER")
	if bannersEnabled == "true" {
		ClassBannerCfg.Enabled = true
	}

	if bannerFooter == "true" {
		ClassBannerCfg.Footer = true
	}

	if bannerText != "" {
		ClassBannerCfg.Text = bannerText
	}

	// ADMIN_APPS_ENABLED env var must match the name in chart/templates/deployment.yaml
	// Only set to false when the env var is explicitly "false"; missing/empty keeps default true.
	AdminAppsEnabled = parseAdminAppsEnabled(os.Getenv("ADMIN_APPS_ENABLED"))
}

// parseAdminAppsEnabled returns true for any value except the literal string "false".
// This means missing, empty, or any other value keeps the default of true.
func parseAdminAppsEnabled(envValue string) bool {
	return envValue != "false"
}

func GenerateBootstrapConfigScript() string {
	var builder strings.Builder
	builder.WriteString("<script>\n")
	builder.WriteString(bootstrapConfigNamespace + " = {};\n")
	builder.WriteString(bootstrapConfigNamespace + ".CLASSIFICATION_BANNER = {};\n")
	fmt.Fprintf(&builder, "%s.CLASSIFICATION_BANNER.enabled = %t;\n", bootstrapConfigNamespace, ClassBannerCfg.Enabled)
	fmt.Fprintf(&builder, "%s.CLASSIFICATION_BANNER.footer = %t;\n", bootstrapConfigNamespace, ClassBannerCfg.Footer)
	fmt.Fprintf(&builder, "%s.CLASSIFICATION_BANNER.text = \"%s\";\n", bootstrapConfigNamespace, template.JSEscapeString(ClassBannerCfg.Text))
	fmt.Fprintf(&builder, "%s.ADMIN_APPS_ENABLED = %t;\n", bootstrapConfigNamespace, AdminAppsEnabled)
	builder.WriteString("</script>\n")

	return builder.String()
}
