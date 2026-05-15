// Copyright 2026 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

import { readable } from 'svelte/store'

import { cleanup, render, screen } from '@testing-library/svelte'
import { afterEach, describe, expect, it, vi } from 'vitest'

import type { ApiApp } from '../../../../routes/types'
import Sidebar from './component.svelte'

vi.mock('$app/stores', () => ({
  page: readable({ url: new URL('http://localhost/') }),
}))

vi.mock('$app/paths', () => ({
  resolve: (p: string) => p,
}))

const tenantApp: ApiApp = { name: 'Tenant', url: 'tenant.uds.dev', gateway: 'tenant' }
const adminApp: ApiApp = { name: 'Grafana', url: 'grafana.admin.uds.dev', gateway: 'admin' }
const customAdminApp: ApiApp = { name: 'Custom', url: 'c.uds.dev', gateway: 'custom-admin-gw' }

describe('Sidebar', () => {
  afterEach(cleanup)

  it('always renders Your Apps link', () => {
    render(Sidebar, { props: { apps: [], adminAppsEnabled: true } })
    expect(screen.getByText('Your Apps')).toBeInTheDocument()
  })

  it('hides Admin Apps when no admin gateway apps exist', () => {
    render(Sidebar, { props: { apps: [tenantApp], adminAppsEnabled: true } })
    expect(screen.queryByText('Admin Apps')).toBeNull()
  })

  it('shows Admin Apps when an admin gateway app exists', () => {
    render(Sidebar, { props: { apps: [adminApp], adminAppsEnabled: true } })
    expect(screen.getByText('Admin Apps')).toBeInTheDocument()
  })

  it('shows Admin Apps for custom admin-named gateway', () => {
    render(Sidebar, { props: { apps: [customAdminApp], adminAppsEnabled: true } })
    expect(screen.getByText('Admin Apps')).toBeInTheDocument()
  })

  it('hides Admin Apps when adminAppsEnabled is false even if admin apps exist', () => {
    render(Sidebar, { props: { apps: [adminApp], adminAppsEnabled: false } })
    expect(screen.queryByText('Admin Apps')).toBeNull()
  })

  it('hides Admin Apps when adminAppsEnabled is false with custom admin gateway', () => {
    render(Sidebar, { props: { apps: [customAdminApp], adminAppsEnabled: false } })
    expect(screen.queryByText('Admin Apps')).toBeNull()
  })
})
