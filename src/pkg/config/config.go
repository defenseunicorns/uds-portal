// Copyright 2025 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

// Package config contains configuration for the application.
package config

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
)

type ClassificationBanners struct {
	Enabled bool   `json:"enabled"`
	Text    string `json:"text"`
	Footer  bool   `json:"footer"`
}

var (
	LocalMode            = true
	InClusterAuthEnabled = false
	ClassBannerCfg       = ClassificationBanners{Enabled: false, Text: "", Footer: false}
)

func init() {
	// configure auth settings
	if os.Getenv("LOCAL_MODE") == "true" && os.Getenv("IN_CLUSTER_AUTH_ENABLED") == "true" {
		slog.Error("Cannot enable both local mode and in-cluster auth")
		os.Exit(1)
	}
	if os.Getenv("LOCAL_MODE") == "false" {
		LocalMode = false
	}
	if os.Getenv("IN_CLUSTER_AUTH_ENABLED") == "true" {
		slog.Info("In-cluster auth enabled")
		InClusterAuthEnabled = true
	}

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

func ServeClassBannerCfg() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")

		data, err := json.Marshal(ClassBannerCfg)
		if err != nil {
			http.Error(w, "failed to marshal classification banners", http.StatusInternalServerError)
			return
		}

		if _, err := w.Write(data); err != nil {
			http.Error(w, "failed to write classification banner config", http.StatusInternalServerError)
		}
	}
}
