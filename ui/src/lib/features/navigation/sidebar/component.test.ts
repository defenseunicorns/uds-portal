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

const tenantApp: ApiApp = { name: 'Tenant', url: 'tenant.uds.dev', gateway: 'tenant', group: 0 }
const adminApp: ApiApp = { name: 'Grafana', url: 'grafana.admin.uds.dev', gateway: 'admin', group: 0 }
const customAdminApp: ApiApp = { name: 'Custom', url: 'c.uds.dev', gateway: 'custom-admin-gw', group: 0 }

describe('Sidebar', () => {
  afterEach(cleanup)

  it('does not render sidebar when no apps exist', () => {
    render(Sidebar, { props: { apps: [], adminAppsEnabled: true } })
    expect(screen.queryByTestId('sidebar')).toBeNull()
  })

  it('does not render sidebar when only tenant apps exist', () => {
    render(Sidebar, { props: { apps: [tenantApp], adminAppsEnabled: true } })
    expect(screen.queryByTestId('sidebar')).toBeNull()
  })

  it('renders sidebar with both links when an admin gateway app exists', () => {
    render(Sidebar, { props: { apps: [adminApp], adminAppsEnabled: true } })
    expect(screen.getByTestId('sidebar')).toBeInTheDocument()
    expect(screen.getByText('Your Apps')).toBeInTheDocument()
    expect(screen.getByText('Admin Apps')).toBeInTheDocument()
  })

  it('renders sidebar with both links for custom admin-named gateway', () => {
    render(Sidebar, { props: { apps: [customAdminApp], adminAppsEnabled: true } })
    expect(screen.getByTestId('sidebar')).toBeInTheDocument()
    expect(screen.getByText('Your Apps')).toBeInTheDocument()
    expect(screen.getByText('Admin Apps')).toBeInTheDocument()
  })

  it('does not render sidebar when adminAppsEnabled is false even if admin apps exist', () => {
    render(Sidebar, { props: { apps: [adminApp], adminAppsEnabled: false } })
    expect(screen.queryByTestId('sidebar')).toBeNull()
  })

  it('does not render sidebar when adminAppsEnabled is false with custom admin gateway', () => {
    render(Sidebar, { props: { apps: [customAdminApp], adminAppsEnabled: false } })
    expect(screen.queryByTestId('sidebar')).toBeNull()
  })
})
