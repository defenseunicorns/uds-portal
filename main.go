// Copyright 2025 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

package main

import (
	"embed"
	"log/slog"
	"os"

	"github.com/defenseunicorns/uds-portal/src/pkg/api"
)

//nolint:typecheck // ui/build is generated on demand by build:ui before packaging
//go:embed all:ui/build
var assets embed.FS

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	slog.Info("Setting up API server")
	r, inCluster, err := api.Setup(&assets)
	if err != nil {
		slog.Warn("failed to start the API server", "error", err)
		os.Exit(1)
	}

	err = api.Serve(r, inCluster)
	if err != nil {
		os.Exit(1)
	}
}
