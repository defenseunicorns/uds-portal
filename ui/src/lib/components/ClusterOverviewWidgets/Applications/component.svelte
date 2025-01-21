<script lang="ts">
  import OpenExternalLinkIcon from '$lib/icons/cluster-overview/OpenExternalLinkIcon.svelte'
  import { type ClusterOverviewUDSPackageType } from '$lib/types'
  import { Cube } from 'carbon-icons-svelte'
  import * as _ from 'lodash'

  import HeaderWithIcon from '../../k8s/HeaderWithIcon/component.svelte'

  import '../styles.postcss'

  export let udsPackages: ClusterOverviewUDSPackageType[] = []

  let transformedPackagesList: { name: string; endpoint: string }[] = []

  $: {
    transformedPackagesList = []

    if (udsPackages) {
      udsPackages.forEach((service) => {
        let name = service.metadata.name
        let endpoint = service.status && service.status.endpoints ? service.status.endpoints[0] : ''

        if (endpoint) {
          transformedPackagesList = [
            ...transformedPackagesList,
            {
              name,
              endpoint,
            },
          ]
        }
      })
    }
  }

  function openLink(endpoint: string) {
    window.open('https://' + endpoint, '_blank')
  }
</script>

<div class="overview-widget">
  <HeaderWithIcon
    height="14"
    title="Applications"
    description="Applications are UDS Packages that are deployed and accessible to users via provided endpoints"
  />

  {#if transformedPackagesList.length === 0}
    <span class="flex self-center">No Application Packages running</span>
  {:else}
    <div class="core-services__rows">
      {#each _.sortBy(transformedPackagesList, 'name') as { name, endpoint }}
        <div class="overview-widget__rows-item">
          <div class="w-8/12 flex items-center space-x-2">
            <div class="overview-widget__name-icon">
              <Cube size={16} class="text-gray-100" />
            </div>
            <div class="truncate">{name}</div>
          </div>
          <div class="w-4/12 flex justify-end" data-testid="app-widget-container">
            <button
              aria-label="open-link-button"
              on:click={() => openLink(endpoint)}
              data-testid={`app-widget-${name}`}
            >
              <OpenExternalLinkIcon />
            </button>
          </div>
        </div>
      {/each}
    </div>
  {/if}
</div>
