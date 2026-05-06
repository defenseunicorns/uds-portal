// Copyright 2025 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

import type { Component } from 'svelte'

export interface BaseRoute {
  name: string
  icon?: Component
  class?: string
  children?: string[]
  path?: string
}

export interface Route {
  name: string
  path: string
  icon?: Component
  class?: string
  children?: RouteChild[]
}

export interface RouteChild {
  name: string
  path: string
}

// UserData is the shape of the user data returned from /api/v1/auth
export interface UserData {
  name: string
  username: string
}
