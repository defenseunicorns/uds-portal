// Copyright 2025-2026 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

import { writable, type Writable } from 'svelte/store'

import { authenticated } from '$features/auth/store'
import type { UserData } from '$features/navigation/types'
import type { ClassBannerCfg } from '$lib/types'

import type { ApiApp } from './types'

export const ssr = false
export const _bannerCfg: Writable<ClassBannerCfg> = writable({ enabled: false, text: '', footer: false })

type AppBootstrapConfig = {
  CLASSIFICATION_BANNER?: ClassBannerCfg
  ADMIN_APPS_ENABLED?: boolean
}

type BrowserWindowWithAppConfig = Window & {
  __APP__?: AppBootstrapConfig
}

interface AuthResponse {
  authenticated: boolean
  userData: UserData
}

async function fetchApps(): Promise<ApiApp[]> {
  try {
    const response = await fetch('/api/v1/apps')
    if (!response.ok) {
      console.error('Failed to fetch apps:', response.status)
      return []
    }
    const appData = await response.json()
    if (!Array.isArray(appData)) {
      console.error('Invalid apps response format')
      return []
    }
    return appData as ApiApp[]
  } catch (error) {
    console.error('Failed to fetch apps:', error)
    return []
  }
}

// auth function that returns both auth status and user data
async function auth(): Promise<AuthResponse> {
  const baseURL = '/api/v1'
  const headers = new Headers({
    'Content-Type': 'application/json',
  })

  try {
    const url = `${baseURL}/auth`
    const response = await fetch(url, {
      method: 'GET',
      headers,
    })
    if (response.ok) {
      const json = await response.json()
      return {
        authenticated: response.ok,
        userData: {
          name: json['name'],
          username: json['username'],
        },
      }
    } else {
      return {
        authenticated: false,
        userData: {
          name: '',
          username: '',
        },
      }
    }
  } catch (error) {
    console.error('Authentication error:', error)
    throw error // Let the caller handle the error
  }
}

// load namespace and auth data before rendering the app
export const load = async () => {
  let userData: UserData = {
    name: '',
    username: '',
  }

  try {
    const authResult = await auth()

    if (authResult.authenticated) {
      authenticated.set(true)
      userData = authResult.userData

      const appConfig = (window as BrowserWindowWithAppConfig).__APP__
      if (appConfig?.CLASSIFICATION_BANNER) {
        _bannerCfg.set(appConfig.CLASSIFICATION_BANNER)
      }
    } else {
      authenticated.set(false)
    }
  } catch (error) {
    console.error('Load error:', error)
    authenticated.set(false)
  }

  const apps = await fetchApps()
  const appConfig = (window as BrowserWindowWithAppConfig).__APP__
  const adminAppsEnabled: boolean = appConfig?.ADMIN_APPS_ENABLED ?? true

  return {
    userData,
    apps,
    adminAppsEnabled,
  }
}
