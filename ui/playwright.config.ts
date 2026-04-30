// Copyright 2025 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

import { defineConfig, PlaywrightTestConfig } from '@playwright/test'

const TEST_CONFIG = process.env.TEST_CFG || 'default'

export const port = '5173'

const configs: Record<string, PlaywrightTestConfig> = {
  // In-cluster: used in CI against a deployed portal (portal.uds.dev)
  default: {
    name: 'default',
    globalSetup: './tests/in-cluster/global-setup',
    timeout: 20 * 1000,
    testDir: 'tests',
    fullyParallel: true,
    retries: process.env.CI ? 3 : 1,
    testMatch: /^(.+\.)?(test|spec)\.[jt]s$/,
    use: {
      baseURL: 'https://portal.uds.dev/',
      storageState: './tests/state.json',
    },
  },

  // Local: run against the vite dev server (pnpm run dev) at localhost:5173
  local: {
    name: 'local',
    timeout: 20 * 1000,
    testDir: 'tests',
    fullyParallel: true,
    retries: 1,
    testMatch: /^(.+\.)?(test|spec)\.[jt]s$/,
    use: {
      baseURL: `http://localhost:${port}/`,
    },
  },
}

export default defineConfig({ ...configs[TEST_CONFIG] })
