// Copyright 2025 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

package apps

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"sync"

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

const devUDSIconAnnotation = "dev.uds.icon"

var (
	informerStoreMu sync.RWMutex
	informerStore   *appInformerStore
)

type appInformerStore struct {
	packageLister cache.GenericLister
	secretLister  cache.GenericLister
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

func iconForPackage(store *appInformerStore, pkg Package) string {
	packageLabel := strings.TrimSpace(pkg.Metadata.Labels[cluster.PackageLabel])
	if packageLabel == "" {
		return ""
	}

	zarfSecretLister := store.secretLister.ByNamespace(state.ZarfNamespaceName)

	secretName := secretNameForPackage(pkg)
	secretObj, err := zarfSecretLister.Get(secretName)
	if err == nil {
		icon := extractIconFromSecret(secretObj)
		if icon != "" {
			return icon
		}
	}

	return ""
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

func secretNameForPackage(pkg Package) string {
	packageLabel := strings.TrimSpace(pkg.Metadata.Labels[cluster.PackageLabel])
	namespaceOverride := strings.TrimSpace(pkg.Metadata.Labels[cluster.NamespaceOverrideLabel])

	deployedPackage := &state.DeployedPackage{
		Name:              packageLabel,
		NamespaceOverride: namespaceOverride,
	}

	return deployedPackage.GetSecretName()
}

func checkCRDExists(client *dynamic.DynamicClient, gvr schema.GroupVersionResource) (bool, error) {
	_, err := client.Resource(gvr).Namespace("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}
		return false, fmt.Errorf("error retrieving CRD: %v", err)
	}

	return true, nil
}
