<!-- Copyright 2025-2026 Defense Unicorns -->
<!-- SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial -->

<script lang="ts">
  import { AppGrid, isAdminGateway } from '$features/apps'
  import { SearchBar } from '$lib/components'

  import type { LayoutData } from './$types'

  export let data: LayoutData

  let searchQuery = ''

  $: tenantApps = data.apps.filter((app) => !isAdminGateway(app.gateway))
  $: filteredApps = tenantApps.filter((app) => app.name.toLowerCase().includes(searchQuery.toLowerCase()))
</script>

<div class="flex w-full flex-col items-center space-y-8">
  <h1 class="text-2xl font-medium text-gray-100">Your Apps</h1>
  <div class="flex w-full max-w-[420px] items-center gap-2.5">
    <SearchBar bind:value={searchQuery} />
  </div>
  <AppGrid apps={filteredApps} />
  {#if filteredApps.length === 0 && searchQuery}
    <img src="/uds_text.svg" alt="" aria-hidden="true" class="h-[200px] w-[200px]" />
  {/if}
</div>
