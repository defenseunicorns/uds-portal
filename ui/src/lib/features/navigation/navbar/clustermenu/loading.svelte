<script lang="ts">
  import { Spinner } from '$components'
  import { CloudOffline, WarningFilled } from 'carbon-icons-svelte'

  import type { ClusterInfo } from './store'
  import { displayClusterName } from './utils'

  export let cluster: ClusterInfo
  export let error: string | null = null
</script>

<section class="h-full flex justify-center items-center">
  {#if error}
    <div class="flex flex-col items-center w-[95%]">
      <div class="flex items-center mb-2">
        <WarningFilled class="text-red-500 w-10 h-10 mr-2" />
        <h1 class="text-gray-100 font-medium text-2xl">Failed to Connect to {displayClusterName(cluster)}</h1>
      </div>
      <CloudOffline class="w-24 h-24 fill-red-500" />
      <h2 class="text-gray-100 text-xl font-medium mt-2">Failure Message:</h2>
      <h3 class="text-red-500 font-medium text-xl mt-2 text-center w-[50%]">{error}</h3>
    </div>
  {:else}
    <div class="flex flex-col items-center">
      <h1 class="text-gray-100 text-2xl font-medium mb-2">Connecting to cluster</h1>
      <div>
        <Spinner />
      </div>
      <h2 class="text-gray-100 text-xl font-medium mt-2">{displayClusterName(cluster)}</h2>
    </div>
  {/if}
</section>
