<!-- Copyright 2025 Defense Unicorns -->
<!-- SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial -->

<script lang="ts">
  import { onMount } from 'svelte'
  import { fade } from 'svelte/transition'

  import { goto } from '$app/navigation'
  import type { UserData } from '$features/navigation/types'
  import { ChevronDown, UserAvatar } from 'carbon-icons-svelte'

  export let userData: UserData

  let dropdownContainer: HTMLElement
  let isOpen = false

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

  function toggleMenu() {
    isOpen = !isOpen
  }

  function signOut() {
    goto('/logout')
  }
</script>

<div class="relative inline-block text-left" bind:this={dropdownContainer}>
  <button
    on:click={toggleMenu}
    class="inline-flex bg-gray-900 dark:hover:bg-gray-700 dark:hover:text-gray-100 dark:focus:text-gray-100 items-center justify-center rounded-md px-4 py-2 text-sm font-medium dark:text-gray-300 focus:outline-none transition-colors duration-200 ease-in-out"
  >
    <UserAvatar class="h-5 w-5 mr-2" />
    <span>{userData.preferredUsername}</span>
    <ChevronDown class="ml-1 h-[12px] w-[12px]" />
  </button>

  {#if isOpen}
    <div
      transition:fade={{ duration: 100 }}
      class="min-w-60 origin-top-right absolute right-0 mt-2 shadow-lg bg-gray-700 focus:outline-none dark:bg-gray-900 rounded-b-lg border dark:border-gray-700"
    >
      <div class="py-1">
        <div class="px-4 py-2 text-sm text-gray-100 font-semibold border-b border-gray-700 truncate">
          <p>{userData.name}</p>
          <p class="text-sm text-gray-300 font-normal mt-1 truncate">User Role: {userData.group}</p>
        </div>
        <button on:click={signOut} class="w-full text-left px-4 py-2 text-sm text-gray-300 hover:bg-gray-600">
          Sign Out
        </button>
      </div>
    </div>
  {/if}
</div>
