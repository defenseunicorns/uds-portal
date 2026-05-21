<!-- Copyright 2026 Defense Unicorns -->
<!-- SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial -->

<script lang="ts">
  export let value = ''
</script>

<div class="flex w-full items-center gap-2.5 rounded-lg border border-gray-700 bg-gray-800 px-4 py-2">
  <svg
    class="h-4 w-4 shrink-0 text-gray-300"
    fill="none"
    viewBox="0 0 16 16"
    xmlns="http://www.w3.org/2000/svg"
    aria-hidden="true"
  >
    <path
      d="M7.333 12.667A5.333 5.333 0 1 0 7.333 2a5.333 5.333 0 0 0 0 10.667ZM14 14l-2.667-2.667"
      stroke="currentColor"
      stroke-linecap="round"
      stroke-linejoin="round"
      stroke-width="1.333"
    />
  </svg>

  <input
    bind:value
    class="flex-1 bg-transparent text-sm font-normal text-gray-300 outline-none placeholder:text-gray-300"
    placeholder="Search"
    type="text"
  />

  <button
    aria-label="clear"
    class="flex shrink-0 items-center justify-center text-gray-300"
    class:hidden={!value}
    on:click={() => (value = '')}
    type="button"
  >
    <svg class="h-2.5 w-2.5" fill="none" viewBox="0 0 10 10" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
      <path d="M1 1l8 8M9 1l-8 8" stroke="currentColor" stroke-linecap="round" stroke-width="1.333" />
    </svg>
  </button>
</div>
