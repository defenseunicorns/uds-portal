// Copyright 2026 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

import { expect, test } from '@playwright/test'

test.describe('endpoint filtering', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/')
  })

  test('annotation-hidden endpoint does not appear as app tile', async ({ page }) => {
    await expect(page.locator('a[href="https://multi-expose-package.uds.dev"]')).toBeVisible()
    await expect(page.locator('a[href="https://multi-expose-package-hidden.uds.dev"]')).toHaveCount(0)
  })

  test('wildcard endpoint does not appear as app tile', async ({ page }) => {
    await expect(page.locator('a[href*="multi-expose-package-wildcard"]')).toHaveCount(0)
  })
})
