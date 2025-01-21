<!-- Copyright 2024 Defense Unicorns -->
<!-- SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial -->

<script lang="ts">
  import { Search, Close, Cube } from 'carbon-icons-svelte'
  import { searchStore } from './store';

  // Local search value bound to the input
  let searchValue = '';

  // Update store when local value changes
  $: {
    $searchStore = searchValue;
  }

  // Sample apps data
  const apps = [
    { id: 1, name: 'UDS Runtime', icon: Cube },
    { id: 2, name: 'Grafana', icon: Cube },
    { id: 3, name: 'Neuvector', icon: Cube },
    { id: 4, name: 'Keycloak Admin', icon: Cube },
    { id: 5, name: 'Matomo', icon: Cube },
    { id: 6, name: 'Jupyter', icon: Cube },
    { id: 7, name: 'App', icon: Cube },
    { id: 8, name: 'App', icon: Cube },
    { id: 9, name: 'App', icon: Cube }
  ];

  $: filteredApps = apps.filter(app =>
    app.name.toLowerCase().includes(searchValue.toLowerCase())
  );

  // Clear search
  function clearSearch() {
    searchValue = '';
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
    <div class="grid w-full max-w-6xl grid-cols-2 gap-8 px-4 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6">
        {#each filteredApps as app}
            <a href="/" class="group flex flex-col items-center space-y-3 text-center">
                <div class="flex h-16 w-16 items-center justify-center rounded-lg bg-gray-800
                           transition-colors duration-200 group-hover:bg-gray-700">
                    <svelte:component this={app.icon} class="h-8 w-8 text-blue-400" />
                </div>
                <span class="text-sm font-medium text-gray-100">{app.name}</span>
            </a>
        {/each}
    </div>
</div>
