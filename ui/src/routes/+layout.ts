// Copyright 2025 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

import { writable, type Writable } from 'svelte/store'

import { authenticated } from '$features/auth/store'
import type { UserData } from '$features/navigation/types'
import type { ClassBannerCfg } from '$lib/types'

export const ssr = false
export const _bannerCfg: Writable<ClassBannerCfg> = writable({ enabled: false, text: '', footer: false })

interface AuthResponse {
  authenticated: boolean
  userData: UserData
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
          preferredUsername: json['preferred-username'],
          group: json['group'],
          inClusterAuth: json['in-cluster-auth'],
        },
      }
    } else {
      return {
        authenticated: false,
        userData: {
          name: '',
          preferredUsername: '',
          group: '',
          inClusterAuth: false,
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
    preferredUsername: '',
    group: '',
    inClusterAuth: false,
  }

  try {
    const authResult = await auth()

    if (authResult.authenticated) {
      authenticated.set(true)
      userData = authResult.userData

      try {
        const classCfg = await fetch('/class-banners')
        _bannerCfg.set(await classCfg.json())
      } catch (e) {
        console.error(e)
      }
    } else {
      authenticated.set(false)
    }
  } catch (error) {
    console.error('Load error:', error)
    authenticated.set(false)
  }

  return {
    userData,
  }
}
