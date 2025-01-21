<!-- Copyright 2024 Defense Unicorns -->
<!-- SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial -->

<script lang="ts">
  import { onDestroy, onMount } from 'svelte'

  import type { KubernetesObject } from '@kubernetes/client-node'
  import { goto } from '$app/navigation'
  import { page } from '$app/stores'
  import type { Row as NamespaceRow } from '$features/k8s/namespaces/store'
  import type { ResourceStoreInterface } from '$features/k8s/types'
  import { addToast } from '$features/toast'
  import { ChevronDown, Information, WarningFilled } from 'carbon-icons-svelte'

  import { clusters, loadingCluster, type ClusterInfo } from './store'
  import { calcClusterDropdownWidth, displayClusterName, noClustersMsg } from './utils'

  let dropdownContainer: HTMLElement
  let dropdown: HTMLDivElement
  let dropdownWidth = 0
  let selected = ''
  let switchErr = false
  let availableClusters: ClusterInfo[] = []
  let searchVal = ''
  let isOpen = false
  let globalNamespaces: ResourceStoreInterface<KubernetesObject, NamespaceRow>

  onMount(() => {
    const handleWindowClick = (event: MouseEvent) => {
      if (dropdownContainer && !dropdownContainer.contains(event.target as Node)) {
        isOpen = false
      }
    }

    window.addEventListener('click', handleWindowClick)

    return () => {
      window.removeEventListener('click', handleWindowClick)
    }
  })

  const unsubscribe = page.subscribe(async ({ data }) => {
    globalNamespaces = data.namespaces as ResourceStoreInterface<KubernetesObject, NamespaceRow>
  })

  onDestroy(() => {
    unsubscribe()
  })

  async function onDropdownClick() {
    searchVal = ''
    if (!isOpen) {
      await fetchClusters()
    }
    isOpen = !isOpen
  }

  async function fetchClusters() {
    try {
      const resp = await fetch('/api/v1/clusters', { method: 'GET' })
      if (resp.ok) {
        clusters.set(await resp.json())
      } else {
        throw new Error(`${await resp.text()}`)
      }
    } catch (e) {
      addToast({ type: 'error', message: e.message, timeoutSecs: 10 })
    }
  }

  async function switchCluster(cluster: ClusterInfo) {
    let resp: Response
    let timeoutPassed = false
    loadingCluster.set({ loading: true, cluster })
    searchVal = ''

    setTimeout(async () => {
      if (resp && resp.ok) {
        loadingCluster.set({ loading: false, cluster })
        goto('/')
      } else {
        timeoutPassed = true
      }
    }, 3000)

    try {
      resp = await fetch(`/api/v1/cluster`, { method: 'POST', body: JSON.stringify({ cluster }) })
      await fetchClusters()
      // close dropdown after getting updated data
      isOpen = false

      if (resp.ok) {
        if (timeoutPassed) {
          loadingCluster.set({ loading: false, cluster })
          goto('/')
        }
        switchErr = false

        if (globalNamespaces.restart) {
          globalNamespaces.restart()
        }
      } else {
        throw new Error(`${await resp.text()}`)
      }
    } catch (e) {
      loadingCluster.update((state) => ({ ...state, err: e.message }))
      switchErr = true
    }
  }

  $: {
    for (const cluster of $clusters) {
      if (cluster.selected) {
        selected = displayClusterName(cluster)
        break
      }
    }
  }

  $: availableClusters = $clusters
    .filter((cluster) => {
      if (!cluster.selected && displayClusterName(cluster).includes(searchVal)) {
        return cluster
      }
    })
    .sort((a, b) => a.name.localeCompare(b.name))

  $: dropdownWidth = calcClusterDropdownWidth($clusters.filter((cluster) => !cluster.selected))
</script>

<div bind:this={dropdownContainer}>
  <button
    on:click={onDropdownClick}
    class="flex items-center justify-between p-2 focus:outline-none transition-colors duration-200 ease-in-out text-sm font-medium text-gray-700 rounded md:border-0 md:w-auto dark:border-gray-700 hover:bg-gray-700 dark:text-gray-300 dark:hover:text-gray-100 dark:focus:text-gray-100"
  >
    {#if switchErr}
      <WarningFilled class="text-red-500 pr-1 h-5 w-5" />
    {/if}
    {selected}
    <ChevronDown class="ml-1 w-[12px] h-[12px]" />
  </button>

  <div
    data-testid="clusterDropdown"
    bind:this={dropdown}
    class="font-normal absolute right-4 z-10 py-2 px-4 mt-2 rounded-b-lg shadow bg-gray-50 dark:bg-gray-900 dark:divide-gray-600 text-sm dark:text-gray-300 border-x-1 border dark:border-gray-700"
    class:hidden={!isOpen}
    style={`min-width: 18rem; width: ${dropdownWidth}px;`}
  >
    {#if (availableClusters.length > 0 && searchVal === '') || searchVal !== ''}
      <div class="p-3">
        <label for="input-group-search" class="sr-only">Search</label>
        <div class="relative">
          <div class="absolute inset-y-0 rtl:inset-r-0 start-0 flex items-center ps-3 pointer-events-none">
            <svg
              class="w-4 h-4 text-gray-500 dark:text-gray-400"
              aria-hidden="true"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 20 20"
            >
              <path
                stroke="currentColor"
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="m19 19-4-4m0-7A7 7 0 1 1 1 8a7 7 0 0 1 14 0Z"
              />
            </svg>
          </div>
          <input
            bind:value={searchVal}
            type="text"
            id="input-group-search"
            data-testid="search"
            class="block w-full p-2 ps-10 text-sm text-gray-900 border border-gray-300 rounded-lg bg-gray-50 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-600 dark:border-gray-500 dark:placeholder-gray-400 dark:text-gray-100 dark:focus:ring-blue-500 dark:focus:border-blue-500"
            placeholder="Search"
          />
        </div>
      </div>
    {/if}
    <ul
      class="min-h-7 px-3 pb-3 max-h-48 overflow-y-auto overflow-x-hidden text-sm text-gray-700 dark:text-gray-200"
      aria-labelledby="dropdownSearchButton"
    >
      {#if availableClusters.length === 0 && searchVal === ''}
        <li>
          <span class="pt-3 text-nowrap cursor-not-allowed flex items-center gap-1">
            <Information class="shrink-0" size={16} />
            {noClustersMsg}
          </span>
        </li>
      {:else}
        {#each availableClusters as cluster}
          <li>
            <button
              class="w-full px-4 py-2 text-left rounded text-nowrap hover:bg-gray-400 dark:hover:bg-gray-600 dark:hover:text-gray-100"
              on:click={() => switchCluster(cluster)}
            >
              {displayClusterName(cluster)}
            </button>
          </li>
        {/each}
      {/if}
    </ul>
  </div>
</div>
