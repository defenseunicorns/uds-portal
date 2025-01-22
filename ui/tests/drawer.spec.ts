// Copyright 2025 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

import { expect, test } from '@playwright/test'

test.describe('Drawer', async () => {
  test.describe('Pods', async () => {
    test.beforeEach(async ({ page }) => {
      await page.goto('/workloads/pods')

      await page.getByTestId(/^podinfo-.*-testid-1$/).click()
    })

    test.describe('is opened when clicking on a table row and', async () => {
      test('will display Details', async ({ page }) => {
        const drawerEl = page.getByTestId('drawer')

        await expect(drawerEl).toBeVisible()
        await expect(drawerEl.getByText('Created')).toBeVisible()
        await expect(drawerEl.getByText('Name', { exact: true })).toBeVisible()
        await expect(drawerEl.getByText('Annotations')).toBeVisible()
        await expect(drawerEl.getByText('podinfo', { exact: true }).first()).toBeVisible()
        await expect(drawerEl.getByText('View Logs')).toBeVisible()
      })

      test('will display Events details', async ({ page }) => {
        const drawerEl = page.getByTestId('drawer')

        await expect(drawerEl).toBeVisible()
        await drawerEl.getByRole('button', { name: 'Events' }).click()

        await expect(drawerEl.getByText('Created container podinfo')).toBeVisible()
      })

      test('will display details', async ({ page }) => {
        const drawerEl = page.getByTestId('drawer')
        const labelKey = await drawerEl.locator('span:text("app.kubernetes.io/name:")').textContent()

        // get sibling span as label value
        const labelValue = await drawerEl
          .locator('span:text("app.kubernetes.io/name:")')
          .locator('xpath=following-sibling::span')
          .textContent()

        const podID = page.url().split('/').pop()
        await drawerEl.getByRole('button', { name: 'YAML' }).click()
        await expect(drawerEl.getByText('namespace:')).toBeVisible()

        // ensure the label key and value are correct
        expect(labelKey!.replace(/\s/g, '')).toEqual(`app.kubernetes.io/name:`)
        expect(labelValue!.replace(/\s/g, '')).toEqual(`podinfo`)
        await expect(drawerEl.locator(`:text("uid: ${podID}")`)).toBeVisible()
      })
    })
  })

  test.describe('Non pods', async () => {
    test('will not display logs button for non-pods', async ({ page }) => {
      await page.goto('/workloads/deployments')
      await page.getByTestId(/^coredns-testid-1$/).click()
      const drawerEl = page.getByTestId('drawer')
      await expect(drawerEl).toBeVisible()
      await expect(drawerEl.getByTestId('view-logs-btn')).not.toBeVisible()
    })
  })
})
