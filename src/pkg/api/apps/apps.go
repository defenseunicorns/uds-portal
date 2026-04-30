// Copyright 2025 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

// Package apps retrieves, filters, and returns UDS packages from the Kubernetes cluster.
package apps

import (
	"context"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"

	"github.com/defenseunicorns/uds-portal/src/pkg/api/auth/incluster"
	"github.com/defenseunicorns/uds-portal/src/pkg/config"
	zarfConfig "github.com/zarf-dev/zarf/src/config"
	"github.com/zarf-dev/zarf/src/pkg/cluster"
	"github.com/zarf-dev/zarf/src/pkg/state"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

const (
	devUDSIconAnnotation = "dev.uds.icon"
	udsPortalPkgName     = "uds-portal"
	myAccountURL         = "sso.uds.dev"
	myAccountName        = "My Account"
)

//go:embed icons/my-account.svg
var myAccountIconSVG []byte
var myAccountIcon = "data:image/svg+xml;base64," + base64.StdEncoding.EncodeToString(myAccountIconSVG)

var (
	informerStoreMu sync.RWMutex
	informerStore   *appInformerStore
)

type appInformerStore struct {
	packageLister cache.GenericLister
	secretLister  cache.GenericLister
}

// GetUDSPackages retrieves UDS packages from the cluster and filters them based on user group membership.
func GetUDSPackages(config *rest.Config, w http.ResponseWriter, r *http.Request) {
	store, err := ensureInformerStore(config)
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
	filteredByEndpoint := filterPackagesWithEndpoints(packages)
	filteredByGroup := filterByUserGroup(r, filteredByEndpoint)
	responseApps := toAPIApps(store, filteredByGroup)

	// return the filtered packages
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseApps); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
	}
}

func filterPackagesWithEndpoints(sourcePackages []Package) []Package {
	packages := make([]Package, 0)
	for _, pkg := range sourcePackages {
		// ensure package has valid endpoints and is not the UDS portal
		if len(pkg.Status.Endpoints) > 0 && pkg.Metadata.Name != udsPortalPkgName {
			packages = append(packages, pkg)
		}
	}
	return packages
}

func filterByUserGroup(r *http.Request, packages []Package) []Package {
	userGroups, _ := r.Context().Value(incluster.GroupKey).([]string)

	if config.LocalMode {
		return packages
	}

	var filteredByGroup []Package
	for _, pkg := range packages {
		filteredPkg := Package{Metadata: pkg.Metadata, Spec: pkg.Spec, Status: pkg.Status}
		// filter out apps that don't match the user group
		if len(filteredPkg.Status.Endpoints) > 0 {
			if len(pkg.Spec.SSO) == 0 {
				continue
			}

			allowed := false
			for _, sso := range pkg.Spec.SSO {
				if len(sso.Groups.AnyOf) == 0 {
					allowed = true
					break
				}

				for _, appGroup := range sso.Groups.AnyOf {
					for _, userGroup := range userGroups {
						if appGroup == userGroup {
							allowed = true
							break
						}
					}
					if allowed {
						break
					}
				}

				if allowed {
					break
				}
			}

			if allowed {
				filteredByGroup = append(filteredByGroup, filteredPkg)
			}
		}
	}
	return filteredByGroup
}

func toAPIApps(store *appInformerStore, packages []Package) []APIApp {
	apiApps := make([]APIApp, 0)
	seen := map[string]struct{}{myAccountURL: {}}
	for _, pkg := range packages {
		icon := iconForPackage(store, pkg)
		for _, url := range dedupeEndpoints(pkg.Status.Endpoints) {
			if _, exists := seen[url]; exists {
				continue
			}
			seen[url] = struct{}{}
			apiApps = append(apiApps, APIApp{
				Name: displayNameForApp(pkg.Metadata.Name, url),
				Icon: icon,
				URL:  url,
			})
		}
	}

	myAccount := APIApp{
		Name: myAccountName,
		Icon: myAccountIcon,
		URL:  myAccountURL,
	}
	apiApps = append([]APIApp{myAccount}, apiApps...)

	return apiApps
}

func displayNameForApp(packageName, url string) string {
	if url == myAccountURL {
		return myAccountName
	}

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

func dedupeEndpoints(endpoints []string) []string {
	if len(endpoints) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(endpoints))
	var urls []string
	for _, endpoint := range endpoints {
		if _, exists := seen[endpoint]; !exists {
			seen[endpoint] = struct{}{}
			urls = append(urls, endpoint)
		}
	}

	return urls
}

func iconForPackage(store *appInformerStore, pkg Package) string {
	packageLabel := strings.TrimSpace(pkg.Metadata.Labels[cluster.PackageLabel])
	if packageLabel == "" {
		return ""
	}

	namespaceLister := store.secretLister.ByNamespace(state.ZarfNamespaceName)

	secretObj, err := namespaceLister.Get(zarfConfig.ZarfPackagePrefix + packageLabel)
	if err == nil {
		icon := extractIconFromSecret(secretObj)
		if icon != "" {
			return icon
		}
	}

	selector := labels.SelectorFromSet(labels.Set{cluster.PackageLabel: packageLabel})
	secretObjs, err := namespaceLister.List(selector)
	if err != nil {
		slog.Debug("Failed to list package secrets", "package", pkg.Metadata.Name, "namespace", state.ZarfNamespaceName, "error", err)
		return ""
	}

	for _, secretObj := range secretObjs {
		icon := extractIconFromSecret(secretObj)
		if icon != "" {
			return icon
		}
	}

	return ""
}

func ensureInformerStore(cfg *rest.Config) (*appInformerStore, error) {
	informerStoreMu.RLock()
	store := informerStore
	informerStoreMu.RUnlock()
	if store != nil {
		return store, nil
	}

	informerStoreMu.Lock()
	defer informerStoreMu.Unlock()

	if informerStore != nil {
		return informerStore, nil
	}

	udsPackageGVR := schema.GroupVersionResource{Group: "uds.dev", Version: "v1alpha1", Resource: "packages"}
	secretsGVR := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "secrets"}

	dynamicClient, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client: %w", err)
	}

	exists, err := checkCRDExists(dynamicClient, udsPackageGVR)
	if err != nil {
		return nil, fmt.Errorf("failed to check package CRD: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("package CRD not found")
	}

	packageFactory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(dynamicClient, 0, metav1.NamespaceAll, nil)
	secretFactory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(dynamicClient, 0, state.ZarfNamespaceName, nil)

	packageInformer := packageFactory.ForResource(udsPackageGVR)
	secretInformer := secretFactory.ForResource(secretsGVR)

	stopCh := make(chan struct{})
	packageFactory.Start(stopCh)
	secretFactory.Start(stopCh)

	if !cache.WaitForCacheSync(stopCh, packageInformer.Informer().HasSynced, secretInformer.Informer().HasSynced) {
		close(stopCh)
		return nil, fmt.Errorf("failed to sync informer caches")
	}

	informerStore = &appInformerStore{
		packageLister: packageInformer.Lister(),
		secretLister:  secretInformer.Lister(),
	}

	return informerStore, nil
}

func listPackages(store *appInformerStore) ([]Package, error) {
	objects, err := store.packageLister.List(labels.Everything())
	if err != nil {
		return nil, err
	}

	packages := make([]Package, 0, len(objects))
	for _, obj := range objects {
		unstructuredObj, ok := obj.(*unstructured.Unstructured)
		if !ok {
			continue
		}

		pkg := Package{}
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredObj.Object, &pkg); err != nil {
			slog.Error("Failed to unmarshal package", "error", err)
			continue
		}

		packages = append(packages, pkg)
	}

	return packages, nil
}

func extractIconFromSecret(obj runtime.Object) string {
	secret, ok := obj.(*unstructured.Unstructured)
	if !ok || secret == nil {
		return ""
	}

	if data, ok := secret.Object["data"].(map[string]interface{}); ok {
		if encodedPayload, ok := data["data"].(string); ok {
			decodedPayload, err := base64.StdEncoding.DecodeString(encodedPayload)
			if err == nil {
				var payload state.DeployedPackage
				if err := json.Unmarshal(decodedPayload, &payload); err == nil {
					icon := strings.TrimSpace(payload.Data.Metadata.Annotations[devUDSIconAnnotation])
					if icon != "" {
						return icon
					}
				}
			}
		}
	}

	return ""
}

func checkCRDExists(client *dynamic.DynamicClient, gvr schema.GroupVersionResource) (bool, error) {
	_, err := client.Resource(gvr).Namespace("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil // CRD not found
		}
		return false, fmt.Errorf("error retrieving CRD: %v", err)
	}

	return true, nil // CRD exists
}
