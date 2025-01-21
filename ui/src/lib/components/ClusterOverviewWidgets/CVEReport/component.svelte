<!-- Copyright 2024 Defense Unicorns -->
<!-- SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial -->

<script lang="ts">
  import { onMount } from 'svelte'

  import { goto } from '$app/navigation'
  import { HeaderWithIcon } from '$components'
  import { type ReportWithVulnerability } from '$features/security'
  import { addToast } from '$features/toast'
  import { formatDetailedAge } from '$lib/utils/helpers'
  import { ChevronRight } from 'carbon-icons-svelte'
  import { format } from 'date-fns'

  let reports: ReportWithVulnerability[] = []
  let loading = true
  let scanTimeInDays: unknown = 0
  let buildDate = ''
  let criticalCount: number = 0
  let highCount: number = 0
  let mediumCount: number = 0

  onMount(async () => {
    try {
      const resp = await fetch('/api/v1/security/cves')

      if (resp.ok) {
        const results = await resp.json()

        reports = results.reports

        buildDate = format(results.buildDate, 'dd MMM yyyy')
        scanTimeInDays = formatDetailedAge(results.scanTime)

        reports.map((report) => {
          if (report.severity === 'CRITICAL') {
            criticalCount++
          } else if (report.severity === 'HIGH') {
            highCount++
          } else if (report.severity === 'MEDIUM') {
            mediumCount++
          }
        })
      } else {
        throw new Error(`${await resp.text()}`)
      }
    } catch (err) {
      console.error('Failed to fetch CVE report', err)
      addToast({
        type: 'error',
        message: 'Failed to fetch CVE report',
        timeoutSecs: 10,
      })
    }

    loading = false
  })
</script>

<div class="px-6 pb-4 space-y-4">
  <div class="mb-6">
    <HeaderWithIcon title="CVE Reports" description="CVE Reports" />
  </div>

  <div class="p-3">
    {#if loading}
      <div class="flex flex-col space-y-6 pb-6 border-b dark:border-gray-700 animate-pulse">
        <div class="grid grid-cols-3 gap-3">
          <div class="flex flex-col h-24 items-center dark:bg-gray-700 rounded-lg overflow-hidden py-3 space-y-2" />
          <div class="flex flex-col h-24 items-center dark:bg-gray-700 rounded-lg overflow-hidden py-3 space-y-2" />
          <div class="flex flex-col h-24 items-center dark:bg-gray-700 rounded-lg overflow-hidden py-3 space-y-2" />
        </div>
      </div>

      <div class="flex flex-col space-y-4 pb-6 mt-3 animate-pulse">
        <div class="flex justify-between">
          <div class="dark:bg-gray-700 h-6 w-[136px]" />
          <div class="dark:bg-gray-700 px-4 rounded-md w-[53px]" />
        </div>

        <div class="flex justify-between">
          <div class="dark:bg-gray-700 h-6 w-[165px]" />
          <div class="dark:bg-gray-700 px-4 rounded-md w-[80px]" />
        </div>

        <div class="flex justify-between">
          <div class="dark:bg-gray-700 h-6 w-[184px]" />
          <div class="dark:bg-gray-700 px-4 rounded-md w-[140px]" />
        </div>
      </div>
    {:else}
      <div class="flex flex-col space-y-6 pb-6 border-b dark:border-gray-700">
        <div class="grid grid-cols-3 gap-3">
          <div class="flex flex-col items-center dark:bg-gray-700 rounded-lg overflow-hidden py-3 space-y-2">
            <div class="flex items-center justify-center h-10 w-10 dark:bg-gray-600 rounded-full">
              <span class="dark:text-red-400">{criticalCount}</span>
            </div>
            <p class="dark:text-red-400">Critical</p>
          </div>
          <div class="flex flex-col items-center dark:bg-gray-700 rounded-lg overflow-hidden py-3 space-y-2">
            <div class="flex items-center justify-center h-10 w-10 dark:bg-gray-600 rounded-full">
              <span class="dark:dark:text-orange-300">{highCount}</span>
            </div>
            <p class="dark:text-orange-300">High</p>
          </div>
          <div class="flex flex-col items-center dark:bg-gray-700 rounded-lg overflow-hidden py-3 space-y-2">
            <div class="flex items-center justify-center h-10 w-10 dark:bg-gray-600 rounded-full">
              <span class="dark:text-yellow-100">{mediumCount}</span>
            </div>
            <p class="dark:text-yellow-100">Medium</p>
          </div>
        </div>
      </div>

      <div class="flex flex-col space-y-4 pb-6 mt-3">
        <div class="flex justify-between">
          <span class="dark:text-gray-400">Total CVE found:</span>
          <span class="dark:bg-gray-700 px-4 rounded-md">{reports.length}</span>
        </div>

        <div class="flex justify-between">
          <span class="dark:text-gray-400">Days since last scan:</span>
          <span class="dark:bg-gray-700 px-4 rounded-md text-center py-0.5">{scanTimeInDays}</span>
        </div>

        <div class="flex justify-between">
          <span class="dark:text-gray-400">CVE database updated:</span>
          <span class="dark:bg-gray-700 px-4 rounded-md text-center py-0.5">{buildDate}</span>
        </div>
      </div>
    {/if}
  </div>

  <div class="bg-white dark:bg-gray-900 flex items-center justify-end pt-5 border-t dark:border-gray-700">
    <button
      class="text-sm text-blue-500 dark:text-blue-300 flex items-center space-x-1"
      on:click={() => goto('/security/cve-report')}
    >
      <span>VIEW REPORT</span>
      <ChevronRight />
    </button>
  </div>
</div>
