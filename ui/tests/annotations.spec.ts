// Copyright 2026 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

import { expect, test } from '@playwright/test'

test.describe('per-endpoint annotations', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/')
  })

  test('endpoint with uds.dev/title shows custom title in tile', async ({ page }) => {
    await expect(page.getByRole('link', { name: 'Custom Annotated App', exact: true })).toBeVisible()
  })

  test('endpoint with uds.dev/icon shows custom icon in tile', async ({ page }) => {
    const tile = page.getByRole('link', { name: 'Custom Annotated App', exact: true })
    await expect(tile.locator('img')).toHaveAttribute(
      'src',
      'data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciLz4=',
    )
  })

  test('portal.uds.dev/title takes precedence over uds.dev/title in tile', async ({ page }) => {
    await expect(page.getByRole('link', { name: 'Portal Override Title', exact: true })).toBeVisible()
    await expect(page.getByRole('link', { name: 'UDS Title', exact: true })).toHaveCount(0)
  })

  test('endpoint with portal.uds.dev/visible=false does not appear as a tile', async ({ page }) => {
    await expect(page.getByRole('link', { name: 'Custom Annotated App', exact: true })).toBeVisible()
    await expect(page.locator('a[href="https://annotated-package-hidden.uds.dev"]')).toHaveCount(0)
  })
})
