<!-- Copyright 2024 Defense Unicorns -->
<!-- SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial -->

<script lang="ts">
  import { routes } from '../routes'
  import './styles.postcss'
  import { _bannerCfg } from '../../../../routes/+layout'

  $: isActive = true
</script>

<aside
  data-testid="main-sidebar"
  id="sidebar"
  aria-labelledby="sidebar-id"
  class="fixed h-full max-h-[calc(100vh-3.5rem)] top-14 left-0 z-40 w-64"
  class:top-20={$_bannerCfg.enabled}
  class:with-header={$_bannerCfg.enabled && !$_bannerCfg.footer}
  class:with-banners={$_bannerCfg.enabled && $_bannerCfg.footer}
>
  <div
    class="h-full overflow-y-auto  py-5 flex flex-col pr-3 pl-1 mt-20"
  >
    <ul class="space-y-2">
      {#each routes as route}
        <li class={route.class}>
          <a
            href={route.path}
            class="relative flex items-center rounded-lg p-2 text-base font-normal transition-colors duration-200
                {isActive && 'active'}"
          >
            <svelte:component this={route.icon} class="icon ml-4 {isActive && 'active'}" />
            <!-- Label with dynamic styling -->
            <span class="ml-3 font-medium"> {route.name} </span>

            <!-- Active indicator bar -->
            {#if isActive}
              <div class="absolute left-0 top-0 h-full w-1 bg-blue-400"></div>
            {/if}
          </a>
        </li>
      {/each}
    </ul>
  </div>
</aside>

<style>
  /* Calcuating height of sidebar based on the height of the banners (header and footer) + navbar */
  .with-banners {
    height: calc(100vh - 6rem);
  }
  /* Calcuating height of sidebar based on the height of the banner (header only) + navbar */
  .with-header {
    height: calc(100vh - 5rem);
  }
</style>
