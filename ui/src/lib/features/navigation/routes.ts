// Copyright 2025-2026 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

import { Apps, Settings } from 'carbon-icons-svelte'

import { isAdminGateway } from '../apps/gateway'
import type { BaseRoute, Route } from './types'

const baseRoutes: BaseRoute[] = [
  {
    name: 'Your Apps',
    icon: Apps,
    path: '/',
  },
  {
    name: 'Admin Apps',
    icon: Settings,
    path: '/admin-apps',
    visible: (apps) => apps.some((app) => isAdminGateway(app.gateway)),
  },
]

// Convert the path to a URL-friendly format
const createPath = (name: string) => `/${name.replace(/\s+/g, '-').toLowerCase()}`

// Convert the base routes to routes
export const routes: Route[] = baseRoutes.map(({ name, children, path, ...rest }) => ({
  ...rest,
  name,
  path: path || createPath(name),
  children: children?.map((name) => ({ name, path: createPath(name) })),
}))
