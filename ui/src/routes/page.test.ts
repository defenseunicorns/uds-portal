// Copyright 2026 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

import { cleanup, render, screen } from '@testing-library/svelte'
import type { ApiApp } from '$lib/types'
import { afterEach, describe, expect, it } from 'vitest'

import Page from './+page.svelte'

const makeApps = (count: number): ApiApp[] =>
  Array.from({ length: count }, (_, i) => ({ name: `App ${i}`, url: `app-${i}.dev` }))

const makeData = (count: number) => ({
  userData: { name: '', username: '' },
  apps: makeApps(count),
  adminAppsEnabled: false,
})

const render39 = () => render(Page, { props: { data: makeData(39) } })
const render40 = () => render(Page, { props: { data: makeData(40) } })

describe('Your Apps page — search bar threshold', () => {
  afterEach(cleanup)

  it('hides search bar with 39 tenant apps', () => {
    render39()
    expect(screen.queryByPlaceholderText('Search')).not.toBeInTheDocument()
  })

  it('shows search bar with 40 tenant apps', () => {
    render40()
    expect(screen.getByPlaceholderText('Search')).toBeInTheDocument()
  })
})
