// Copyright 2024 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

import { expect, test } from '@playwright/test'

// NOTE: these tests should only be run against the in-cluster environment
// check global-setup.ts for the login process
test.describe('in-cluster', async () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/')
  })

  test.describe('home page', async () => {
    test('user menu expands and collapses when clicked on', async ({ page }) => {
      const userMenuButton = page.getByText('doug', { exact: true })

      await userMenuButton.click()
      const signOutButton = page.getByText('Sign Out')

      await expect(signOutButton).toBeVisible()

      await userMenuButton.click()
      await expect(signOutButton).not.toBeVisible()
    })

    test('user menu expands and collapses when clicked away from', async ({ page }) => {
      const userMenuButton = page.getByText('doug', { exact: true })

      await userMenuButton.click()
      const signOutButton = page.getByText('Sign Out')

      await expect(signOutButton).toBeVisible()

      await page.getByRole('link', { name: 'Overview' }).click()
      await expect(signOutButton).not.toBeVisible()
    })

    test('classification banners display', async ({ page }) => {
      const classificationBanner = page.getByTestId('classification-header')
      await expect(classificationBanner).toBeVisible()

      const classificationFooter = page.getByTestId('classification-footer')
      await expect(classificationFooter).toBeVisible()
      await expect(page.getByText('UNCLASSIFIED').first()).toBeVisible()
    })
  })

  // we perform these tests in-cluster because the tables aren't populated in the "runtime" test cluster
  test.describe('Applications pages', async () => {
    test('application packages page', async ({ page }) => {
      await page.getByRole('button', { name: 'Applications' }).click()
      await page.getByRole('link', { name: /^Packages$/ }).click()

      // click on row
      await page.getByText('init').click()

      // ensure details are visible
      await expect(page.getByText('Architecture')).toBeVisible()

      // use exact: true to ensure it isn't the text in the table
      await expect(page.getByText('zarf-injector', { exact: true })).toBeVisible()
    })

    test('applications endpoints page', async ({ page }) => {
      await page.getByRole('button', { name: 'Applications' }).click()
      await page.getByRole('link', { name: /^Endpoints$/ }).click()

      await expect(page.getByText('sso.uds.dev')).toBeVisible()
    })
  })

  test.describe('cves page', async () => {
    test('cves page loads', async ({ page }) => {
      await page.getByRole('button', { name: 'Security' }).click()
      await page.getByRole('link', { name: 'CVE Report' }).click()
      await expect(page.getByTestId(/^cve-report(?!-search).*/).first()).toBeVisible()

      // filter by severity
      await page.getByRole('radio', { name: 'Critical' }).click()
      let cveReports = await page.getByTestId(/^cve-report(?!-search).*/).all()

      for (const report of cveReports) {
        const cvssScore = await report.getByTestId('cvss-score').innerText()
        expect(cvssScore).toContain('critical')
      }

      // reset filter
      await page.getByRole('radio', { name: 'All' }).click()

      // filter by search
      await page.getByTestId('cve-report-search').fill('GHSA')
      cveReports = await page.getByTestId(/^cve-report(?!-search).*/).all()

      for (const report of cveReports) {
        const cveId = await report.getByTestId('cve_id').innerText()
        expect(cveId).not.toContain('CVE')
      }
    })

    test.describe('CVE Drawer', async () => {
      test.beforeEach(async ({ page }) => {
        await page.goto('/security/cve-report')

        await page
          .getByTestId(/^cve-report(?!-search).*/)
          .first()
          .click()
      })

      test('will display CVE details', async ({ page }) => {
        const drawerEl = page.getByTestId('drawer')

        await expect(drawerEl.getByText(/CVSS Score/)).toBeVisible()
        await expect(drawerEl.getByText(/Status/).first()).toBeVisible()
        await expect(drawerEl.getByText(/Description/)).toBeVisible()
        await expect(drawerEl.getByText(/Build Date/)).toBeVisible()
      })
    })
  })
})
