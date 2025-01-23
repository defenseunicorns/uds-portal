// Copyright 2025 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

// Package apps retrieves, filters, and returns UDS packages from the Kubernetes cluster.
package apps

type Metadata struct {
	Name string `json:"name"`
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

type Spec struct {
	SSO []SSO `json:"sso"`
}

type App struct {
	Metadata Metadata `json:"metadata"`
	Spec     Spec     `json:"spec,omitempty"`
	Status   Status   `json:"status"`
}
