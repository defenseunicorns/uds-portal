// Copyright 2025 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

package apps

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/zarf-dev/zarf/src/pkg/cluster"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/cache"
)

func TestListPackages(t *testing.T) {
	tests := []struct {
		name        string
		listObjects []runtime.Object
		listErr     error
		wantLen     int
		wantErr     bool
	}{
		{
			name: "returns packages from valid unstructured objects",
			listObjects: []runtime.Object{
				&unstructured.Unstructured{Object: map[string]interface{}{
					"metadata": map[string]interface{}{
						"name":   "podinfo",
						"labels": map[string]interface{}{"zarf.dev/package": "podinfo"},
					},
					"status": map[string]interface{}{"endpoints": []interface{}{"podinfo.uds.dev"}},
				}},
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name:    "returns error when lister returns error",
			listErr: fmt.Errorf("list failed"),
			wantLen: 0,
			wantErr: true,
		},
		{
			name: "skips non-unstructured objects",
			listObjects: []runtime.Object{
				&runtime.Unknown{},
			},
			wantLen: 0,
			wantErr: false,
		},
		{
			name: "skips objects that fail package conversion",
			listObjects: []runtime.Object{
				&unstructured.Unstructured{Object: map[string]interface{}{
					"metadata": "invalid-metadata-shape",
					"status":   map[string]interface{}{"endpoints": []interface{}{"bad.uds.dev"}},
				}},
			},
			wantLen: 0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &appInformerStore{
				packageLister: fakeGenericLister{
					listObjects: tt.listObjects,
					listErr:     tt.listErr,
				},
			}

			got, err := listPackages(store)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
			if len(got) != tt.wantLen {
				t.Fatalf("expected %d packages, got %d", tt.wantLen, len(got))
			}
		})
	}
}

func TestIconForPackage(t *testing.T) {
	const expectedIcon = "data:image/svg+xml;base64,icon"

	pkgWithoutLabel := Package{}
	pkgWithLabel := Package{Metadata: Metadata{Labels: map[string]string{cluster.PackageLabel: "podinfo"}}}
	pkgWithOverride := Package{Metadata: Metadata{Labels: map[string]string{
		cluster.PackageLabel:           "podinfo",
		cluster.NamespaceOverrideLabel: "podinfo-system",
	}}}

	tests := []struct {
		name       string
		pkg        Package
		nilStore   bool
		getObjects map[string]runtime.Object
		want       string
	}{
		{
			name:     "returns empty when store is nil",
			pkg:      pkgWithLabel,
			nilStore: true,
			want:     "",
		},
		{
			name: "returns empty when package label missing",
			pkg:  pkgWithoutLabel,
			want: "",
		},
		{
			name: "returns icon from direct secret lookup",
			pkg:  pkgWithLabel,
			getObjects: map[string]runtime.Object{
				secretNameForPackage(pkgWithLabel): secretObjectWithIcon(expectedIcon),
			},
			want: expectedIcon,
		},
		{
			name: "returns empty when direct secret lookup misses",
			pkg:  pkgWithLabel,
			want: "",
		},
		{
			name: "uses override secret name when namespace override label present",
			pkg:  pkgWithOverride,
			getObjects: map[string]runtime.Object{
				secretNameForPackage(pkgWithOverride): secretObjectWithIcon(expectedIcon),
			},
			want: expectedIcon,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var store *appInformerStore
			if !tt.nilStore {
				store = &appInformerStore{
					secretLister: fakeGenericLister{
						namespaceLister: fakeNamespaceLister{
							getObjects: tt.getObjects,
						},
					},
				}
			}

			got := iconForPackage(store, tt.pkg)
			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestSecretNameForPackage(t *testing.T) {
	tests := []struct {
		name string
		pkg  Package
		want string
	}{
		{
			name: "without namespace override uses base secret name",
			pkg: Package{Metadata: Metadata{Labels: map[string]string{
				cluster.PackageLabel: "podinfo",
			}}},
			want: "zarf-package-podinfo",
		},
		{
			name: "with namespace override uses override suffix",
			pkg: Package{Metadata: Metadata{Labels: map[string]string{
				cluster.PackageLabel:           "podinfo",
				cluster.NamespaceOverrideLabel: "podinfo-system",
			}}},
			want: "zarf-package-podinfo-override-podinfo-system",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := secretNameForPackage(tt.pkg)
			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

type fakeGenericLister struct {
	namespaceLister cache.GenericNamespaceLister
	listObjects     []runtime.Object
	listErr         error
}

func (f fakeGenericLister) List(_ labels.Selector) (ret []runtime.Object, err error) {
	if f.listErr != nil {
		return nil, f.listErr
	}
	return f.listObjects, nil
}

func (f fakeGenericLister) Get(_ string) (runtime.Object, error) {
	return nil, fmt.Errorf("not implemented")
}

func (f fakeGenericLister) ByNamespace(_ string) cache.GenericNamespaceLister {
	return f.namespaceLister
}

type fakeNamespaceLister struct {
	getObjects  map[string]runtime.Object
	listObjects []runtime.Object
	listErr     error
}

func (f fakeNamespaceLister) List(_ labels.Selector) (ret []runtime.Object, err error) {
	if f.listErr != nil {
		return nil, f.listErr
	}
	return f.listObjects, nil
}

func (f fakeNamespaceLister) Get(name string) (runtime.Object, error) {
	if obj, found := f.getObjects[name]; found {
		return obj, nil
	}
	return nil, fmt.Errorf("%s not found", name)
}

func secretObjectWithIcon(icon string) runtime.Object {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"metadata": map[string]interface{}{
				"annotations": map[string]interface{}{
					devUDSIconAnnotation: icon,
				},
			},
		},
	}
	payloadBytes, _ := json.Marshal(payload)

	return &unstructured.Unstructured{Object: map[string]interface{}{
		"data": map[string]interface{}{
			"data": base64.StdEncoding.EncodeToString(payloadBytes),
		},
	}}
}
