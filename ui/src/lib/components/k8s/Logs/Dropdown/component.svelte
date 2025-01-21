<!-- Copyright 2024 Defense Unicorns -->
<!-- SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial -->

<script lang="ts">
  import { Dropdown } from 'flowbite'

  export let items: string[]
  export let selectedValue: string | null = null
  export let width: string = 'w-32' // must be a tailwind class
  export let maxHeight: string = 'max-h-32' // must be a tailwind class
  export let id: string
  export let placeholder: string = 'Select an option' // new placeholder prop
  export let disabled: boolean = false
  export let label: string = ''

  function handleSelect(item: string) {
    selectedValue = item

    // hide flowbite dropdown because it won't do it for us
    const dropdownButton = document.getElementById(id)
    const dropDownMenu = document.getElementById(`${id}-menu`)
    const dropdown = new Dropdown(dropDownMenu, dropdownButton)
    dropdown.hide()
  }

  // Helper to determine what text to display
  $: displayText = selectedValue || placeholder
  $: isPlaceholderShown = !selectedValue
</script>

<div class="relative">
  <span class="absolute -top-2 left-1 px-1 text-xs text-gray-100 bg-gray-800 rounded-md">
    {label}
  </span>
  <button
    {id}
    data-dropdown-toggle="{id}-menu"
    class="{width} justify-between truncate bg-gray-800 text-gray-100 hover:border-gray-500 border border-gray-700 text-sm px-2 py-1 rounded-md inline-flex items-center h-10 focus:outline-none focus:ring-0 focus:border-blue-600"
    type="button"
    {disabled}
  >
    <span class="truncate mr-1" class:text-gray-500={isPlaceholderShown} class:text-gray-100={!isPlaceholderShown}>
      {displayText}
    </span>
    <svg
      class="w-2.5 h-2.5 flex-shrink-0"
      aria-hidden="true"
      xmlns="http://www.w3.org/2000/svg"
      fill="none"
      viewBox="0 0 10 6"
    >
      <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="m1 1 4 4 4-4" />
    </svg>
  </button>

  <!-- Dropdown Menu -->
  <div id="{id}-menu" class="z-50 hidden absolute truncate bg-gray-800 rounded-md border border-gray-700 shadow-lg">
    <ul class="{maxHeight} {width} overflow-y-auto text-sm text-gray-100" aria-labelledby={id}>
      {#each items as item}
        <li>
          <button
            on:click={() => handleSelect(item)}
            class="w-full px-3 py-1 text-left hover:bg-gray-700 whitespace-nowrap"
          >
            {item}
          </button>
        </li>
      {/each}
    </ul>
  </div>
</div>
