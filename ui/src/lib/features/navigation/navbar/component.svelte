<!-- Copyright 2024 Defense Unicorns -->
<!-- SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial -->

<script lang="ts">
  import { get } from 'svelte/store'

  import { authenticated } from '$features/auth/store'
  import { UserMenu } from '$features/navigation'
  import type { UserData } from '$features/navigation/types'

  import { isSidebarExpanded } from '../store'
  import { default as ClusterMenu } from './clustermenu/component.svelte'
  import { clusters } from './clustermenu/store'

  export let userData: UserData

  const inClusterAuth = (userData && userData.inClusterAuth) ?? false

  // Don't expand sidebar if api auth is enabled and user is unauthenticated
  $: {
    if ($authenticated) {
      isSidebarExpanded.set(true)
    } else {
      isSidebarExpanded.set(false)
    }
  }
</script>

<div class="bg-gray-50 antialiased">
  <nav
    class="fixed left-0 right-0 z-50 border-b border-gray-200 bg-white px-4 py-2.5 dark:border-gray-700 dark:bg-gray-900"
  >
    <div class="flex flex-wrap items-center justify-between">
      <div class="flex items-center justify-start">
        {#if $authenticated}
          <button
            id="toggle-sidebar-id"
            data-testid="toggle-sidebar"
            aria-label="Toggle Sidebar"
            aria-expanded="true"
            aria-controls="sidebar"
            on:click={() => isSidebarExpanded.update((v) => !v)}
            class="mr-3 hidden cursor-pointer rounded p-2 text-gray-600 hover:bg-gray-100 hover:text-gray-900 lg:inline dark:text-gray-100 dark:hover:bg-gray-700 dark:hover:text-gray-100"
          >
            <svg class="h-5 w-5" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 16 12">
              <path
                stroke="currentColor"
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M1 1h14M1 6h14M1 11h7"
              />
            </svg>
          </button>
        {/if}
        <a href="/" class="mr-4 flex">
          <img src="/doug.svg" class="mr-3 h-8" alt="FlowBite Logo" />
          <span class="self-center whitespace-nowrap text-2xl font-semibold dark:text-gray-100">UDS</span>
        </a>
      </div>
      <div class="flex items-center lg:order-2">
        {#if get(clusters).length > 0}
          <ClusterMenu />
        {/if}
        {#if inClusterAuth}
          <UserMenu {userData} />
        {/if}
      </div>
    </div>
  </nav>
</div>
