// Copyright 2025 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

package middleware

import (
	"log/slog"
	"net/http"

	"github.com/defenseunicorns/uds-portal/src/pkg/api/auth/incluster"
)

// Auth is a middleware that handles all API authentication for UDS Portal
func Auth(inCluster bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if inCluster {
				req, valid := incluster.ValidateJWT(w, r)
				if valid {
					next.ServeHTTP(w, req)
					return
				}
				// token invalid
				slog.Debug("Token invalid")
				return
			}
			// no auth enabled
			next.ServeHTTP(w, r)
		})
	}
}
