// Copyright 2024 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

package middleware

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/defenseunicorns/uds-runtime/src/pkg/api/auth/incluster"
	localAuth "github.com/defenseunicorns/uds-runtime/src/pkg/api/auth/local"
	"github.com/defenseunicorns/uds-runtime/src/pkg/config"
)

// Auth is a middleware that handles all API authentication for UDS Runtime
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if config.LocalAuthEnabled {
			// only /api/ and /swagger behind auth
			if strings.HasPrefix(r.URL.Path, "/api/") {
				if valid := localAuth.ValidateSessionCookie(w, r); valid {
					next.ServeHTTP(w, r)
					return
				}
				// session invalid
				slog.Debug("Session invalid")
				return
			}
		} else if config.InClusterAuthEnabled {
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
