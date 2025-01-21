<!-- Copyright 2024 Defense Unicorns -->
<!-- SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial -->

<script lang="ts">
  import { onMount } from 'svelte'

  import type { CoreV1Event, KubernetesObject } from '@kubernetes/client-node'
  import { goto } from '$app/navigation'
  import { EventList } from '$components'
  import { LogsIcon } from '$lib/icons'
  import { Close } from 'carbon-icons-svelte'
  import DOMPurify from 'dompurify'
  import hljs from 'highlight.js/lib/core'
  import yaml from 'highlight.js/lib/languages/yaml'
  import * as YAML from 'yaml'

  import './styles.postcss'

  import type { DrawerDetails, DrawerOpts } from '$components/k8s/Drawer/types'

  import { _bannerCfg } from '../../../../routes/+layout'

  export let resource: KubernetesObject
  export let baseURL: string

  export let drawerOpts: DrawerOpts = {}
  const { addDetails, includeEvents = true, transformResource } = drawerOpts

  type Tab = 'details' | 'yaml' | 'events'

  let events: CoreV1Event[] = []
  let eventSource: EventSource

  onMount(() => {
    // initialize highlight language
    hljs.registerLanguage('yaml', yaml)

    if (includeEvents) {
      const path: string = '/api/v1/resources/events?fields=.count,.involvedObject,.message,.source,.type'
      eventSource = new EventSource(path)
      eventSource.onmessage = (event) => {
        events = JSON.parse(event.data) as CoreV1Event[]
      }
    }

    const handleKeydown = (e: KeyboardEvent) => {
      const tabList: Tab[] = includeEvents ? ['details', 'events', 'yaml'] : ['details', 'yaml']
      let targetTab: string | undefined

      switch (e.key) {
        // If the Escape key is pressed, close the panel by navigating to the base URL
        case 'Escape':
          goto(baseURL)
          return

        // If the left arrow key is pressed, move to the previous tab
        case 'ArrowLeft':
          targetTab = tabList[tabList.indexOf(activeTab) - 1]
          break

        // If the right arrow key is pressed, move to the next tab
        case 'ArrowRight':
          targetTab = tabList[tabList.indexOf(activeTab) + 1]
          break
      }

      // Only update the active tab if the target tab is valid
      if (targetTab) {
        activeTab = targetTab as Tab
      }
    }

    // Add the event listener when the component is mounted
    window.addEventListener('keydown', handleKeydown)

    // Clean up the event listener when the component is destroyed
    return () => {
      window.removeEventListener('keydown', handleKeydown)
      if (eventSource) {
        eventSource.close()
      }
    }
  })

  function formatDate(dateString: string) {
    return new Date(dateString).toLocaleString()
  }

  const defaultDetails = [
    { label: 'Created', value: formatDate(resource.metadata?.creationTimestamp as unknown as string) },
    { label: 'Name', value: resource.metadata?.name },
    { label: 'Namespace', value: resource.metadata?.namespace },
  ] as DrawerDetails

  $: resource = (transformResource && (transformResource(resource) as KubernetesObject)) || resource
  $: details = addDetails ? addDetails(resource) : defaultDetails

  if ((resource.metadata?.ownerReferences?.length && details) || 0 > 0) {
    details.push({
      label: 'Controlled By',
      value: `${resource.metadata?.ownerReferences?.[0]?.kind} ${resource.metadata?.ownerReferences?.[0]?.name}`,
    })
  }

  let activeTab: Tab = 'details'

  function setActiveTab(evt: MouseEvent) {
    const target = evt.target as HTMLButtonElement
    activeTab = target.id as Tab
  }
</script>

<div
  data-testid="drawer"
  class="fixed top-14 right-0 z-40 h-full max-h-[calc(100vh-3.5rem)] overflow-y-auto w-1/2 dark:bg-gray-800 shadow-2xl shadow-black/80 transform transition-transform duration-300 ease-in-out"
  class:top-20={$_bannerCfg.enabled}
  class:with-header={$_bannerCfg.enabled && !$_bannerCfg.footer}
  class:with-banners={$_bannerCfg.enabled && $_bannerCfg.footer}
>
  <div class="flex flex-col h-full">
    <!-- Dark header -->
    <div class="bg-gray-900 text-gray-100 p-4 pb-0">
      <div class="flex justify-between items-center">
        <h2 class="text-xl">
          <span class="font-semibold">{resource.kind}:</span>
          <span>{resource.metadata?.name}</span>
        </h2>
        <div class="flex space-x-1">
          <!-- Show logs button only for Pods -->
          {#if resource.kind === 'Pod'}
            <div class="relative group">
              <button
                type="button"
                data-testid="view-logs-btn"
                class="text-gray-100 bg-transparent hover:bg-gray-200 hover:text-gray-900 rounded-lg text-sm w-8 h-8 top-2.5 end-2.5 inline-flex items-center justify-center dark:hover:bg-gray-600 dark:hover:text-gray-100"
                on:click={() =>
                  goto(
                    `${baseURL}/${resource.metadata?.uid}/logs?pod=${resource.metadata?.name}&namespace=${resource.metadata?.namespace}`,
                  )}
              >
                <LogsIcon class="w-5 h-5" />
              </button>
              <div class="details-tooltip bg-gray-700 whitespace-nowrap">
                <span>View Logs</span>
              </div>
            </div>
          {/if}
          <div class="relative group">
            <button
              type="button"
              class="text-gray-100 bg-transparent hover:bg-gray-200 hover:text-gray-900 rounded-lg text-sm w-8 h-8 top-2.5 end-2.5 inline-flex items-center justify-center dark:hover:bg-gray-600 dark:hover:text-gray-100"
              on:click={() => goto(baseURL)}
            >
              <Close class="w-5 h-5" />
            </button>
            <div class="details-tooltip bg-gray-700 whitespace-nowrap">
              <span>Close</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Tabs -->
      <div class="flex font-medium pt-3">
        <ul class="flex w-full" id="drawer-tabs">
          <li class="flex-1">
            <button id="details" class:active={activeTab === 'details'} on:click={setActiveTab}>Details</button>
          </li>
          {#if includeEvents}
            <li class="flex-1">
              <button id="events" class:active={activeTab === 'events'} on:click={setActiveTab}>Events</button>
            </li>
          {/if}
          <li class="flex-1">
            <button id="yaml" class:active={activeTab === 'yaml'} on:click={setActiveTab}>YAML</button>
          </li>
        </ul>
      </div>
    </div>

    <!-- Content -->
    <div class="flex-grow overflow-y-auto dark:text-gray-100 pb-20">
      {#if activeTab === 'details'}
        <!-- Details tab -->
        <div class="bg-gray-800 text-gray-100 p-6 rounded-lg">
          <dl class="mt-6 space-y-4">
            {#each details as { label, value, isList }}
              {#if isList}
                <div class="flex flex-col sm:flex-row gap-9 border-b border-gray-700 pb-2">
                  <dt class="font-bold text-sm flex-none w-[180px]">{label}</dt>
                  <dd class="text-gray-300">
                    <ul class="list-disc list-inside">
                      {#each value as item}
                        <li>{item}</li>
                      {/each}
                    </ul>
                  </dd>
                </div>
              {:else}
                <div class="flex flex-col sm:flex-row gap-9 border-b border-gray-700 pb-2">
                  <dt class="font-bold text-sm flex-none w-[180px]">{label}</dt>
                  <dd class="text-gray-300">{value || 'Not provided'}</dd>
                </div>
              {/if}
            {/each}

            {#if resource.metadata?.labels}
              <div class="flex flex-col sm:flex-row gap-9 border-b border-gray-700 pb-2">
                <dt class="font-bold text-sm flex-none w-[180px]">Labels</dt>
                <dd class="overflow-x-auto">
                  <div class="flex flex-wrap gap-2">
                    {#each Object.entries(resource.metadata?.labels || {}) as [key, value]}
                      <span class="bg-gray-700 px-2 py-0.5 rounded-md text-xs">
                        <span class="text-gray-100">{key}:</span>
                        <span class="text-gray-300">{value}</span>
                      </span>
                    {/each}
                  </div>
                </dd>
              </div>
            {/if}

            {#if resource.metadata?.annotations}
              <div class="flex flex-col sm:flex-row gap-9">
                <dt class="font-bold text-sm flex-none w-[180px]">Annotations</dt>
                <dd class="overflow-x-auto">
                  <div class="flex flex-wrap gap-2">
                    {#each Object.entries(resource.metadata?.annotations || {}) as [key, value]}
                      <span class="bg-gray-700 px-2 py-0.5 rounded-md text-xs">
                        <span class="text-gray-100">{key}:</span>
                        <span class="text-gray-300">{value}</span>
                      </span>
                    {/each}
                  </div>
                </dd>
              </div>
            {/if}
          </dl>
        </div>
      {:else if activeTab === 'events'}
        <EventList {events} {resource} />
      {:else if activeTab === 'yaml'}
        <!-- YAML tab -->
        <div class="text-gray-100 p-4">
          <code class="text-sm text-gray-500 dark:text-gray-100 whitespace-pre w-full block">
            <!-- Disable svelte/no-at-html-tags eslint rule here because we are using DOMPurify to sanitize -->
            <!-- eslint-disable-next-line svelte/no-at-html-tags -->
            {@html DOMPurify.sanitize(hljs.highlight(YAML.stringify(resource), { language: 'yaml' }).value)}
          </code>
        </div>
      {/if}
    </div>
  </div>
</div>

<style lang="postcss">
  /* Tooltip container styles */
  .details-tooltip {
    @apply absolute top-full left-1/2 transform -translate-x-[50%] mt-2 p-2 bg-gray-700 text-gray-100 text-xs rounded shadow-lg opacity-0 transition-opacity duration-300 z-50 pointer-events-none whitespace-nowrap;
  }

  /* Show tooltip on hover */
  .group:hover .details-tooltip {
    @apply opacity-90 pointer-events-auto;
  }

  /* Calcuating height of sidebar based on the height of the banner (header only) + navbar */
  .with-header {
    height: calc(100vh - 5rem);
  }

  /* Calcuating height of sidebar based on the height of the banners (header and footer) + navbar */
  .with-banners {
    height: calc(100vh - 6.5rem);
  }
</style>
