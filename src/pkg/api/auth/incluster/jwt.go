// Copyright 2025 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

package incluster

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	GroupKey             contextKey = "group"
	PreferredUserNameKey contextKey = "preferred_username"
	NameKey              contextKey = "name"
)

// ValidateJWT checks if the request has a valid JWT token with the required groups.
func ValidateJWT(w http.ResponseWriter, r *http.Request) (*http.Request, bool) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return r, false
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// parse the JWT token without validation
	token, _, err := jwt.NewParser(jwt.WithoutClaimsValidation()).ParseUnverified(tokenString, jwt.Claims(jwt.MapClaims{}))
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return r, false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return r, false
	}

	// extract groups claim (optional) and set all groups into context
	groups := []string{}
	if rawGroups, exists := claims["groups"]; exists {
		parsedGroups, ok := rawGroups.([]interface{})
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return r, false
		}

		for _, group := range parsedGroups {
			groupStr, ok := group.(string)
			if !ok {
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
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
		return r, false
	}

	// set additional user details in context
	r = r.WithContext(context.WithValue(r.Context(), PreferredUserNameKey, preferredUsername))
	r = r.WithContext(context.WithValue(r.Context(), NameKey, name))

	return r, true
}
