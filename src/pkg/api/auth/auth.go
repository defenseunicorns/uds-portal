// Copyright 2025 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

package auth

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/defenseunicorns/uds-portal/src/pkg/api/auth/incluster"
)

// UserResponse is the response for the /auth endpoint
type UserResponse struct {
	Username string `json:"username"`
	Name     string `json:"name"`
}

// RequestHandler is the main handler for the /auth endpoint; it returns a userResponse struct
// indicating whether the request was authenticated via local or in-cluster auth, and relevant user data
func RequestHandler(w http.ResponseWriter, r *http.Request, inCluster bool) {
	resp := UserResponse{}

	if !inCluster {
		resp.Username = "local"
		resp.Name = "First Last"
	} else {
		// grab values from context set by Auth JWT middleware
		username := r.Context().Value(incluster.PreferredUserNameKey)
		name := r.Context().Value(incluster.NameKey)

		// ensure values are valid
		if username != nil && name != nil {
			resp.Name = name.(string)
			resp.Username = username.(string)
		} else {
			slog.Warn("Failed to get username from context")
			http.Error(w, "authorization failure", http.StatusInternalServerError)
			return
		}
	}

	// write response
	bodyBytes, err := json.Marshal(resp)
	if err != nil {
		slog.Debug("failed to marshal response", "error", err)
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(bodyBytes)
	if err != nil {
		slog.Debug("failed to write response", "error", err)
		http.Error(w, "failed to write response", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
