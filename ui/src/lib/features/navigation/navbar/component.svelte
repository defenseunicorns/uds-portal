<!-- Copyright 2025 Defense Unicorns -->
<!-- SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial -->

<script lang="ts">
  import { resolve } from '$app/paths'
  import { UserMenu } from '$features/navigation'
  import type { UserData } from '$features/navigation/types'

  export let userData: UserData
</script>

<div class="bg-gray-50 antialiased">
  <nav class="fixed right-0 left-0 z-50 bg-white px-4 py-2.5 dark:bg-gray-900">
    <div class="flex flex-wrap items-center justify-between">
      <div class="flex items-center justify-start">
        <a href={resolve('/')} class="mr-4 flex">
          <img src="/uds.svg" class="mr-3 h-10" alt="Unicorn Delivery Service Logo" />
        </a>
      </div>
      <div class="flex items-center lg:order-2">
        <UserMenu {userData} />
      </div>
    </div>
  </nav>
</div>
