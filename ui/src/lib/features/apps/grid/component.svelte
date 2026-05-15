<!-- Copyright 2026 Defense Unicorns -->
<!-- SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial -->

<script lang="ts">
  import { Cube } from 'carbon-icons-svelte'

  import type { ApiApp } from '../../../../routes/types'

  export let apps: ApiApp[]
</script>

<div class="flex w-full justify-center">
  <div class="flex w-full max-w-6xl flex-wrap justify-start gap-8">
    {#each apps as app (app.url)}
      <div
        class="flex h-33 w-33 flex-col items-center justify-center rounded-lg transition-colors duration-200 hover:bg-gray-800"
      >
        <a
          href="https://{app.url}"
          target="_blank"
          rel="noopener noreferrer"
          class="flex h-full w-full flex-col items-center justify-center"
        >
          <div class="flex items-center justify-center p-2">
            {#if app.icon}
              <img src={app.icon} alt="{app.name} icon" class="h-8 w-8" />
            {:else}
              <svelte:component this={Cube} class="h-8 w-8 text-blue-400" />
            {/if}
          </div>
          <span class="mt-2 w-full truncate text-center text-base text-gray-100">
            {app.name}
          </span>
        </a>
      </div>
    {/each}
  </div>
</div>
