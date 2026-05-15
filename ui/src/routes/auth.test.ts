// Copyright 2025-2026 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

import { authenticated } from '$features/auth/store'
import { beforeEach, describe, expect, vi } from 'vitest'

import { load } from './+layout'

// Mock stores
vi.mock('$features/auth/store', () => ({
  authenticated: {
    set: vi.fn(),
    subscribe: vi.fn(() => () => {}),
    update: vi.fn(),
  },
}))

describe('load function', () => {
  // Mock fetch
  const fetchMock = vi.fn()
  global.fetch = fetchMock

  beforeEach(() => {
    vi.clearAllMocks()
    ;(
      window as Window & { __APP__?: { CLASSIFICATION_BANNER?: { enabled: boolean; text: string; footer: boolean } } }
    ).__APP__ = {
      CLASSIFICATION_BANNER: {
        enabled: false,
        text: '',
        footer: false,
      },
    }
  })

  test('successful local authentication with token in URL', async () => {
    const mockUserData = {
      name: 'First Last',
      username: 'local',
    }
    fetchMock.mockResolvedValueOnce({
      ok: true,
      json: () => mockUserData,
    })
    fetchMock.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve([]),
    })

    const result = await load()

    // Verify fetch was called correctly
    expect(fetchMock).toHaveBeenNthCalledWith(1, '/api/v1/auth', {
      method: 'GET',
      headers: new Headers({
        'Content-Type': 'application/json',
      }),
    })

    // Verify store operations
    expect(authenticated.set).toHaveBeenCalledWith(true)

    // Verify return value
    expect(result).toEqual({
      userData: {
        name: 'First Last',
        username: 'local',
      },
      apps: [],
      adminAppsEnabled: true,
    })
  })

  test('successful in-cluster authentication (without token in URL)', async () => {
    const mockUserData = {
      name: 'Doug Unicorn',
      username: 'doug@defenseunicorns.com',
    }
    fetchMock.mockResolvedValueOnce({
      ok: true,
      json: () => mockUserData,
    })
    fetchMock.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve([]),
    })

    const result = await load()

    // Verify fetch was called with empty token
    expect(fetchMock).toHaveBeenNthCalledWith(1, '/api/v1/auth', {
      method: 'GET',
      headers: new Headers({
        'Content-Type': 'application/json',
      }),
    })

    // Verify authenticated state was set to true
    expect(authenticated.set).toHaveBeenCalledWith(true)

    // Verify return value
    expect(result).toEqual({
      userData: {
        name: 'Doug Unicorn',
        username: 'doug@defenseunicorns.com',
      },
      apps: [],
      adminAppsEnabled: true,
    })
  })

  test('authentication failure with invalid token', async () => {
    fetchMock.mockResolvedValueOnce({
      ok: false,
      status: 401,
      statusText: 'Unauthorized',
    })
    fetchMock.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve([]),
    })

    const result = await load()

    // Verify authenticated state was set to false
    expect(authenticated.set).toHaveBeenCalledWith(false)

    // Verify return value
    expect(result).toEqual({
      userData: {
        name: '',
        username: '',
      },
      apps: [],
      adminAppsEnabled: true,
    })
  })

  test('network errors during authentication', async () => {
    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {})
    const networkError = new Error('Network error')

    fetchMock.mockRejectedValueOnce(networkError)
    fetchMock.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve([]),
    })

    const result = await load()

    // Verify error was logged
    expect(consoleSpy).toHaveBeenCalledWith('Load error:', expect.any(Error))

    // Verify authenticated state was set to false
    expect(authenticated.set).toHaveBeenCalledWith(false)

    // Verify return value
    expect(result).toEqual({
      userData: {
        name: '',
        username: '',
      },
      apps: [],
      adminAppsEnabled: true,
    })

    consoleSpy.mockRestore()
  })

  test('malformed JSON in successful response', async () => {
    fetchMock.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.reject(new Error('Invalid JSON')),
    })
    fetchMock.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve([]),
    })

    const result = await load()

    // Verify authenticated state was set to false
    expect(authenticated.set).toHaveBeenCalledWith(false)

    // Verify return value
    expect(result).toEqual({
      userData: {
        name: '',
        username: '',
      },
      apps: [],
      adminAppsEnabled: true,
    })
  })
})
