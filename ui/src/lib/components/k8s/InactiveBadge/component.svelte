<!-- Copyright 2025 Defense Unicorns -->
<!-- SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial -->

<script>
  import { Information } from 'carbon-icons-svelte'

  export let tooltipDirection = 'tooltip-right'
  export let tooltipText = ''
</script>

<span
  class="bg-gray-100 text-gray-700 text-xs font-medium inline-flex items-center px-2.5 py-0.5 rounded me-2 dark:bg-gray-100 dark:text-gray-700 border border-gray-700"
  data-testid="unavailable-tag"
>
  <div class="relative group mr-2">
    <Information class="w-4 h-4 text-gray-700" />
    <div class={`tooltip ${tooltipDirection} min-w-56`}>
      <div class="whitespace-normal">
        {tooltipText}
      </div>
    </div>
  </div>
  Unavailable
</span>
