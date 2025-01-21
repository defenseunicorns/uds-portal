<!-- Copyright 2024 Defense Unicorns -->
<!-- SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial -->

<script lang="ts">
  import { DataTable } from '$components'
  import { resourceDescriptions } from '$lib/utils/descriptions'

  import { createStore, type Columns } from './store'

  export let columns: Columns = [
    ['name', 'w-1/12'],
    ['namespace', 'w-2/12'],
    ['status', 'w-1/12'],
    ['endpoints', 'w-2/12'],
    ['monitors', 'w-2/12'],
    ['ssoClients', 'w-2/12'],
    ['networkPolicies', 'w-1/12'],
    ['age', 'w-1/12'],
  ]

  const name = 'UDS Packages'
  const description = resourceDescriptions[name.split(' ').join('')]
</script>

<DataTable {columns} {createStore} {name} {description} minWidth={2000} />
