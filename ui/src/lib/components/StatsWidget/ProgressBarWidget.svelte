<!-- Copyright 2024 Defense Unicorns -->
<!-- SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial -->

<script lang="ts">
  import { Card, InactiveBadge, ProgressBar } from '$components'

  import type { BarSizeType, UnitType } from './types.ts'

  export let capacity: number = 0
  export let progress: number = 0
  export let barSize: BarSizeType = 'sm'
  export let statText: string
  export let unit: UnitType
  export let value: number | string
  export let deactivated: boolean = false
</script>

<Card>
  <div class="w-full">
    <div class="w-full">
      <div class="flex justify-between items-center">
        <span class="text-sm font-medium text-gray-400 dark:text-gray-400 truncate">{statText}</span>
        {#if deactivated}
          <div class="flex justify-end">
            <InactiveBadge
              tooltipDirection="tooltip-left"
              tooltipText="Metrics Server is unavailable.
              Ensure Metrics Server is running in the cluster."
            />
          </div>
        {/if}
      </div>
      <span
        class="mt-1 text-3xl font-semibold"
        class:text-gray-900={!deactivated}
        class:text-gray-400={deactivated}
        class:dark:text-gray-100={!deactivated}
      >
        {value.toString()}
      </span>
    </div>

    {#if !deactivated}
      <ProgressBar size={barSize} {progress} {capacity} {unit} />
    {/if}
  </div>
</Card>
