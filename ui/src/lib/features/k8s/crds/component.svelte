<!-- Copyright 2024 Defense Unicorns -->
<!-- SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial -->

<script lang="ts">
  import { DataTable } from '$components'
  import { resourceDescriptions } from '$lib/utils/descriptions'

  import { createStore, type Columns } from './store'

  export let columns: Columns = [
    ['name', 'w-3/12'],
    ['group', 'w-3/12 xl:w-2/12'],
    ['kind', 'w-2/12'],
    ['versions', 'line-clamp-4 w-1/12 xl:w-2/12'],
    ['scope', 'w-2/12'],
    ['age', 'w-1/12'],
  ]

  const name = 'CustomResourceDefinitions'
  const description = resourceDescriptions[name]
</script>

<DataTable {columns} {createStore} {name} {description} isNamespaced={false} />
