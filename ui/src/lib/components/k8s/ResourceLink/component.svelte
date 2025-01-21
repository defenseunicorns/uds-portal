<!-- Copyright 2024 Defense Unicorns -->
<!-- SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial -->

<script lang="ts">
  import type { KubernetesObject } from '@kubernetes/client-node'
  import { goto } from '$app/navigation'
  import { addToast } from '$features/toast'

  export let name = ''
  export let route = ''
  export let resourcePath = ''
  export let resources: KubernetesObject[] | undefined = undefined

  const gotoResource = () => {
    const resource = resources?.find((r) => r.metadata?.name === name)
    if (resource) {
      goto(`/${route}/${resource.metadata?.uid}`)
    } else {
      addToast({
        timeoutSecs: 5,
        message: `Resource ${name} not found`,
        type: 'error',
      })
    }
  }

  const fetchResources = async () => {
    try {
      const resp = await fetch(`${resourcePath}?once=true&fields=.metadata`)
      if (resp.ok) {
        resources = await resp.json()
      } else {
        throw new Error(`${await resp.text()}`)
      }
    } catch (error) {
      console.error(`Failed to fetch resources: ${error}`)
      addToast({
        timeoutSecs: 5,
        message: `Failed to fetch ${name} resource`,
        type: 'error',
      })
    }
  }

  const handleClick = async (event: MouseEvent) => {
    event.stopPropagation()
    if (resources) {
      gotoResource()
    } else {
      await fetchResources()
      gotoResource()
    }
  }
</script>

<button
  on:click|self={(event) => handleClick(event)}
  class="font-medium text-blue-600 dark:text-blue-500 hover:underline pr-4 text-left"
>
  {name}
</button>
