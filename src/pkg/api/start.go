// Copyright 2025 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

package api

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/defenseunicorns/pkg/exec"
	"github.com/defenseunicorns/uds-portal/src/pkg/api/k8s/session"
	udsMiddleware "github.com/defenseunicorns/uds-portal/src/pkg/api/middleware"
	"github.com/defenseunicorns/uds-portal/src/pkg/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Setup(assets *embed.FS) (*chi.Mux, bool, error) {
	k8sSession, err := session.CreateK8sSession()
	if err != nil {
		return nil, false, fmt.Errorf("failed to setup k8s session: %w", err)
	}

	inCluster := k8sSession.InCluster

	r := chi.NewRouter()

	// Add middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Health check is unauthenticated
	r.Get("/healthz", healthz)

	r.Group(func(r chi.Router) {
		r.Use(udsMiddleware.Auth(inCluster))
		r.Route("/api/v1", func(r chi.Router) {
			r.Get("/auth", authHandler(inCluster))
			r.Get("/apps", getUDSPackages(k8sSession))
		})
	})

	// launch app in local mode
	if !inCluster {
		port := "8080"
		host := "127.0.0.1"
		colorYellow := "\033[33m"
		colorReset := "\033[0m"
		url := fmt.Sprintf("http://%s:%s", host, port)
		log.Printf("%sConnect to UDS Portal: %s%s", colorYellow, url, colorReset)
		err := exec.LaunchURL(url)
		if err != nil {
			// if the browser fails to open we intentionally fail the start command
			// local mode is intended for local development where a browser is available
			// deploying outside of kubernetes for production has security implications
			// and should be avoided
			return nil, inCluster, fmt.Errorf("failed to launch URL: %w", err)
		}
	}

	// Serve static files from embed.FS
	if assets != nil {
		staticFS, err := fs.Sub(assets, "ui/build")
		if err != nil {
			return nil, inCluster, fmt.Errorf("failed to create static file system: %w", err)
		}

		if err := fileServer(r, http.FS(staticFS)); err != nil {
			return nil, inCluster, fmt.Errorf("failed to serve static files: %w", err)
		}
	}
	return r, inCluster, nil
}

// fileServer is a custom file server handler for embedded files
func fileServer(r chi.Router, root http.FileSystem) error {
	const bootstrapConfigMarker = "<!-- uds-portal-bootstrap-config -->"

	// Load index.html content and modification time at startup
	f, err := root.Open("index.html")
	if err != nil {
		return err
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return err
	}
	indexModTime := stat.ModTime()

	indexHTML, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	indexWithBootstrap := strings.Replace(string(indexHTML), bootstrapConfigMarker, config.GenerateBootstrapConfigScript(), 1)

	// Create a new file server handler
	fsHandler := http.FileServer(root)

	// Serve the index.html file if the requested file doesn't exist
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		requestPath := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
		if requestPath == "." || requestPath == "" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			http.ServeContent(w, r, "index.html", indexModTime, strings.NewReader(indexWithBootstrap))
			return
		}

		// Try to open the file from the embedded filesystem
		file, err := root.Open(requestPath)
		if err != nil {
			// If the file doesn't exist, serve the pre-loaded index.html
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			// Serve the index.html file with the pre-loaded content
			http.ServeContent(w, r, "index.html", indexModTime, strings.NewReader(indexWithBootstrap))
			return
		}
		file.Close()

		// If the file exists, serve it using the http.FileServer
		fsHandler.ServeHTTP(w, r)
	})

	return nil
}

func Serve(r *chi.Mux, inCluster bool) error {
	server := &http.Server{
		Addr:              "127.0.0.1:8080",
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}

	if inCluster {
		slog.Info("Starting server in in-cluster mode on 0.0.0.0:8080")

		server.Addr = "0.0.0.0:8080"
		if err := server.ListenAndServe(); err != nil {
			slog.Warn("server failed to start", "error", err)
			return err
		}

		return nil
	}

	slog.Info("Starting server in local mode on 127.0.0.1:8080")

	if err := server.ListenAndServe(); err != nil {
		slog.Warn("server failed to start", "err", err)
		return err
	}

	return nil
}
