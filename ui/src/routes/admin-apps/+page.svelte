<!-- Copyright 2026 Defense Unicorns -->
<!-- SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial -->

<script lang="ts">
  import { AppGrid, isAdminGateway } from '$features/apps'

  import type { LayoutData } from '../$types'

  export let data: LayoutData

  $: adminApps = data.apps.filter((app) => isAdminGateway(app.gateway))
</script>

{#if !data.adminAppsEnabled}
  <div class="flex w-full flex-col items-center space-y-8">
    <span class="text-2xl font-medium text-gray-100">Admin Apps</span>
    <p class="text-base text-gray-400">Admin Apps is disabled in this deployment.</p>
  </div>
{:else if adminApps.length > 0}
  <div class="flex w-full flex-col items-center space-y-8">
    <span class="text-2xl font-medium text-gray-100">Admin Apps</span>
    <div
      data-testid="admin-access-banner"
      class="w-full max-w-6xl rounded-md border border-yellow-700/40 bg-yellow-900/20 px-4 py-3 text-sm text-yellow-200"
    >
      These apps may require additional network access (e.g. VPN or a jumphost).
    </div>
    <AppGrid apps={adminApps} />
  </div>
{/if}
