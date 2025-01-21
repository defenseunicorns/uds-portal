<script lang="ts">
  import type { V1Pod } from '@kubernetes/client-node'
  import type { UDSPackageStatus } from '$features/k8s/types'
  import { getColorForStatus } from '$lib/features/k8s/helpers'
  import { type ClusterOverviewUDSPackageType } from '$lib/types'
  import { resourceDescriptions } from '$lib/utils/descriptions'
  import { Cube } from 'carbon-icons-svelte'
  import * as _ from 'lodash'

  import HeaderWithIcon from '../../k8s/HeaderWithIcon/component.svelte'

  import '../styles.postcss'

  export let coreServices: ClusterOverviewUDSPackageType[] = []
  export let pods: V1Pod[] = []

  const coreServicesMapping: Record<string, string> = {
    authservice: 'Authorization',
    grafana: 'Monitoring',
    istio: 'Service Mesh',
    keycloak: 'Identity Access Management',
    loki: 'Log Aggregation',
    'metrics-server': 'Metrics',
    neuvector: 'Container Security',
    'prometheus-stack': 'Monitoring',
    vector: 'Log Aggregation',
    velero: 'Backup & Restore',
    'uds-runtime': 'Frontend Views & Insights',
  }

  let transformedCoreServiceList: { name: string; status: UDSPackageStatus }[] = []
  let hasPolicyEngineOperator: boolean = false
  let hasServiceMesh: boolean = false

  $: {
    transformedCoreServiceList = []
    hasPolicyEngineOperator = pods.filter((pod: V1Pod) => pod?.metadata?.name?.match(/^pepr-uds-core/)).length > 0
    hasServiceMesh =
      pods
        .filter((pod: V1Pod) => pod?.metadata?.name?.match(/^istiod/))
        .filter((pod) => pod.status && pod.status?.phase === 'Running').length === 1

    if (coreServices) {
      coreServices.forEach((service) => {
        let name = coreServicesMapping[service.metadata.name]
        let status = service.status ? service.status.phase : 'Pending'

        if (Object.keys(coreServicesMapping).includes(service.metadata.name)) {
          transformedCoreServiceList = [
            ...transformedCoreServiceList,
            {
              name,
              status,
            },
          ]
        }
      })
    }

    // If we have pepr uds core pods then we have a Policy Engine & Operator in the cluster
    if (hasPolicyEngineOperator) {
      transformedCoreServiceList = [
        ...transformedCoreServiceList,
        {
          name: 'Policy Engine & Operator',
          status: 'Ready',
        },
      ]
    }

    if (hasServiceMesh) {
      transformedCoreServiceList = [
        ...transformedCoreServiceList,
        {
          name: 'Service Mesh',
          status: 'Ready',
        },
      ]
    }
  }
</script>

<div class="overview-widget">
  <HeaderWithIcon height="14" title="Core Services" description={resourceDescriptions['UDSPackage']} />

  {#if transformedCoreServiceList.length === 0}
    <span class="flex self-center">No Core Services running</span>
  {:else}
    <div class="overview-widget__rows">
      <!-- Remove duplicates and sort by name -->
      {#each _.sortBy(_.uniqBy(transformedCoreServiceList, 'name'), 'name') as { name, status }}
        <div class="overview-widget__rows-item">
          <div class="w-11/12 flex items-center space-x-2">
            <div class="overview-widget__name-icon">
              <Cube size={16} class="text-gray-100" />
            </div>

            <div class="truncate">{name}</div>
          </div>

          <div class={`w-1/12 ${getColorForStatus('UDSPackage', status)}`}>
            {status}
          </div>
        </div>
      {/each}
    </div>
  {/if}
</div>
