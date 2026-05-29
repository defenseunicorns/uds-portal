// Copyright 2025-2026 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

import { cleanup, render, screen } from '@testing-library/svelte'
import type { ApiApp } from '$lib/types'
import { afterEach, describe, expect, it } from 'vitest'

import AppGrid from './component.svelte'

const appWithIcon: ApiApp = {
  name: 'Grafana',
  url: 'grafana.uds.dev',
  icon: 'data:image/svg+xml;base64,abc123',
}
const appWithoutIcon: ApiApp = { name: 'Podinfo', url: 'podinfo.uds.dev' }

describe('AppGrid', () => {
  afterEach(cleanup)

  it('renders app name', () => {
    render(AppGrid, { props: { apps: [appWithIcon] } })
    expect(screen.getByText('Grafana')).toBeInTheDocument()
  })

  it('renders img when app has icon', () => {
    const { container } = render(AppGrid, { props: { apps: [appWithIcon] } })
    const img = container.querySelector('img')
    expect(img).toBeInTheDocument()
    expect(img).toHaveAttribute('src', appWithIcon.icon)
  })

  it('renders fallback when app has no icon', () => {
    const { container } = render(AppGrid, { props: { apps: [appWithoutIcon] } })
    const fallback = container.querySelector('img')
    expect(fallback).toBeInTheDocument()
    expect(fallback).toHaveAttribute('src', '/default_logo.svg')
    expect(screen.getByText('Podinfo')).toBeInTheDocument()
  })

  it('renders multiple apps', () => {
    render(AppGrid, { props: { apps: [appWithIcon, appWithoutIcon] } })
    expect(screen.getByText('Grafana')).toBeInTheDocument()
    expect(screen.getByText('Podinfo')).toBeInTheDocument()
  })

  it('renders app as external link', () => {
    render(AppGrid, { props: { apps: [appWithIcon] } })
    const link = screen.getByRole('link', { name: /grafana/i })
    expect(link).toHaveAttribute('href', `https://${appWithIcon.url}`)
    expect(link).toHaveAttribute('target', '_blank')
    expect(link).toHaveAttribute('rel', 'noopener noreferrer')
  })
})
