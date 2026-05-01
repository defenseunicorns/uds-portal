// Copyright 2025 Defense Unicorns
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

var ClassBannerCfg = ClassificationBanners{Enabled: false, Text: "", Footer: false}

const (
	bootstrapConfigNamespace = "window.__APP__"
)

func init() {
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
}

func GenerateBootstrapConfigScript() string {
	var builder strings.Builder
	builder.WriteString("<script>\n")
	builder.WriteString(bootstrapConfigNamespace + " = {};\n")
	builder.WriteString(bootstrapConfigNamespace + ".CLASSIFICATION_BANNER = {};\n")
	fmt.Fprintf(&builder, "%s.CLASSIFICATION_BANNER.enabled = %t;\n", bootstrapConfigNamespace, ClassBannerCfg.Enabled)
	fmt.Fprintf(&builder, "%s.CLASSIFICATION_BANNER.footer = %t;\n", bootstrapConfigNamespace, ClassBannerCfg.Footer)
	fmt.Fprintf(&builder, "%s.CLASSIFICATION_BANNER.text = \"%s\";\n", bootstrapConfigNamespace, template.JSEscapeString(ClassBannerCfg.Text))
	builder.WriteString("</script>\n")

	return builder.String()
}
