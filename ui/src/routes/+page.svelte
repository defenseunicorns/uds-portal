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

        // validate response format
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

        // flatten apps by endpoint such that each App has only one endpoint
        apps = splitAppsByEndpoint(appData)

        // process app data for rendering
        apps = apps.map((app) => {
          app.metadata.name = formatAppName(app)
          return app
        })

        return
      }
      console.error('Failed to fetch services')
    } catch (error) {
      console.error('error:', error)
    }
  })

  $: filteredApps = apps.filter((app) => app.metadata.name.toLowerCase().includes(searchValue.toLowerCase())) as App[]

  function clearSearch() {
    searchValue = ''
  }

  function splitAppsByEndpoint(apps: App[]): App[] {
    return apps.flatMap((app) =>
      app.status.endpoints.map((endpoint) => ({
        metadata: { ...app.metadata },
        status: { endpoints: [endpoint] },
        icon: app.icon,
      })),
    )
  }

  function formatAppName(app: App): string {
    if (!app.metadata.name) return ''

    // replace hyphens with spaces
    let formattedName = app.metadata.name.replace(/-/g, ' ')

    // capitalize the first letter of each word
    formattedName = formattedName.replace(/\b\w/g, (char) => char.toUpperCase())

    // capitalize 'uds' if it starts with 'uds'
    if (formattedName.toLowerCase().startsWith('uds')) {
      formattedName = formattedName.replace(/^uds/i, 'UDS')
    }

    // rename sso endpoint to "My Account"
    if (app.status.endpoints[0].startsWith('sso.')) {
      formattedName = 'My Account'
    }

    return formattedName
  }
</script>

<div class="w-full flex flex-col items-center space-y-8">
  <!-- Title -->
  <span class="text-2xl font-medium text-gray-100">My Apps</span>

  <!-- Search Input -->
  <div class="md:w-full max-w-md sm:w-[50%]">
    <div class="relative">
      <div class="pointer-events-none absolute inset-y-0 left-3 flex items-center">
        <Search class="h-4 w-4 text-gray-400" />
      </div>
      <input
        type="text"
        bind:value={searchValue}
        name="input-search"
        autocomplete="off"
        class="block w-full rounded-lg border border-transparent bg-gray-800 py-2 pl-10 pr-8 text-sm text-gray-200
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

  <!-- App Grid -->
  <div class="w-full flex justify-center">
    <div class="flex flex-wrap justify-start gap-8 w-full max-w-3xl">
      {#each filteredApps.sort((a, b) => (a.metadata.name > b.metadata.name ? 1 : -1)) as app}
        <div
          class="w-28 h-24 flex flex-col items-center justify-center rounded-lg hover:bg-gray-700 transition-colors duration-200"
        >
          <a
            href="https://{app.status.endpoints[0]}"
            target="_blank"
            rel="noopener noreferrer"
            class="flex flex-col items-center justify-center w-full h-full"
          >
            <div class="flex items-center justify-center">
              <svelte:component this={app.icon} class="h-8 w-8 text-blue-400" />
            </div>
            <span class="text-md text-gray-100 text-center w-full truncate mt-3">
              {app.metadata.name}
            </span>
          </a>
        </div>
      {/each}
    </div>
  </div>
</div>
