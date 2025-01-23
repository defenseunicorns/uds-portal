<!-- Copyright 2025 Defense Unicorns -->
<!-- SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial -->

<script lang="ts">
  import { onMount } from 'svelte'

  import { Close, Cube, Search } from 'carbon-icons-svelte'

  import { searchStore } from './store'
  import type { App } from './types'

  // Local search value bound to the input
  let searchValue = ''

  // Update store when local value changes
  $: {
    $searchStore = searchValue
  }

  let apps: App[] = []

  onMount(async () => {
    try {
      const url = '/api/v1/apps'
      const response = await fetch(url)
      if (response.ok) {
        const json = await response.json()
        let appData = json as App[]

        // Simple validation that the response matches our type
        if (
          !Array.isArray(appData) ||
          !appData.every((app) => app.metadata?.name && Array.isArray(app.status?.endpoints))
        ) {
          console.error('Invalid response format')
        }
        appData = appData.map((app) => ({
          ...app,
          icon: Cube,
        }))
        apps = appData
        return
      }
      console.error('Failed to fetch services')
    } catch (error) {
      console.error('error:', error)
    }
  })

  $: filteredApps = apps.filter((app) => app.metadata.name.toLowerCase().includes(searchValue.toLowerCase())) as App[]

  // Clear search
  function clearSearch() {
    searchValue = ''
  }
</script>

<div class="flex flex-col items-center space-y-8">
  <span class="text-2xl font-medium text-gray-100">My Apps</span>

  <div class="w-full max-w-md">
    <div class="relative">
      <div class="pointer-events-none absolute inset-y-0 left-3 flex items-center">
        <Search class="h-4 w-4 text-gray-400" />
      </div>

      <input
        type="text"
        bind:value={searchValue}
        name="input-search"
        autocomplete="off"
        class="block w-full rounded-lg border border-transparent bg-gray-900 py-2 pl-10 pr-8 text-sm text-gray-200
                       placeholder-gray-400 outline-none transition-colors duration-200
                       focus:border-gray-700 focus:bg-gray-800 focus:ring-0"
        placeholder="Search"
        data-testid="datatable-search"
      />

      {#if searchValue}
        <button
          class="absolute inset-y-0 right-3 flex items-center text-gray-400 hover:text-gray-200"
          on:click={clearSearch}
        >
          <Close class="h-4 w-4" />
        </button>
      {/if}
    </div>
  </div>

  <!-- Apps Grid -->
  <div
    class="grid w-full max-w-6xl grid-cols-2 gap-8 px-4 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-[repeat(auto-fit, minmax(150px, 1fr))] justify-center place-items-center"
  >
    {#each filteredApps as app}
      {#each app.status.endpoints as endpoint}
        <a
          href="https://{endpoint}"
          target="_blank"
          rel="noopener noreferrer"
          class="group flex flex-col items-center space-y-3 text-center"
        >
          <div
            class="flex h-16 w-16 items-center justify-center rounded-lg
                       transition-colors duration-200 group-hover:bg-gray-700"
          >
            <svelte:component this={app.icon} class="h-8 w-8 text-blue-400" />
          </div>
          <span class="text-sm font-medium text-gray-100">{app.metadata.name}</span>
        </a>
      {/each}
    {/each}
  </div>
</div>
