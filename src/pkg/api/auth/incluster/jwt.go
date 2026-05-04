// Copyright 2025 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

package incluster

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/defenseunicorns/uds-portal/src/pkg/config"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	GroupKey             contextKey = "group"
	PreferredUserNameKey contextKey = "preferred_username"
	NameKey              contextKey = "name"
	oidcClientID                    = "uds-portal"
)

var issuerKeySets = struct {
	mu   sync.RWMutex
	sets map[string]oidc.KeySet
}{
	sets: map[string]oidc.KeySet{},
}

// ValidateJWT checks if the request has a valid JWT token with the required groups.
func ValidateJWT(w http.ResponseWriter, r *http.Request) (*http.Request, bool) {
	return validateJWT(w, r, parseAndValidateJWTClaims)
}

func validateJWT(
	w http.ResponseWriter,
	r *http.Request,
	parseClaims func(string) (jwt.MapClaims, error),
) (*http.Request, bool) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		slog.Debug("JWT validation failed", "reason", "missing authorization header")
		return r, false
	}

	authParts := strings.Fields(authHeader)
	if len(authParts) != 2 || !strings.EqualFold(authParts[0], "Bearer") {
		http.Error(w, "Invalid Authorization header", http.StatusUnauthorized)
		slog.Debug("JWT validation failed", "reason", "invalid authorization header format")
		return r, false
	}

	tokenString := authParts[1]

	claims, err := parseClaims(tokenString)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		slog.Debug("JWT validation failed", "reason", "token verification failed", "error", err)
		return r, false
	}

	// extract groups claim (optional) and set all groups into context
	groups := []string{}
	if rawGroups, exists := claims["groups"]; exists {
		parsedGroups, ok := rawGroups.([]interface{})
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			slog.Debug("JWT validation failed", "reason", "invalid groups claim type")
			return r, false
		}

		for _, group := range parsedGroups {
			groupStr, ok := group.(string)
			if !ok {
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				slog.Debug("JWT validation failed", "reason", "invalid group entry type")
				return r, false
			}
			if groupStr != "" {
				groups = append(groups, groupStr)
			}
		}
	}

	r = r.WithContext(context.WithValue(r.Context(), GroupKey, groups))

	// extract and validate preferred username and name
	preferredUsername, preferredUsernameOk := claims["preferred_username"].(string)
	name, nameOk := claims["name"].(string)

	if !preferredUsernameOk || !nameOk {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		slog.Debug("JWT validation failed", "reason", "missing required user claims")
		return r, false
	}

	// set additional user details in context
	r = r.WithContext(context.WithValue(r.Context(), PreferredUserNameKey, preferredUsername))
	r = r.WithContext(context.WithValue(r.Context(), NameKey, name))

	return r, true
}

func parseAndValidateJWTClaims(tokenString string) (jwt.MapClaims, error) {
	udsDomain := strings.TrimSpace(config.UDSDomain)
	if udsDomain == "" {
		return nil, fmt.Errorf("UDS_DOMAIN is not configured")
	}

	issuer := fmt.Sprintf("https://sso.%s/realms/uds", udsDomain)

	verifier := oidc.NewVerifier(issuer, getIssuerKeySet(issuer), &oidc.Config{
		ClientID:             oidcClientID,
		SupportedSigningAlgs: []string{"RS256", "RS384", "RS512"},
	})

	idToken, err := verifier.Verify(context.Background(), tokenString)
	if err != nil {
		return nil, err
	}

	claims := jwt.MapClaims{}
	if err := idToken.Claims(&claims); err != nil {
		return nil, err
	}

	return claims, nil
}

func getIssuerKeySet(issuer string) oidc.KeySet {
	issuerKeySets.mu.RLock()
	keySet, found := issuerKeySets.sets[issuer]
	issuerKeySets.mu.RUnlock()
	if found {
		return keySet
	}

	issuerKeySets.mu.Lock()
	defer issuerKeySets.mu.Unlock()

	if keySet, found = issuerKeySets.sets[issuer]; found {
		return keySet
	}

	keySet = oidc.NewRemoteKeySet(context.Background(), issuer+"/protocol/openid-connect/certs")
	issuerKeySets.sets[issuer] = keySet

	return keySet
}
