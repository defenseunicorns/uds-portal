// Copyright 2025-2026 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

// Package apps retrieves, filters, and returns UDS packages from the Kubernetes cluster.
package apps

type Metadata struct {
	Name        string            `json:"name"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

type Groups struct {
	AnyOf []string `json:"anyOf"`
}

type Status struct {
	Endpoints []string `json:"endpoints"`
}

type SSO struct {
	Groups Groups `json:"groups"`
}

type Expose struct {
	Host    string `json:"host,omitempty"`
	Gateway string `json:"gateway,omitempty"`
}

type Network struct {
	Expose []Expose `json:"expose,omitempty"`
}

type Spec struct {
	SSO     []SSO   `json:"sso"`
	Network Network `json:"network,omitempty"`
}

type Package struct {
	Metadata Metadata `json:"metadata"`
	Spec     Spec     `json:"spec,omitempty"`
	Status   Status   `json:"status"`
}

type APIApp struct {
	Name    string `json:"name"`
	Icon    string `json:"icon,omitempty"`
	URL     string `json:"url"`
	Gateway string `json:"gateway,omitempty"`
	Group   int    `json:"-"`
}
