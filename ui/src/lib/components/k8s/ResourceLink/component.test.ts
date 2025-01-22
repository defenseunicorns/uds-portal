// Copyright 2025 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

import type { KubernetesObject } from '@kubernetes/client-node'
import { fireEvent, render, waitFor } from '@testing-library/svelte'
import { goto } from '$app/navigation'
import * as toastStore from '$features/toast/store'
import { vi, type MockInstance } from 'vitest'

import ResourceLink from './component.svelte'

vi.mock('$app/navigation', () => ({ goto: vi.fn() }))

describe('ResourceLink Component', () => {
  let addToast: MockInstance
  beforeEach(() => {
    addToast = vi.spyOn(toastStore, 'addToast').mockImplementation(() => {})
  })

  afterEach(() => {
    vi.clearAllMocks()
  })

  it('navigates to a resource when resources are provided', async () => {
    const resources = [
      {
        metadata: { name: 'test', uid: '12345' },
      },
    ]
    const { getByText } = render(ResourceLink, {
      props: {
        name: 'test',
        route: 'resource',
        resources: resources as KubernetesObject[],
      },
    })

    const button = getByText('test')
    await fireEvent.click(button)

    expect(goto).toHaveBeenCalledWith('/resource/12345')
  })

  it('fetches resources if resources not provided, then navigates to resource route', async () => {
    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => [
        {
          metadata: { name: 'test', uid: '12345' },
        },
      ],
    })

    const { getByText } = render(ResourceLink, {
      props: { name: 'test', route: 'resource', resourcePath: '/api/resources' },
    })

    const button = getByText('test')
    await fireEvent.click(button)

    expect(fetch).toHaveBeenCalledWith('/api/resources?once=true&fields=.metadata')
    await waitFor(() => expect(goto).toHaveBeenCalledWith('/resource/12345'))
  })

  it('does not navigate to target if no match is found in resources', async () => {
    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => [
        {
          metadata: { name: 'test', uid: '12345' },
        },
      ],
    })

    const { getByText } = render(ResourceLink, {
      props: { name: 'not-found', route: 'resource', resourcePath: '/api/resources' },
    })

    const button = getByText('not-found')
    fireEvent.click(button)

    await waitFor(() => expect(goto).not.toHaveBeenCalled())
  })

  it('shows a toast message on fetch error', async () => {
    global.fetch = vi.fn().mockResolvedValue({ ok: false, text: async () => 'Error fetching resources' })

    const { getByText } = render(ResourceLink, {
      props: { name: 'test', route: 'resource', resourcePath: '/api/resources' },
    })

    const button = getByText('test')
    fireEvent.click(button)

    await waitFor(() => {
      expect(goto).not.toHaveBeenCalled()
      expect(addToast).toHaveBeenCalledWith({
        timeoutSecs: 5,
        message: `Failed to fetch test resource`,
        type: 'error',
      })
    })
  })
})
