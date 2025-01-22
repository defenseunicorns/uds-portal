// Copyright 2025 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

// Package config contains configuration for the application.
package config

import (
	"crypto/rand"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/defenseunicorns/uds-app-portal/src/pkg/api/auth/local"
)

type ClassificationBanners struct {
	Enabled bool   `json:"enabled"`
	Text    string `json:"text"`
	Footer  bool   `json:"footer"`
}

var (
	LocalAuthEnabled     = true
	InClusterAuthEnabled = false
	ClassBannerCfg       = ClassificationBanners{Enabled: false, Text: "", Footer: false}
)

func init() {
	// configure auth settings
	if os.Getenv("LOCAL_AUTH_ENABLED") == "true" && os.Getenv("IN_CLUSTER_AUTH_ENABLED") == "true" {
		slog.Error("Cannot enable both local and in-cluster auth")
		os.Exit(1)
	}
	if os.Getenv("LOCAL_AUTH_ENABLED") == "false" {
		LocalAuthEnabled = false
	}
	if os.Getenv("IN_CLUSTER_AUTH_ENABLED") == "true" {
		slog.Info("In-cluster auth enabled")
		InClusterAuthEnabled = true
	}

	// If local auth is enabled, generate a token
	if LocalAuthEnabled {
		slog.Info("Local auth enabled")
		token, err := genToken(96)
		if err != nil {
			slog.Error("Failed to generate local auth token")
			os.Exit(1)
		}
		local.AuthToken = token
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

// genToken generates a secure random string of the specified length.
func genToken(length int) (string, error) {
	// safe chars for token generation: https://owasp.org/www-community/password-special-characters
	const randomStringChars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ~-"
	bytes := make([]byte, length)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	for i, b := range bytes {
		bytes[i] = randomStringChars[b%byte(len(randomStringChars))]
	}

	return string(bytes), nil
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
