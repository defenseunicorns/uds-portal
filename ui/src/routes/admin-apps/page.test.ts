// Copyright 2026 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

import { cleanup, render, screen } from '@testing-library/svelte'
import type { ApiApp } from '$lib/types'
import { afterEach, describe, expect, it } from 'vitest'

import Page from './+page.svelte'

const makeAdminApps = (count: number): ApiApp[] =>
  Array.from({ length: count }, (_, i) => ({ name: `Admin App ${i}`, url: `admin-${i}.dev`, gateway: 'admin' }))

const renderPage = (count: number) =>
  render(Page, {
    props: { data: { userData: { name: '', username: '' }, apps: makeAdminApps(count), adminAppsEnabled: true } },
  })

describe('Admin Apps page — search bar threshold', () => {
  afterEach(cleanup)

  it('hides search bar with 39 admin apps', () => {
    renderPage(39)
    expect(screen.queryByPlaceholderText('Search')).not.toBeInTheDocument()
  })

  it('shows search bar with 40 admin apps', () => {
    renderPage(40)
    expect(screen.getByPlaceholderText('Search')).toBeInTheDocument()
  })
})
