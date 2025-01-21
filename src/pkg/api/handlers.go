// Copyright 2024 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/defenseunicorns/uds-app-portal/src/pkg/api/auth"
	"github.com/defenseunicorns/uds-app-portal/src/pkg/config"
	"github.com/defenseunicorns/uds-runtime/src/pkg/api/rest"
	"github.com/defenseunicorns/uds-runtime/src/pkg/k8s/session"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

// todo: helper fn, move to utils or something
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

func checkClusterConnection(k8sSession *session.K8sSession) http.HandlerFunc {
	return k8sSession.ServeConnStatus()
}

func getUDSPackages(k8sSession *session.K8sSession) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// define GVR for UDS Packages
		udsPackageGVR := schema.GroupVersionResource{
			Group:    "uds.dev",
			Version:  "v1alpha1",
			Resource: "packages",
		}

		// create dynamic client to get CRD directly and check if the CRD exists
		dynamicClient, err := dynamic.NewForConfig(k8sSession.Clients.Config)
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
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	auth.RequestHandler(w, r)
}

func getClassBannerCfg() func(w http.ResponseWriter, r *http.Request) {
	return config.ServeClassBannerCfg()
}

func healthz(w http.ResponseWriter, _ *http.Request) {
	slog.Debug("Health check called")

	response := map[string]interface{}{
		"status":    "UP",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("Failed to encode health response", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
