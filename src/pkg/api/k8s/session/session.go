// Copyright 2025 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

package session

import (
	"fmt"
	"log/slog"
	"os"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sSession struct {
	Config         *rest.Config
	CurrentCtx     string
	CurrentCluster string
	InCluster      bool
}

// CreateK8sSession creates a new k8s session
func CreateK8sSession() (*K8sSession, error) {
	inCluster := isRunningInCluster()

	config, err := loadK8sConfig(inCluster)
	if err != nil {
		return nil, fmt.Errorf("failed to create k8s config: %w", err)
	}

	var currentCtx, currentCluster string
	if !inCluster { // get the current context and cluster
		currentCtx, currentCluster, err = currentContext()
		slog.Info("Current context", "context", currentCtx, "cluster", currentCluster)
		if err != nil {
			return nil, fmt.Errorf("failed to get current context: %w", err)
		}
	}

	k8sSession := &K8sSession{
		Config:         config,
		CurrentCtx:     currentCtx,
		CurrentCluster: currentCluster,
		InCluster:      inCluster,
	}

	return k8sSession, nil
}

func isRunningInCluster() bool {
	_, hostExists := os.LookupEnv("KUBERNETES_SERVICE_HOST")
	_, portExists := os.LookupEnv("KUBERNETES_SERVICE_PORT")
	return hostExists && portExists
}

func loadK8sConfig(inCluster bool) (*rest.Config, error) {
	if inCluster {
		return rest.InClusterConfig()
	}

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		loadingRules,
		&clientcmd.ConfigOverrides{},
	)

	return clientConfig.ClientConfig()
}

func currentContext() (string, string, error) {
	rawConfig, err := clientcmd.NewDefaultClientConfigLoadingRules().Load()
	if err != nil {
		return "", "", fmt.Errorf("failed to load kubeconfig: %w", err)
	}

	ctxName := rawConfig.CurrentContext
	if ctxName == "" {
		return "", "", fmt.Errorf("no current context set in kubeconfig")
	}

	ctx, exists := rawConfig.Contexts[ctxName]
	if !exists {
		return "", "", fmt.Errorf("current context %q not found in kubeconfig", ctxName)
	}

	return ctxName, ctx.Cluster, nil
}
