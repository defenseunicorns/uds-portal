// Copyright 2025 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

// Package apps retrieves, filters, and returns UDS packages from the Kubernetes cluster.
package apps

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/defenseunicorns/uds-app-portal/src/pkg/api/auth/incluster"
	"github.com/defenseunicorns/uds-app-portal/src/pkg/config"
	"github.com/defenseunicorns/uds-runtime/src/pkg/api/rest"
	"github.com/defenseunicorns/uds-runtime/src/pkg/k8s/client"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

// GetUDSPackages retrieves UDS packages from the cluster and filters them based on user group membership.
func GetUDSPackages(client *client.Clients, w http.ResponseWriter, r *http.Request) {
	udsPackageGVR := schema.GroupVersionResource{
		Group:    "uds.dev",
		Version:  "v1alpha1",
		Resource: "packages",
	}

	// create client and ensure CRD exists
	dynamicClient, err := dynamic.NewForConfig(client.Config)
	if err != nil {
		http.Error(w, "cluster error", http.StatusInternalServerError)
		slog.Error("dynamic client error", "error", err)
		return
	}
	exists, err := checkCRDExists(dynamicClient, udsPackageGVR)
	if err != nil {
		http.Error(w, "crd error", http.StatusInternalServerError)
		slog.Error("error checking CRD", "error", err)
		return
	}
	if !exists {
		http.Error(w, "crd error", http.StatusInternalServerError)
		slog.Error("CRD not found", "error", err)
		return
	}

	// get packages
	packages, err := dynamicClient.Resource(udsPackageGVR).List(r.Context(), metav1.ListOptions{})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get packages: %v", err), http.StatusInternalServerError)
		return
	}

	// filter packages
	filteredByEndpoint := filterPackages(packages)
	filteredByGroup := filterByUserGroup(r, filteredByEndpoint)

	// return the filtered packages
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(filteredByGroup); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
	}
}

func filterPackages(packages *unstructured.UnstructuredList) []App {
	fields := []string{"metadata.name", "spec.sso[].groups.anyOf", "status.endpoints"}
	filteredFields := rest.FilterItemsByFields(packages.Items, fields)

	var apps []App
	for _, item := range filteredFields {
		app := App{}
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(item, &app); err != nil {
			slog.Error("Failed to unmarshal package", "error", err)
			continue
		}
		// ensure app has valid endpoints and is not the admin portal
		if len(app.Status.Endpoints) > 0 && app.Metadata.Name != "uds-app-portal" {
			apps = append(apps, app)
		}
	}
	return apps
}

func filterByUserGroup(r *http.Request, apps []App) []App {
	userGroup := r.Context().Value(incluster.GroupKey)

	if !config.InClusterAuthEnabled || userGroup == "/UDS Core/Admin" || userGroup == "/UDS Core/Auditor" {
		return apps
	}

	var filteredByGroup []App
	for _, app := range apps {
		var filteredApp App
		// filter out apps on the admin gw
		for _, endpoint := range app.Status.Endpoints {
			if !strings.Contains(endpoint, ".admin.") {
				filteredApp.Status.Endpoints = append(filteredApp.Status.Endpoints, endpoint)
			}
		}
		// filter out apps that don't match the user group
		if len(filteredApp.Status.Endpoints) > 0 {
			if len(app.Spec.SSO) > 0 {
				for _, sso := range app.Spec.SSO {
					for _, appGroup := range sso.Groups.AnyOf {
						if appGroup == userGroup {
							// user is in the app's group
							filteredApp.Metadata = app.Metadata
							filteredByGroup = append(filteredByGroup, filteredApp)
							break
						}
					}
				}
			} else {
				// no SSO groups are defined, allow all users
				filteredApp.Metadata = app.Metadata
				filteredByGroup = append(filteredByGroup, filteredApp)
			}
		}
	}
	return filteredByGroup
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
