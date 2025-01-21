// Copyright 2024 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

import { fireEvent, render, screen } from '@testing-library/svelte'
import type { UDSPackageStatus } from '$features/k8s/types'

import '@testing-library/jest-dom'

import { vi } from 'vitest'

import Component from './component.svelte'

const mockPackages = [
  {
    metadata: { name: 'App1', namespace: 'App1' },
    status: { endpoints: ['app1.example.com'], phase: 'Ready' as UDSPackageStatus },
  },
  {
    metadata: { name: 'App2', namespace: 'App2' },
    status: { endpoints: ['app2.example.com'], phase: 'Ready' as UDSPackageStatus },
  },
]

describe('ClusterOverviewWidgets/Applications', () => {
  test('renders correctly with no packages', () => {
    render(Component, { udsPackages: [] })

    expect(screen.getByText('No Application Packages running')).toBeInTheDocument()
  })

  test('renders correctly with packages', () => {
    render(Component, { udsPackages: mockPackages })

    expect(screen.getByText('App1')).toBeInTheDocument()
    expect(screen.getByText('App2')).toBeInTheDocument()
  })

  test('opens correct link when "Open App" button is clicked', async () => {
    window.open = vi.fn()
    render(Component, { udsPackages: mockPackages })

    const openLink = screen.getByTestId('app-widget-App1')
    await fireEvent.click(openLink)

    expect(window.open).toHaveBeenCalledWith('https://app1.example.com', '_blank')
  })

  test('renders the correct number of applications', () => {
    render(Component, { udsPackages: mockPackages })

    const appElements = screen.getAllByRole('button', { name: 'open-link-button' })

    expect(appElements.length).toBe(mockPackages.length)
  })

  test('handles packages with no endpoints', () => {
    const packagesWithEmptyEndpoint = [
      {
        metadata: { name: 'App3', namespace: 'App3' },
        status: { endpoints: [], phase: 'Pending' as UDSPackageStatus },
      },
    ]

    render(Component, { udsPackages: packagesWithEmptyEndpoint })

    expect(screen.queryByText('App3')).not.toBeInTheDocument()
  })
})
