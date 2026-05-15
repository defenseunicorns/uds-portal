// Copyright 2026 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

import { expect, test } from '@playwright/test'

test.describe('Admin Apps view', () => {
  test('Admin Apps sidebar link appears when an admin app exists', async ({ page }) => {
    await page.goto('/')
    await expect(page.getByRole('link', { name: 'Admin Apps' })).toBeVisible()
  })

  test('/admin-apps renders banner and admin app', async ({ page }) => {
    await page.goto('/admin-apps')
    await expect(page.getByTestId('admin-access-banner')).toContainText('additional network access')
    await expect(page.getByText('Admin Apps').first()).toBeVisible()
    await expect(page.getByRole('link', { name: 'Podinfo Admin' })).toBeVisible()
  })

  test('/ does not show the admin gateway app', async ({ page }) => {
    await page.goto('/')
    await expect(page.getByRole('link', { name: 'Podinfo Admin' })).toHaveCount(0)
  })
})
