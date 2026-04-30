// Copyright 2025 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

import { expect, test } from '@playwright/test'

test.describe('portal', async () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/')
  })

  test.describe('home page', async () => {
    test('user menu expands and collapses when clicked', async ({ page }) => {
      const userMenuButton = page.getByTestId('user-menu-button')
      await userMenuButton.click()
      await expect(page.getByText('Sign Out')).toBeVisible()
      await userMenuButton.click()
      await expect(page.getByText('Sign Out')).not.toBeVisible()
    })

    test('user menu closes when clicking away', async ({ page }) => {
      await page.getByTestId('user-menu-button').click()
      await expect(page.getByText('Sign Out')).toBeVisible()
      await page.getByText('Your Apps').click()
      await expect(page.getByText('Sign Out')).not.toBeVisible()
    })
  })

  test.describe('app grid', async () => {
    test('page loads with app grid heading', async ({ page }) => {
      await expect(page.getByText('Your Apps')).toBeVisible()
    })

    test('My Account is the first app', async ({ page }) => {
      const firstApp = page.locator('a[target="_blank"]').first()
      await expect(firstApp).toContainText('My Account')
    })

    test('Podinfo is in the app grid', async ({ page }) => {
      await expect(page.getByRole('link', { name: 'Podinfo' })).toBeVisible()
    })

    test('clicking Podinfo opens the podinfo page', async ({ page, context }) => {
      const [newPage] = await Promise.all([
        context.waitForEvent('page'),
        page.getByRole('link', { name: 'Podinfo' }).click(),
      ])
      await newPage.waitForLoadState('networkidle')
      await expect(newPage.getByRole('button', { name: 'Ping' })).toBeVisible()
    })
  })
})
