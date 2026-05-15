<!-- Copyright 2026 Defense Unicorns -->
<!-- SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial -->

<script lang="ts">
  import { resolve } from '$app/paths'
  import { page } from '$app/stores'

  import type { ApiApp } from '../../../../routes/types'
  import { routes } from '../routes'
  import type { Route } from '../types'

  export let apps: ApiApp[]
  export let adminAppsEnabled: boolean

  const isAdminAppsRoute = (r: Route) => r.path === '/admin-apps'

  $: visibleRoutes = routes.filter(
    (r: Route) => (!r.visible || r.visible(apps)) && (!isAdminAppsRoute(r) || adminAppsEnabled),
  )
  $: activePath = $page.url.pathname
</script>

<aside class="flex w-44 flex-col gap-1 py-4 pl-4">
  {#each visibleRoutes as route (route.path)}
    {@const active = activePath === route.path}
    <a
      href={resolve(route.path as Parameters<typeof resolve>[0])}
      data-testid="sidebar-link"
      class="flex items-center gap-2 rounded px-2 py-2 text-sm transition-colors duration-150"
      class:text-blue-400={active}
      class:text-gray-200={!active}
      class:hover:bg-gray-800={!active}
    >
      {#if route.icon}
        <svelte:component this={route.icon} class="h-4 w-4" />
      {/if}
      <span>{route.name}</span>
    </a>
  {/each}
</aside>
