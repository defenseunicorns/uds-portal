// Copyright 2025 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/defenseunicorns/uds-portal/src/pkg/api/apps"
	"github.com/defenseunicorns/uds-portal/src/pkg/api/auth"
	"github.com/defenseunicorns/uds-portal/src/pkg/api/k8s/session"
)

func getUDSPackages(k8sSession *session.K8sSession) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		apps.GetUDSPackages(k8sSession.Config, w, r)
	}
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	auth.RequestHandler(w, r)
}

func healthz(w http.ResponseWriter, _ *http.Request) {
	slog.Debug("Health check called")

	response := map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("Failed to encode health response", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
