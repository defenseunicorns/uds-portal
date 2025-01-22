<!-- Copyright 2025 Defense Unicorns -->
<!-- SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial -->

<script lang="ts">
  export let primaryMsg = ''
  export let secondaryMsg = ''
</script>

<div class="flex flex-col items-center justify-center text-gray-400 space-y-4">
  <div class="text-center space-y-1">
    <h3 class="text-lg font-medium text-gray-300">{primaryMsg}</h3>
    <p class="text-sm">{secondaryMsg}</p>
  </div>
</div>
