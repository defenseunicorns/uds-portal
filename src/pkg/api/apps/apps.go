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
	"sort"
	"strings"

	"github.com/defenseunicorns/uds-portal/src/pkg/api/auth/incluster"
	"github.com/defenseunicorns/uds-portal/src/pkg/config"
	"k8s.io/client-go/rest"
)

const (
	udsPortalPkgName         = "uds-portal"
	myAccountName            = "My Account"
	portalHideAppsAnnotation = "uds.dev/portal-hide-apps"

	portalNSVisibleAnnotation = "portal.uds.dev/visible"
	udsDevVisibleAnnotation   = "uds.dev/visible"
	portalNSTitleAnnotation   = "portal.uds.dev/title"
	udsDevTitleAnnotation     = "uds.dev/title"
	portalNSIconAnnotation    = "portal.uds.dev/icon"
	udsDevIconAnnotation      = "uds.dev/icon"

	groupMyAccount   = 0
	groupUDSPlatform = 10
	groupUDSCore     = 20
	groupOther       = 30
)

var udsPlatformPackages = map[string]struct{}{
	"cat":          {},
	"uds-registry": {},
	"uds-ui":       {},
}

// udsCorePackages lists uds-core packages that expose a gateway endpoint visible in the portal.
// Other uds-core packages (loki, prometheus-stack, velero, etc.) have no expose entries and never produce a tile.
var udsCorePackages = map[string]struct{}{
	"grafana": {},
}

func groupForPackage(pkg Package) int {
	if _, ok := udsPlatformPackages[pkg.Metadata.Name]; ok {
		return groupUDSPlatform
	}
	if _, ok := udsCorePackages[pkg.Metadata.Name]; ok {
		return groupUDSCore
	}
	return groupOther
}

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
	filteredHidden := filterHiddenEndpoints(filteredExposed)
	filteredByGroup := filterByUserGroup(r, filteredHidden, inCluster)
	myAccountURL := ""
	if config.UDSDomain != "" {
		myAccountURL = "sso." + config.UDSDomain
	}
	responseApps := toAPIApps(store, filteredByGroup, myAccountURL, config.UDSDomain, config.UDSAdminDomain)

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

func filterHiddenEndpoints(packages []Package) []Package {
	result := make([]Package, 0, len(packages))
	for _, pkg := range packages {
		hidden := hiddenHostSet(pkg)
		var filtered []Expose
		for _, e := range pkg.Spec.Network.Expose {
			if strings.Contains(e.Host, "*") {
				continue
			}
			if _, ok := hidden[e.Host]; ok {
				continue
			}
			if !isEndpointVisible(e) {
				continue
			}
			filtered = append(filtered, e)
		}
		if len(filtered) == 0 {
			continue
		}
		pkg.Spec.Network.Expose = filtered
		result = append(result, pkg)
	}
	return result
}

func hiddenHostSet(pkg Package) map[string]struct{} {
	hostSet := map[string]struct{}{}
	raw := strings.TrimSpace(pkg.Metadata.Annotations[portalHideAppsAnnotation])
	if raw == "" {
		return hostSet
	}
	for _, h := range strings.Split(raw, ",") {
		if h = strings.TrimSpace(h); h != "" {
			hostSet[h] = struct{}{}
		}
	}
	return hostSet
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

// endpointURL builds the tile URL from an expose entry's host, gateway, tenantDomain, and adminDomain.
// Returns "" when the URL cannot be determined (i.e., the required domain is empty), signaling that
// the tile should be skipped.
//
//   - tenant gateway (or empty): uses tenantDomain; returns "" if tenantDomain is empty.
//   - admin gateway: uses adminDomain when set; falls back to "admin."+tenantDomain when only
//     tenantDomain is set; returns "" when both are empty.
//   - custom gateway: uses host.<gateway>.<tenantDomain>; returns "" if tenantDomain is empty.
func endpointURL(host, gateway, tenantDomain, adminDomain string) string {
	switch gateway {
	case "", "tenant":
		if tenantDomain == "" {
			return ""
		}
		return host + "." + tenantDomain
	case "admin":
		domain := adminDomain
		if domain == "" && tenantDomain != "" {
			domain = "admin." + tenantDomain
		}
		if domain == "" {
			return ""
		}
		return host + "." + domain
	default:
		// Custom gateway: host.<gateway>.<tenantDomain>; skip if domain empty.
		if tenantDomain == "" {
			return ""
		}
		return host + "." + gateway + "." + tenantDomain
	}
}

func toAPIApps(
	store *appInformerStore,
	packages []Package,
	myAccountURL string,
	tenantDomain, adminDomain string,
) []APIApp {
	apiApps := make([]APIApp, 0)
	seen := map[string]struct{}{}

	if myAccountURL != "" {
		apiApps = append(apiApps, APIApp{
			Name:  myAccountName,
			Icon:  myAccountIcon,
			URL:   myAccountURL,
			Group: groupMyAccount,
		})
		// pre-seed seen so a package endpoint matching the SSO host isn't listed twice
		seen[myAccountURL] = struct{}{}
	}

	for _, pkg := range packages {
		meta := metadataForPackage(store, pkg)
		group := groupForPackage(pkg)
		for _, e := range pkg.Spec.Network.Expose {
			if e.Host == "" {
				continue
			}
			url := endpointURL(e.Host, e.Gateway, tenantDomain, adminDomain)
			if url == "" {
				continue
			}
			if _, exists := seen[url]; exists {
				continue
			}
			seen[url] = struct{}{}
			title := firstNonEmpty(e.Annotations[portalNSTitleAnnotation], e.Annotations[udsDevTitleAnnotation], meta.title)
			icon := firstNonEmpty(e.Annotations[portalNSIconAnnotation], e.Annotations[udsDevIconAnnotation], meta.icon)
			apiApps = append(apiApps, APIApp{
				Name:    displayNameForApp(title, pkg.Metadata.Name),
				Icon:    icon,
				URL:     url,
				Gateway: e.Gateway,
				Group:   group,
			})
		}
	}

	sort.Slice(apiApps, func(i, j int) bool {
		if apiApps[i].Group != apiApps[j].Group {
			return apiApps[i].Group < apiApps[j].Group
		}
		if apiApps[i].Name != apiApps[j].Name {
			return apiApps[i].Name < apiApps[j].Name
		}
		return apiApps[i].URL < apiApps[j].URL
	})

	return apiApps
}

// isEndpointVisible returns false only when the first visibility annotation found
// (portal.uds.dev/visible checked before uds.dev/visible) is explicitly "false"
// (case-insensitive). Missing annotation defaults to visible.
func isEndpointVisible(e Expose) bool {
	for _, key := range []string{portalNSVisibleAnnotation, udsDevVisibleAnnotation} {
		if val, ok := e.Annotations[key]; ok {
			return strings.ToLower(strings.TrimSpace(val)) != "false"
		}
	}
	return true
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v = strings.TrimSpace(v); v != "" {
			return v
		}
	}
	return ""
}

func formatPackageName(packageName string) string {
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

func displayNameForApp(title, packageName string) string {
	if title != "" {
		return title
	}
	return formatPackageName(packageName)
}
