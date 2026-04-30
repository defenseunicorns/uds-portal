// Copyright 2025 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

import { cleanup, fireEvent, render, screen } from '@testing-library/svelte'
import { UserMenu } from '$features/navigation'
import type { UserData } from '$features/navigation/types'
import { afterEach, beforeEach, describe, expect, vi } from 'vitest'

describe('UserMenu', () => {
  const mockUserData: UserData = {
    name: 'Doug Unicorn',
    username: 'doug@uds.dev',
  }

  const assignMock = vi.fn()

  beforeEach(() => {
    // Reset all mocks before each test
    vi.clearAllMocks()

    Object.defineProperty(window, 'location', {
      configurable: true,
      value: {
        ...window.location,
        assign: assignMock,
      },
    })
  })

  afterEach(() => {
    cleanup()
  })

  test('renders with initial closed state', () => {
    render(UserMenu, { props: { userData: mockUserData } })

    // Check if the main button is rendered with user name
    expect(screen.getByText(mockUserData.username)).toBeInTheDocument()

    // Verify dropdown is not visible initially
    expect(screen.queryByText('Sign Out')).not.toBeInTheDocument()
  })

  test('opens dropdown when clicked', async () => {
    render(UserMenu, { props: { userData: mockUserData } })

    const button = screen.getByText(mockUserData.username)
    await fireEvent.click(button)

    // Check if dropdown content is visible
    expect(screen.getByText('Sign Out')).toBeInTheDocument()
    expect(screen.getByText(mockUserData.name)).toBeInTheDocument()
  })

  test('navigates to logout page when signing out', async () => {
    render(UserMenu, { props: { userData: mockUserData } })

    // Open dropdown
    const button = screen.getByText(mockUserData.username)
    await fireEvent.click(button)

    // Click sign out
    const signOutButton = screen.getByText('Sign Out')
    await fireEvent.click(signOutButton)

    // Verify navigation
    expect(assignMock).toHaveBeenCalledWith('/logout')
  })

  test('closes dropdown when clicking outside', async () => {
    render(UserMenu, { props: { userData: mockUserData } })

    // Open dropdown
    const button = screen.getByText(mockUserData.username)
    await fireEvent.click(button)

    // Click outside the dropdown
    await fireEvent.pointerDown(document.body)

    // Verify dropdown closes
    expect(screen.queryByText('Sign Out')).not.toBeInTheDocument()
  })

  test('keeps dropdown open when clicking inside', async () => {
    render(UserMenu, { props: { userData: mockUserData } })

    // Open dropdown
    const button = screen.getByText(mockUserData.username)
    await fireEvent.click(button)

    // Click inside the dropdown content
    const dropdownName = screen.getByText(mockUserData.name)
    await fireEvent.pointerDown(dropdownName)

    // Verify dropdown remains open
    expect(screen.getByText('Sign Out')).toBeInTheDocument()
  })
})
