// Copyright 2025 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

package incluster

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

type expectedContext struct {
	Groups            []string
	PreferredUsername string
	Name              string
}

func TestValidateJWT(t *testing.T) {
	parseClaims := func(tokenString string) (jwt.MapClaims, error) {
		token, _, err := jwt.NewParser(jwt.WithoutClaimsValidation()).ParseUnverified(tokenString, jwt.MapClaims{})
		if err != nil {
			return nil, err
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, jwt.ErrTokenInvalidClaims
		}

		return claims, nil
	}

	// Helper function to create a JWT token without signing
	createToken := func(groups []any) string {
		claims := jwt.MapClaims{
			"preferred_username": "testuser",
			"name":               "Test User",
		}
		if groups != nil {
			claims["groups"] = groups
		}
		token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
		tokenString, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
		return tokenString
	}

	tests := []struct {
		name            string
		token           string
		authorization   string
		expectedStatus  int
		expectedContext *expectedContext
	}{
		{
			name:            "Valid token with single group",
			token:           createToken([]any{"/UDS Core/Admin"}),
			authorization:   "Bearer %s",
			expectedStatus:  http.StatusOK,
			expectedContext: &expectedContext{Groups: []string{"/UDS Core/Admin"}, PreferredUsername: "testuser", Name: "Test User"},
		},
		{
			name:            "Valid token with multiple groups",
			token:           createToken([]any{"/UDS Core/Admin", "/Unicorn-Squad"}),
			authorization:   "Bearer %s",
			expectedStatus:  http.StatusOK,
			expectedContext: &expectedContext{Groups: []string{"/UDS Core/Admin", "/Unicorn-Squad"}, PreferredUsername: "testuser", Name: "Test User"},
		},
		{
			name:            "Valid token with empty groups",
			token:           createToken([]any{}),
			authorization:   "Bearer %s",
			expectedStatus:  http.StatusOK,
			expectedContext: &expectedContext{Groups: []string{}, PreferredUsername: "testuser", Name: "Test User"},
		},
		{
			name:           "Invalid token with non-string group",
			token:          createToken([]any{"/UDS Core/Admin", 42}),
			authorization:  "Bearer %s",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:            "Valid token with no groups claim",
			token:           createToken(nil),
			authorization:   "Bearer %s",
			expectedStatus:  http.StatusOK,
			expectedContext: &expectedContext{Groups: []string{}, PreferredUsername: "testuser", Name: "Test User"},
		},
		{
			name:            "Valid token with lowercase bearer scheme",
			token:           createToken([]any{"/UDS Core/Admin"}),
			authorization:   "bearer   %s",
			expectedStatus:  http.StatusOK,
			expectedContext: &expectedContext{Groups: []string{"/UDS Core/Admin"}, PreferredUsername: "testuser", Name: "Test User"},
		},
		{
			name:           "Invalid token",
			token:          "invalid.token.string",
			authorization:  "Bearer %s",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid authorization header format",
			token:          createToken([]any{"/UDS Core/Admin"}),
			authorization:  "Token %s",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Missing token",
			token:          "",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request to pass to our handler
			req, _ := http.NewRequest("GET", "/", nil)
			if tt.token != "" {
				authorizationHeader := "Bearer %s"
				if tt.authorization != "" {
					authorizationHeader = tt.authorization
				}
				req.Header.Set("Authorization", fmt.Sprintf(authorizationHeader, tt.token))
			}

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()

			// Call the function directly
			request, result := validateJWT(rr, req, parseClaims)
			if tt.expectedContext != nil {
				require.Equal(t, tt.expectedContext.Groups, request.Context().Value(GroupKey), "group and user not set together")
				require.Equal(t, tt.expectedContext.PreferredUsername, request.Context().Value(PreferredUserNameKey), "group and user not set together")
				require.Equal(t, tt.expectedContext.Name, request.Context().Value(NameKey), "group and user not set together")
			}

			// Check the status code
			require.Equal(t, tt.expectedStatus, rr.Code, "handler returned wrong status code")

			// Check the return value
			expectedResult := tt.expectedStatus == http.StatusOK
			require.Equal(t, expectedResult, result, "ValidateJWT returned unexpected result")
		})
	}
}
