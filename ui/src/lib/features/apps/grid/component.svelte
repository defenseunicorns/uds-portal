<!-- Copyright 2026 Defense Unicorns -->
<!-- SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial -->

<script lang="ts">
  import type { ApiApp } from '$lib/types'

  export let apps: ApiApp[]
</script>

<div class="flex w-full justify-center">
  <div class="flex w-full max-w-6xl flex-wrap justify-start gap-8">
    {#each apps as app (app.url)}
      <div
        class="flex h-33 w-33 flex-col items-center justify-center rounded-lg px-4 py-2.5 shadow-sm transition-colors duration-200 hover:bg-gray-800"
      >
        <a
          href="https://{app.url}"
          target="_blank"
          rel="noopener noreferrer"
          class="flex h-full w-full flex-col items-center justify-center"
        >
          <div class="flex w-full flex-col items-center gap-0.5">
            <div class="flex h-12 w-12 items-center justify-center rounded p-2">
              {#if app.icon}
                <img src={app.icon} alt="" class="h-8 w-8" />
              {:else}
                <img src="/default_logo.svg" alt="" class="h-8 w-8" />
              {/if}
            </div>
            <span class="line-clamp-2 w-full text-center text-sm leading-normal font-normal text-gray-300">
              {app.name}
            </span>
          </div>
        </a>
      </div>
    {/each}
  </div>
</div>
