// Copyright 2025-2026 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

import adapter from '@sveltejs/adapter-static'
import { sveltekit } from '@sveltejs/kit/vite'
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte'
import { defineConfig } from 'vitest/config'

export default defineConfig(({ mode }) => ({
  plugins: [
    sveltekit({
      preprocess: vitePreprocess(),
      adapter: adapter({
        pages: 'build',
        assets: 'build',
        fallback: 'index.html',
        strict: true,
      }),
      version: {
        // Preserve SvelteKit 2's disabled background polling.
        pollInterval: 0,
      },
      onwarn: (warning, handler) => {
        if (warning.code.startsWith('a11y-')) {
          throw new Error(warning.message)
        }
        handler(warning)
      },
    }),
  ],
  server: {
    proxy: {
      // Proxy all requests starting with /api to the go server
      // noting that we use http and 8080 in local mode
      '/api': {
        target: 'http://127.0.0.1:8080',
        changeOrigin: true,
      },
      '/healthz': {
        target: 'http://127.0.0.1:8080',
        changeOrigin: true,
      },
    },
  },
  test: {
    include: ['src/**/*.test.{js,ts}'],
    environment: 'jsdom',
    globals: true,
    setupFiles: ['src/setupTests.ts'],
  },
  resolve: mode === 'test' ? { conditions: ['browser'] } : undefined,
}))
