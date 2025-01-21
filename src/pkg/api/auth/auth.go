// Copyright 2024 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

package auth

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/defenseunicorns/uds-runtime/src/pkg/api/auth/incluster"
	"github.com/defenseunicorns/uds-runtime/src/pkg/api/auth/local"
	"github.com/defenseunicorns/uds-runtime/src/pkg/config"
)

// UserResponse is the response for the /auth endpoint
type UserResponse struct {
	Group             string `json:"group"`
	PreferredUsername string `json:"preferred-username"`
	Name              string `json:"name"`
	InClusterAuth     bool   `json:"in-cluster-auth"`
}

// RequestHandler is the main handler for the /auth endpoint; it returns a userResponse struct
// indicating whether the request was authenticated via local or in-cluster auth, and relevant user data
func RequestHandler(w http.ResponseWriter, r *http.Request) {
	resp := UserResponse{
		InClusterAuth: false,
	}
	if config.LocalAuthEnabled && !local.Auth(w, r) {
		// auth failed, response is already written
		return
	}

	if config.InClusterAuthEnabled {
		// grab values from context set by Auth JWT middleware
		group := r.Context().Value(incluster.GroupKey)
		username := r.Context().Value(incluster.PreferredUserNameKey)
		name := r.Context().Value(incluster.NameKey)

		resp.InClusterAuth = true

		// ensure values are valid
		if group != nil && username != nil && name != nil {
			resp.Group = group.(string)
			resp.Name = name.(string)
			resp.PreferredUsername = username.(string)
		} else {
			slog.Warn("Failed to get group and username from context")
			http.Error(w, "authorization failure", http.StatusInternalServerError)
			return
		}
	}

	// write response
	w.WriteHeader(http.StatusOK)
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
}
