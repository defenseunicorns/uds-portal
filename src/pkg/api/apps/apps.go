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

	"github.com/defenseunicorns/uds-app-portal/src/pkg/api/auth/incluster"
	"github.com/defenseunicorns/uds-runtime/src/pkg/api/rest"
	"github.com/defenseunicorns/uds-runtime/src/pkg/k8s/client"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

func GetUDSPackages(client *client.Clients, w http.ResponseWriter, r *http.Request) {
	udsPackageGVR := schema.GroupVersionResource{
		Group:    "uds.dev",
		Version:  "v1alpha1",
		Resource: "packages",
	}

	// create dynamic client to get CRD directly and check if the CRD exists
	dynamicClient, err := dynamic.NewForConfig(client.Config)
	exists, err := checkCRDExists(dynamicClient, udsPackageGVR)
	if err != nil {
		return
	}

	// return packages if the CRD exists
	if !exists {
		http.Error(w, "CRD not found", http.StatusNotFound)
		return
	}

	// get all Packages using the dynamic client
	packages, err := dynamicClient.Resource(udsPackageGVR).List(r.Context(), metav1.ListOptions{})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get packages: %v", err), http.StatusInternalServerError)
		return
	}

	// todo: filter packages further based on group

	// retrieve groups from auth middleware context
	group := r.Context().Value(incluster.GroupKey)
	username := r.Context().Value(incluster.PreferredUserNameKey)
	name := r.Context().Value(incluster.NameKey)

	// log group, username, and name
	slog.Warn("Group", "group", group)
	slog.Warn("Username", "username", username)
	slog.Warn("Name", "name", name)

	// filter the Packages to send only the necessary fields
	fields := []string{"metadata.name", "spec.sso.groups.anyOf", "status.endpoints"}
	filteredFields := rest.FilterItemsByFields(packages.Items, fields)

	type Status struct {
		Endpoints []string `json:"endpoints"`
	}

	// remove packages without endpoints
	var filteredData []map[string]interface{}
	for _, item := range filteredFields {
		// marshall unstructured status into Status struct to check endpoints
		status := Status{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(item["status"].(map[string]interface{}), &status)
		if err != nil {
			slog.Error("Failed to unmarshal status", "error", err)
			continue
		}
		if len(status.Endpoints) > 0 {
			filteredData = append(filteredData, item)
		}
	}

	// convert to JSON and write response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(filteredData); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
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
