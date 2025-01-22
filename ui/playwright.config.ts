// Copyright 2025 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

import { defineConfig, PlaywrightTestConfig } from '@playwright/test'
import { loadEnv } from 'vite'

const { VITE_PORT_ENV } = loadEnv('dev', process.cwd())
const TEST_CONFIG = process.env.TEST_CFG || 'default'

// use port 8443 because by default we use TLS when running locally
const port = VITE_PORT_ENV ?? '8443'
const protocol = 'https'
const host = 'runtime-local.uds.dev'

const configs: Record<string, PlaywrightTestConfig> = {
  // For all default E2E tests
  default: {
    name: 'default',
    webServer: {
      command: '../build/uds-runtime',
      url: `${protocol}://${host}:${port}`,
      reuseExistingServer: !process.env.CI,
      env: { LOCAL_AUTH_ENABLED: 'false' },
    },
    timeout: 20 * 1000,
    testDir: 'tests',
    fullyParallel: true,
    retries: process.env.CI ? 2 : 1,
    testMatch: /^(?!.*local-auth|.*connections|.*in-cluster)(.+\.)?(test|spec)\.[jt]s$/,
    use: {
      baseURL: `${protocol}://${host}:${port}/`,
    },
  },

  // For testing reconnecting to the cluster
  connections: {
    name: 'connections',
    webServer: {
      command: '../build/uds-runtime',
      url: `${protocol}://${host}:${port}`,
      reuseExistingServer: !process.env.CI,
      env: { LOCAL_AUTH_ENABLED: 'false' },
    },
    timeout: 10 * 1000,
    testDir: 'tests',
    fullyParallel: false,
    retries: process.env.CI ? 2 : 1,
    testMatch: 'connections.spec.ts',
    use: {
      baseURL: `https://runtime-local.uds.dev:${port}/`,
    },
  },

  // For testing local auth only
  localAuth: {
    name: 'local-auth',
    timeout: 60 * 1000,
    testDir: 'tests',
    fullyParallel: false,
    retries: process.env.CI ? 2 : 1,
    testMatch: 'local-auth.spec.ts',
    use: {
      baseURL: `https://runtime-local.uds.dev:${port}/`,
    },
  },

  // For running E2E tests against Runtime in-cluster
  inCluster: {
    name: 'in-cluster',
    globalSetup: './tests/global-setup',
    timeout: 40 * 1000,
    testDir: 'tests',
    fullyParallel: true,
    retries: process.env.CI ? 3 : 1,
    testMatch: /^(?!.*local-auth|.*connections)(.+\.)?(test|spec)\.[jt]s$/,
    use: {
      baseURL: `${protocol}://runtime.admin.uds.dev/`,
      storageState: './tests/state.json',
    },
  },
}

export default defineConfig({ ...configs[TEST_CONFIG] })

export { port }
