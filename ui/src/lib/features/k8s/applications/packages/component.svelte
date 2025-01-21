<!-- Copyright 2024 Defense Unicorns -->
<!-- SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial -->

<script lang="ts">
  import type { KubernetesObject, V1Secret } from '@kubernetes/client-node'
  import { DataTable } from '$components'
  import type { DrawerDetails, DrawerOpts } from '$components/k8s/Drawer/types'
  import type { DeployedPackage } from '$features/k8s/applications/packages/types'

  import { createStore, type Columns } from './store'

  export let columns: Columns = [
    ['name', 'w-2/12'],
    ['version', ' w-2/12'],
    ['description', 'italic w-4/12 line-clamp-2'],
    ['arch', ' w-1/12'],
    ['flavor', ' w-2/12'],
    ['age', ' w-1/12'],
  ]
  const name = 'Application Packages'
  const description = 'Packages are either UDS or Zarf packages that are deployed in the cluster.'
  const includeEvents = false
  const isNamespaced = false
  const minWidth = 1250

  const addDetails = (resource: KubernetesObject | DeployedPackage): DrawerDetails => {
    resource = resource as DeployedPackage
    return [
      { label: 'Name', value: resource.name },
      { label: 'Version', value: (resource.data.metadata && resource.data.metadata.version) ?? '' },
      { label: 'Description', value: (resource.data.metadata && resource.data.metadata.description) ?? '' },
      { label: 'Flavor', value: (resource.data.build && resource.data.build.flavor) ?? '' },
      { label: 'Architecture', value: (resource.data.metadata && resource.data.metadata.architecture) ?? '' },
      { label: 'Components', value: resource.deployedComponents.map((c) => c.name), isList: true },
    ]
  }

  const transformResource = (resource: V1Secret) => {
    // grab data field from Secret
    let data = JSON.parse(atob(resource.data?.data ?? ''))

    // add metadata and kind fields to data so it renders properly in Details panel
    data.metadata = { name: data.data.metadata.name }
    data.kind = data.data.kind
    return data
  }

  const drawerOpts: DrawerOpts = { includeEvents, addDetails, transformResource }
</script>

<DataTable {columns} {createStore} {name} {description} {isNamespaced} {minWidth} {drawerOpts} />
