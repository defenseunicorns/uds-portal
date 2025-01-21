import { execSync } from 'child_process'

import { expect, test } from '@playwright/test'
import { K8s, kind } from 'kubernetes-fluent-client'

// Utility function to run shell commands
function execCommand(command: string) {
  execSync(command, { stdio: 'inherit' })
}

async function createPod() {
  await K8s(kind.Pod).Apply({
    metadata: {
      name: 'new-pod',
      namespace: 'default',
    },
    spec: {
      containers: [
        {
          name: 'my-container',
          image: 'nginx',
        },
      ],
    },
  })
}

test.describe('Cluster connection interactions -- switching, reconnecting', () => {
  test('should handle cluster disconnection and reconnection', async ({ page }) => {
    test.setTimeout(120000)
    await page.goto('/workloads/pods')

    // Stop the cluster
    execCommand('k3d cluster stop runtime-2')

    // Wait for disconnection to be detected
    await expect(page.getByText('Cluster health check failed: no connection')).toBeVisible({ timeout: 20000 })

    // Start the cluster again
    execCommand('k3d cluster start runtime-2')

    // Wait for the reconnection to be detected
    await expect(page.getByText('Cluster connection restored')).toBeVisible({ timeout: 15000 })

    // Use KFC to create a new pod
    await createPod()

    // ensure stream is using latest cache, meaning view should show new pod
    await expect(page.getByText('new-pod')).toBeVisible({ timeout: 15000 })
  })

  test('should switch between clusters', async ({ page }) => {
    await page.goto('/workloads/pods')

    // Click clustermenu dropdown
    await page.getByRole('button', { name: 'k3d-runtime-2', exact: true }).click()

    // Click on runtime since current-context is runtime-2
    await page.getByRole('button', { name: 'k3d-runtime', exact: true }).click()

    // Expect loading.svelte to be visible
    await expect(page.getByText('Connecting to cluster')).toBeVisible()

    // Wait for 3 seconds
    await page.waitForTimeout(3000)

    // Expect loading.svelte to be hidden and goto('/') which means cluster-overview
    await expect(page.getByText('Connecting to cluster')).toBeHidden()
    await expect(page).toHaveURL('/')
    await expect(page.getByRole('button', { name: 'k3d-runtime', exact: true })).toBeVisible()
  })

  test('switch after failure', async ({ page }) => {
    test.setTimeout(120000)

    await page.goto('/workloads/pods')

    await expect(page.getByRole('button', { name: 'k3d-runtime', exact: true })).toBeVisible()

    // Stop the cluster
    execCommand('k3d cluster stop runtime')

    // Wait for disconnection to be detected
    await expect(page.getByText('Cluster health check failed: no connection')).toBeVisible({ timeout: 15000 })

    // Click clustermenu dropdown
    await page.getByRole('button', { name: 'k3d-runtime', exact: true }).click()

    // Click on runtime since current-context is runtime-2
    await page.getByRole('button', { name: 'k3d-runtime-2', exact: true }).click()

    // Expect loading.svelte to be visible
    await expect(page.getByText('Connecting to cluster')).toBeVisible()

    // Wait for 3 seconds
    await page.waitForTimeout(3000)

    // Expect loading.svelte to be hidden and goto('/') which means cluster-overview
    await expect(page.getByText('Connecting to cluster')).toBeHidden()
    await expect(page.getByText('Cluster health check failed: no connection')).toBeHidden()
    await expect(page).toHaveURL('/')
    await expect(page.getByRole('button', { name: 'k3d-runtime-2', exact: true })).toBeVisible()
  })

  test('switch to failed cluster returns error', async ({ page }) => {
    test.setTimeout(120000)
    await page.goto('/')

    await expect(page.getByRole('button', { name: 'k3d-runtime-2', exact: true })).toBeVisible()

    // Click clustermenu dropdown
    await page.getByRole('button', { name: 'k3d-runtime-2', exact: true }).click()

    // Click on runtime since that cluster was stopped in previous test
    await page.getByRole('button', { name: 'k3d-runtime', exact: true }).click()

    // Expect loading.svelte to be visible
    await expect(page.getByText('connection refused')).toBeVisible()
  })
})
