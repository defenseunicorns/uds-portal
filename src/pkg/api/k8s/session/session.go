// Copyright 2024 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

package session

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/defenseunicorns/uds-runtime/src/pkg/k8s/client"
	"github.com/defenseunicorns/uds-runtime/src/pkg/k8s/session"
	"k8s.io/client-go/tools/clientcmd"
)

// CreateK8sSession creates a new k8s session
func CreateK8sSession() (*session.K8sSession, error) {
	k8sClient, err := client.New(&clientcmd.ConfigOverrides{})
	if err != nil {
		return nil, fmt.Errorf("failed to create k8s client: %w", err)
	}

	_, cancel := context.WithCancel(context.Background())

	inCluster, err := client.IsRunningInCluster()
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to check if running in cluster: %w", err)
	}

	var currentCtx, currentCluster string
	if !inCluster { // get the current context and cluster
		currentCtx, currentCluster, err = client.CurrentContext()
		slog.Info("Current context", "context", currentCtx, "cluster", currentCluster)
		if err != nil {
			cancel()
			return nil, fmt.Errorf("failed to get current context: %w", err)
		}
	}

	k8sSession := &session.K8sSession{
		Clients:        k8sClient,
		CurrentCtx:     currentCtx,
		CurrentCluster: currentCluster,
		InCluster:      inCluster,
		Status:         make(chan string),
		Cancel:         cancel,
	}

	return k8sSession, nil
}
