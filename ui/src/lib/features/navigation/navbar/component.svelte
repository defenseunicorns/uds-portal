<!-- Copyright 2025 Defense Unicorns -->
<!-- SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial -->

<script lang="ts">
  import { get } from 'svelte/store'

  import { UserMenu } from '$features/navigation'
  import type { UserData } from '$features/navigation/types'

  export let userData: UserData

  const inClusterAuth = (userData && userData.inClusterAuth) ?? false
</script>

<div class="bg-gray-50 antialiased">
  <nav
    class="fixed left-0 right-0 z-50 border-b border-gray-200 bg-white px-4 py-2.5 dark:border-gray-700 dark:bg-gray-900"
  >
    <div class="flex flex-wrap items-center justify-between">
      <div class="flex items-center justify-start">
        <a href="/" class="mr-4 flex">
          <img src="/doug.svg" class="mr-3 h-8" alt="FlowBite Logo" />
          <span class="self-center whitespace-nowrap text-2xl font-semibold dark:text-gray-100">UDS</span>
        </a>
      </div>
      <div class="flex items-center lg:order-2">
        {#if inClusterAuth}
          <UserMenu {userData} />
        {/if}
      </div>
    </div>
  </nav>
</div>
