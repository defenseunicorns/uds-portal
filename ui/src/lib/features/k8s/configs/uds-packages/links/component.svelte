<!-- Copyright 2024 Defense Unicorns -->
<!-- SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial -->

<script lang="ts">
  import { Link } from '$components'

  export let endpoints: string[]
</script>

{#if endpoints.length < 1}
  <p class="text-gray-500 dark:text-gray-100">-</p>
{/if}

<div class="flex flex-col">
  {#each endpoints as endpoint}
    <div>
      <Link href={`https://${endpoint}`} text={endpoint} />
    </div>
  {/each}
</div>
