// Copyright 2025 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

import fs from 'node:fs'

import { expect, test } from '@playwright/test'

test.describe('Logs', async () => {
  // open logs view and load logs
  test.beforeEach(async ({ page }) => {
    await page.goto('/workloads/pods')

    // click on pod row
    await page.getByTestId(/^podinfo-.*-testid-1$/).click()

    // click on logs button
    const logsButton = page.getByTestId('view-logs-btn')
    await expect(logsButton).toBeVisible()
    await logsButton.click()

    // ensure the podinfo pod row is active/highlighted
    const podRowContent = page.getByTestId(/^podinfo-.*-testid-1$/)
    const podRow = podRowContent.locator('xpath=../..') // get grandparent element
    const podRowClasses = await podRow.getAttribute('class')
    expect(podRowClasses).toContain('active')

    // ensure "no container selected" message
    await page.getByText('No container selected')

    // select podinfo container
    await page.locator('#container-dropdown').click()
    await page.locator('text=podinfo').last().click()

    // wait for logs to load
    await page.waitForSelector('.log-line')
  })

  test('autoload logs for pods with a single container', async ({ page }) => {
    // click on a pod with a single container (local-path-provisioner)
    await page.getByTestId(/^local-path-provisioner-.*-testid-1$/).click()

    // click on logs button
    const logsButton = page.getByTestId('view-logs-btn')
    await expect(logsButton).toBeVisible()
    await logsButton.click()

    // wait for logs to autoload
    await page.waitForSelector('.log-line')
  })

  test('arrows keys do not break logs view', async ({ page }) => {
    await page.keyboard.press('ArrowUp')

    // assert that logs are still visible and the app didn't throw a 404
    let logLines = await page.locator('.log-line').count()
    expect(logLines).toBeGreaterThan(0)

    await page.keyboard.press('ArrowDown')
    logLines = await page.locator('.log-line').count()
    expect(logLines).toBeGreaterThan(0)
  })

  test('changing namespace, pod, and container to update logs', async ({ page }) => {
    // assert logs are visible
    let logLines = await page.locator('.log-line').count()
    expect(logLines).toBeGreaterThan(0)

    // get Pepr pod names so we can easily select them later
    const peprPods = await page.locator('text=pepr-uds-core-').all()
    const peprPodNames = await Promise.all(peprPods.map(async (el) => el.textContent()))
    peprPodNames.sort()

    // click on namespace dropdown
    await page.locator('#ns-dropdown').click()
    await page.locator('text=pepr-system').last().click()

    await page.waitForTimeout(200)

    // wait for logs to load
    logLines = await page.locator('.log-line').count()
    expect(logLines).toBe(0)

    // select a new pod
    await page.locator('#pod-dropdown').click()
    await page.locator(`text=${peprPodNames[0]}`).last().click()
    logLines = await page.locator('.log-line').count()
    expect(logLines).toBe(0)

    // select a new container and wait for new logs to load
    await page.locator('#container-dropdown').click()
    await page.locator('text=server').last().click()
    await page.waitForTimeout(500)
    logLines = await page.locator('.log-line').count()
    expect(logLines).toBeGreaterThan(0)

    // select another pod from the same namespace
    await page.locator('#pod-dropdown').click()
    await page
      .locator(`text=${peprPodNames.at(-1)}`)
      .last()
      .click() // should be watcher pod because we sorted the names
    logLines = await page.locator('.log-line').count()
    expect(logLines).toBe(0)

    // select container and wait for new logs to load
    await page.locator('#container-dropdown').click()
    await page.locator('text=watcher').last().click()
    await page.waitForTimeout(500)
    logLines = await page.locator('.log-line').count()
    expect(logLines).toBeGreaterThan(0)
  })

  test('highlighted search term matches log and resets when empty', async ({ page }) => {
    const searchTerm = 'Starting HTTP Server'
    await page.getByTestId('logs-search').fill(searchTerm)

    await page.waitForSelector('.log-line')
    const logLines = page.getByText(searchTerm)
    await expect(logLines).toBeVisible()

    // expect a mark tag on the log line that we searched for
    const highlighted = page.locator('div.log-line mark').first()
    expect(await highlighted.textContent()).toBe(searchTerm)

    // reset search term and assert that the mark tag is gone
    await page.getByTestId('logs-search').fill('')
    await page.waitForSelector('.log-line')
    const numHighlights = page.locator('div.log-line mark').count()
    expect(await numHighlights).toBe(0)
  })

  test('copy logs to clipboard', async ({ page, context }) => {
    await context.grantPermissions(['clipboard-read', 'clipboard-write']) // need this to copy to clipboard
    const copyButton = page.getByTestId('copy-logs-btn')
    await expect(copyButton).toBeVisible()
    await copyButton.click()

    const notification = await page.getByTestId('logs-notification')
    await expect(notification).toHaveText('Copied!')

    // assert on clipboard contents
    const clipboardText = await page.evaluate(() => navigator.clipboard.readText())
    expect(clipboardText).toContain('Starting HTTP Server')
  })

  test('download logs', async ({ page }) => {
    const downloadButton = page.getByTestId('download-logs-btn')
    await expect(downloadButton).toBeVisible()
    const [download] = await Promise.all([page.waitForEvent('download'), downloadButton.click()])

    const notification = await page.getByTestId('logs-notification')
    await expect(notification).toHaveText('Downloading')

    // assert on downloaded file contents
    const path = await download.path()
    await download.saveAs(path)
    fs.readFile(path, 'utf8', (err, data) => {
      if (err) {
        console.error('Error reading the file:', err)
        return
      }
      // assert on file contents
      expect(data).toContain('Starting HTTP Server')
    })
  })

  test('resizing logs view', async ({ page }) => {
    const podTable = await page.$('.table-content')
    const logs = await page.$('.logs')
    const resizeBar = await page.$('.resize-bar')

    if (!podTable || !logs || !resizeBar) {
      throw new Error('Required elements not found')
    }

    const initialPodTableHeight = await podTable.evaluate((node) => getComputedStyle(node).height)
    const initialLogsHeight = await logs.evaluate((node) => getComputedStyle(node).height)

    // Simulate drag event
    const box = await resizeBar.boundingBox()
    if (!box) {
      throw new Error('Bounding box not found')
    }
    await page.mouse.move(box.x + box.width / 2, box.y + box.height / 2)
    await page.mouse.down()
    await page.mouse.move(box.x + box.width / 2, box.y + box.height / 2 + 50)
    await page.mouse.up()

    const newPodTableHeight = await podTable.evaluate((node) => getComputedStyle(node).height)
    const newLogsHeight = await logs.evaluate((node) => getComputedStyle(node).height)

    expect(newPodTableHeight).not.toBe(initialPodTableHeight)
    expect(newLogsHeight).not.toBe(initialLogsHeight)
  })

  test('escape key closes logs view', async ({ page }) => {
    await expect(page.getByTestId('logs-view')).toBeVisible()
    await page.keyboard.press('Escape')
    await expect(page.getByTestId('logs-view')).not.toBeVisible()
  })

  test('close button closes logs view', async ({ page }) => {
    await expect(page.getByTestId('logs-view')).toBeVisible()
    const closeButton = page.getByTestId('close-logs-btn')
    await expect(closeButton).toBeVisible()
    await closeButton.click()
    await expect(page.getByTestId('logs-view')).not.toBeVisible()
  })
})
