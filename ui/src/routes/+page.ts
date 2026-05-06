// Copyright 2025 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

import type { ApiApp } from './types'

export const load = async () => {
  try {
    const response = await fetch('/api/v1/apps')
    if (!response.ok) {
      console.error('Failed to fetch apps:', response.status)
      return { apps: [] as ApiApp[] }
    }
    const appData = await response.json()
    if (!Array.isArray(appData)) {
      console.error('Invalid apps response format')
      return { apps: [] as ApiApp[] }
    }
    return { apps: appData as ApiApp[] }
  } catch (error) {
    console.error('Failed to fetch apps:', error)
    return { apps: [] as ApiApp[] }
  }
}
