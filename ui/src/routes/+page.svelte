<!-- Copyright 2025-2026 Defense Unicorns -->
<!-- SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial -->

<script lang="ts">
  import { AppGrid, isAdminGateway } from '$features/apps'

  import type { LayoutData } from './$types'

  export let data: LayoutData

  $: tenantApps = data.apps.filter((app) => !isAdminGateway(app.gateway))
</script>

<div class="flex w-full flex-col items-center space-y-8">
  <span class="text-2xl font-medium text-gray-100">Your Apps</span>
  <AppGrid apps={tenantApps} />
</div>
