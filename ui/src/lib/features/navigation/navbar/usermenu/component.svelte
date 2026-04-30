<!-- Copyright 2025 Defense Unicorns -->
<!-- SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial -->

<script lang="ts">
  import { fade } from 'svelte/transition'

  import type { UserData } from '$features/navigation/types'
  import { ChevronDown, UserAvatar } from 'carbon-icons-svelte'

  export let userData: UserData

  let dropdownContainer: HTMLElement
  let isOpen = false

  function handleWindowPointerDown(event: PointerEvent) {
    if (dropdownContainer && !dropdownContainer.contains(event.target as Node)) {
      isOpen = false
    }
  }

  function toggleMenu() {
    isOpen = !isOpen
  }

  function signOut() {
    window.location.assign('/logout')
  }
</script>

<svelte:window on:pointerdown|capture={handleWindowPointerDown} />

<div class="relative inline-block text-left" bind:this={dropdownContainer}>
  <button
    on:click={toggleMenu}
    class="inline-flex items-center justify-center rounded-md bg-gray-900 px-4 py-2 text-sm font-medium transition-colors duration-200 ease-in-out focus:outline-none dark:text-gray-300 dark:hover:bg-gray-700 dark:hover:text-gray-100 dark:focus:text-gray-100"
  >
    <UserAvatar class="mr-2 h-5 w-5" />
    <span>{userData.username}</span>
    <ChevronDown class="ml-1 h-[12px] w-[12px]" />
  </button>

  {#if isOpen}
    <div
      transition:fade={{ duration: 100 }}
      class="absolute right-0 mt-2 min-w-60 origin-top-right rounded-b-lg border bg-gray-700 shadow-lg focus:outline-none dark:border-gray-700 dark:bg-gray-900"
    >
      <div class="py-1">
        <div class="truncate border-b border-gray-700 px-4 py-2 text-sm font-semibold text-gray-100">
          <p>{userData.name}</p>
        </div>
        <button on:click={signOut} class="w-full px-4 py-2 text-left text-sm text-gray-300 hover:bg-gray-600">
          Sign Out
        </button>
      </div>
    </div>
  {/if}
</div>
