// Copyright 2026 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

import { tick } from 'svelte'

import { cleanup, render, screen } from '@testing-library/svelte'
import { afterEach, describe, expect, it, vi } from 'vitest'

import type { ApiApp } from '../../../../routes/types'
import Sidebar from './component.svelte'

const navigation = vi.hoisted(() => ({
  callbacks: [] as Array<(navigation: { to: { url: URL } | null }) => void>,
}))

vi.mock('$app/navigation', () => ({
  afterNavigate: (callback: (navigation: { to: { url: URL } | null }) => void) => {
    navigation.callbacks.push(callback)
  },
}))

vi.mock('$app/paths', () => ({
  resolve: (p: string) => p,
}))

const tenantApp: ApiApp = { name: 'Tenant', url: 'tenant.uds.dev', gateway: 'tenant' }
const adminApp: ApiApp = { name: 'Grafana', url: 'grafana.admin.uds.dev', gateway: 'admin' }
const customAdminApp: ApiApp = { name: 'Custom', url: 'c.uds.dev', gateway: 'custom-admin-gw' }

describe('Sidebar', () => {
  afterEach(() => {
    cleanup()
    navigation.callbacks.length = 0
  })

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

  it('updates the active link when navigating between routes', async () => {
    window.history.replaceState({}, '', '/')
    render(Sidebar, { props: { apps: [adminApp], adminAppsEnabled: true } })

    const appsLink = screen.getByRole('link', { name: 'Your Apps' })
    const adminLink = screen.getByRole('link', { name: 'Admin Apps' })
    expect(appsLink).toHaveAttribute('aria-current', 'page')
    expect(adminLink).not.toHaveAttribute('aria-current')

    navigation.callbacks[0]({ to: { url: new URL('http://localhost/admin-apps') } })
    await tick()
    expect(adminLink).toHaveAttribute('aria-current', 'page')
    expect(appsLink).not.toHaveAttribute('aria-current')

    navigation.callbacks[0]({ to: { url: new URL('http://localhost/') } })
    await tick()
    expect(appsLink).toHaveAttribute('aria-current', 'page')
    expect(adminLink).not.toHaveAttribute('aria-current')
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
