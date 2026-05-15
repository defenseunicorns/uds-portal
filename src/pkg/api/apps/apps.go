// Copyright 2025-2026 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

// Package apps retrieves, filters, and returns UDS packages from the Kubernetes cluster.
package apps

import (
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/defenseunicorns/uds-portal/src/pkg/api/auth/incluster"
	"github.com/defenseunicorns/uds-portal/src/pkg/config"
	"k8s.io/client-go/rest"
)

const (
	udsPortalPkgName = "uds-portal"
	myAccountName    = "My Account"
)

//go:embed icons/my-account.svg
var myAccountIconSVG []byte
var myAccountIcon = "data:image/svg+xml;base64," + base64.StdEncoding.EncodeToString(myAccountIconSVG)

// GetUDSPackages retrieves UDS packages from the cluster and filters them based on user group membership.
func GetUDSPackages(restConfig *rest.Config, inCluster bool, w http.ResponseWriter, r *http.Request) {
	store, err := ensureInformerStore(restConfig)
	if err != nil {
		http.Error(w, "cluster error", http.StatusInternalServerError)
		slog.Error("informer init error", "error", err)
		return
	}

	packages, err := listPackages(store)
	if err != nil {
		http.Error(w, "cluster error", http.StatusInternalServerError)
		slog.Error("package list error", "error", err)
		return
	}

	// filter packages and transform into API response shape
	filteredExposed := filterExposedPackages(packages)
	filteredByGroup := filterByUserGroup(r, filteredExposed, inCluster)
	myAccountURL := ""
	if config.UDSDomain != "" {
		myAccountURL = "sso." + config.UDSDomain
	}
	responseApps := toAPIApps(store, filteredByGroup, myAccountURL, config.UDSDomain)

	// return the filtered packages
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseApps); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
	}
}

func filterExposedPackages(sourcePackages []Package) []Package {
	packages := make([]Package, 0)
	for _, pkg := range sourcePackages {
		if len(pkg.Spec.Network.Expose) > 0 && pkg.Metadata.Name != udsPortalPkgName {
			packages = append(packages, pkg)
		}
	}
	return packages
}

func filterByUserGroup(r *http.Request, packages []Package, inCluster bool) []Package {
	userGroups, _ := r.Context().Value(incluster.GroupKey).([]string)

	if !inCluster {
		return packages
	}

	var filteredByGroup []Package
	for _, pkg := range packages {
		filteredPkg := Package{Metadata: pkg.Metadata, Spec: pkg.Spec, Status: pkg.Status}
		// filter out apps that don't match the user group
		if len(filteredPkg.Spec.Network.Expose) > 0 {
			if len(pkg.Spec.SSO) == 0 {
				continue
			}

			allowed := false
		ssoLoop:
			for _, sso := range pkg.Spec.SSO {
				if len(sso.Groups.AnyOf) == 0 {
					allowed = true
					break ssoLoop
				}

				for _, appGroup := range sso.Groups.AnyOf {
					for _, userGroup := range userGroups {
						if appGroup == userGroup {
							allowed = true
							break ssoLoop
						}
					}
				}
			}

			if allowed {
				filteredByGroup = append(filteredByGroup, filteredPkg)
			}
		}
	}
	return filteredByGroup
}

// endpointURL builds the tile URL from an expose entry's host, gateway, and domain.
// A tenant gateway (or empty gateway) produces host.domain; any other gateway name
// produces host.<gateway>.domain. If domain is empty, just the host is returned.
func endpointURL(host, gateway, domain string) string {
	if domain == "" {
		return host
	}
	if gateway == "" || gateway == "tenant" {
		return host + "." + domain
	}
	return host + "." + gateway + "." + domain
}

func toAPIApps(store *appInformerStore, packages []Package, myAccountURL string, domain string) []APIApp {
	apiApps := make([]APIApp, 0)
	seen := map[string]struct{}{}

	if myAccountURL != "" {
		apiApps = append(apiApps, APIApp{
			Name: myAccountName,
			Icon: myAccountIcon,
			URL:  myAccountURL,
		})
		// pre-seed seen so a package endpoint matching the SSO host isn't listed twice
		seen[myAccountURL] = struct{}{}
	}

	for _, pkg := range packages {
		icon := iconForPackage(store, pkg)
		for _, e := range pkg.Spec.Network.Expose {
			if e.Host == "" {
				continue
			}
			url := endpointURL(e.Host, e.Gateway, domain)
			if _, exists := seen[url]; exists {
				continue
			}
			seen[url] = struct{}{}
			apiApps = append(apiApps, APIApp{
				Name:    displayNameForApp(pkg.Metadata.Name),
				Icon:    icon,
				URL:     url,
				Gateway: e.Gateway,
			})
		}
	}

	return apiApps
}

func displayNameForApp(packageName string) string {
	// TODO: (@wstarr) - this is a temporary function to normalize Package names until a better solution is designed
	normalized := strings.ReplaceAll(strings.TrimSpace(packageName), "-", " ")
	words := strings.Fields(normalized)
	for i, word := range words {
		lower := strings.ToLower(word)
		if lower == "uds" {
			words[i] = "UDS"
			continue
		}

		if lower == "" {
			continue
		}

		words[i] = strings.ToUpper(lower[:1]) + lower[1:]
	}

	return strings.Join(words, " ")
}
