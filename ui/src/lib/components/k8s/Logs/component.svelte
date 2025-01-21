<!-- Copyright 2024 Defense Unicorns -->
<!-- SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial -->

<script lang="ts">
  import { onDestroy, onMount } from 'svelte'
  import { get } from 'svelte/store'
  import { fade } from 'svelte/transition'

  import { goto } from '$app/navigation'
  import { Dropdown } from '$components'
  import { HintMessage } from '$components/index.js'
  import { getAllLogs, getLastNLines } from '$components/k8s/Logs/helpers'
  import {
    autoScroll,
    containerName,
    filteredContainers,
    filteredPods,
    handleAutoscroll,
    handleContainer,
    handleNamespace,
    handlePod,
    handleSearch,
    logElements,
    logEventSource,
    namespaces,
    podName,
    pods,
    scrollTimeout,
    searchTerm,
    selectedNamespace,
  } from '$components/k8s/Logs/store'
  import type { Pod } from '$features/k8s/types'
  import { LogsIcon } from '$lib/icons'
  import { CheckmarkOutline, Close, Copy, Download, Search } from 'carbon-icons-svelte'

  let podEventSource: EventSource

  // manage notifications
  let showNotification = false
  let notificationTimeout: ReturnType<typeof setTimeout>
  let notificationMessage = ''

  const handleCopy = async () => {
    const lastLines = getLastNLines(100)
    try {
      await navigator.clipboard.writeText(lastLines)
      notificationMessage = 'Copied!'
      showNotification = true

      // Clear existing timeout if there is one
      if (notificationTimeout) {
        clearTimeout(notificationTimeout)
      }

      // Hide notification after 2 seconds
      notificationTimeout = setTimeout(() => {
        showNotification = false
      }, 2000)
    } catch (err) {
      console.error('Failed to copy logs:', err)
    }
  }

  export const handleDownload = () => {
    try {
      const logs = getAllLogs()
      const blob = new Blob([logs], { type: 'text/plain' })
      const url = URL.createObjectURL(blob)

      const link = document.createElement('a')
      const pod = get(podName)
      link.href = url
      link.download = `logs-${pod}-${new Date().toISOString().slice(0, 19).replace(/[:]/g, '-')}.txt`

      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)

      URL.revokeObjectURL(url)

      notificationMessage = 'Downloading'
      showNotification = true

      if (notificationTimeout) {
        clearTimeout(notificationTimeout)
      }

      notificationTimeout = setTimeout(() => {
        showNotification = false
      }, 2000)
    } catch (err) {
      console.error('Failed to download logs:', err)
    }
  }

  onDestroy(cleanupPage)

  // set up store handlers
  handleNamespace()
  handlePod()
  handleContainer()
  handleSearch()
  handleAutoscroll()

  onMount(() => {
    // check url query params for name and namespace
    const urlParams = new URLSearchParams(window.location.search)
    const podParam = urlParams.get('pod')
    const nsParam = urlParams.get('namespace')

    // set up sse stream for pods
    const fieldSelectors = 'metadata.name,metadata.namespace,spec.containers[].name'
    const podURL = `/api/v1/resources/workloads/pods?fields=${fieldSelectors}`
    podEventSource = new EventSource(podURL)

    podEventSource.onmessage = (event) => {
      const data = JSON.parse(event.data)
      const newNamespaces = new Set<string>()

      // Grab namespaces from pods
      data.forEach((p: Pod) => {
        newNamespaces.add(p.metadata.namespace)
      })

      // Update stores
      pods.set(data)
      namespaces.set(new Set(Array.from(newNamespaces).sort()))

      // Set from URL params after data is loaded
      if (podParam && nsParam) {
        selectedNamespace.set(nsParam)
        podName.set(podParam)
      }
    }

    podEventSource.onerror = (event) => {
      console.error('EventSource failed:', event)
    }

    // add window event handler for keydown (used to track esc)
    const handleKeydown = (e: KeyboardEvent) => {
      switch (e.key) {
        // If the Escape key is pressed, close the Logs view
        case 'Escape':
          goto('/workloads/pods')
          return
      }
    }
    window.addEventListener('keydown', handleKeydown)

    return () => {
      window.removeEventListener('keydown', handleKeydown)
    }
  })

  function cleanupPage() {
    // clean up event sources
    if (podEventSource) {
      podEventSource.close()
    }
    const logES = get(logEventSource)
    if (logES) {
      logES.close()
      logEventSource.set(null)
    }

    // clean up timeouts
    if (notificationTimeout) {
      clearTimeout(notificationTimeout)
    }
    const st = get(scrollTimeout)
    if (st) {
      clearTimeout(st)
      scrollTimeout.set(null)
    }
  }
</script>

<section class="h-full flex-shrink-0 mt-1" data-testid="logs-view">
  <div class="flex h-full w-full flex-col bg-gray-900 rounded-lg overflow-y-auto">
    <!-- Header -->
    <div class="flex items-center justify-between py-3 border-b border-gray-700 mx-4">
      <!-- Left section: title, search, and log controls -->
      <div class="flex items-center space-x-6">
        <!-- Logo and title -->
        <div class="flex items-center">
          <LogsIcon class="w-5 h-5" />
          <span class="font-medium text-gray-100 text-nowrap ml-1">Logs</span>
        </div>

        <!-- Search -->
        <div class="relative w-72">
          <div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
            <Search class="text-gray-100" />
          </div>
          <input
            type="text"
            name="input-search"
            data-testid="logs-search"
            autocomplete="off"
            class="h-9 block w-full rounded-md border-gray-700 bg-gray-800 p-2.5 pl-9 text-gray-100 placeholder-gray-300 focus:ring-primary-500 focus:border-primary-500 focus:ring-blue-600 text-sm"
            placeholder="Search"
            bind:value={$searchTerm}
          />
        </div>

        <!-- Autoscroll toggle -->
        <div class="flex items-center">
          <label class="relative inline-flex items-center cursor-pointer">
            <input type="checkbox" bind:checked={$autoScroll} data-testid="autoscroll-btn" class="sr-only peer" />
            <div
              class="w-9 h-5 bg-gray-800 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-800 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-blue-600"
            />
            <span class="ml-2 text-sm font-medium text-gray-100">Auto-scroll</span>
          </label>
        </div>
      </div>

      <!-- Right section: Action buttons -->
      <div class="flex items-center space-x-2">
        <div class="relative">
          <button
            class="p-1.5 text-gray-100 hover:text-gray-100 rounded-lg hover:bg-gray-800"
            data-testid="copy-logs-btn"
            on:click={handleCopy}
          >
            <Copy class="w-4 h-4" />
          </button>
          {#if showNotification}
            <div
              transition:fade={{ duration: 200 }}
              data-testid="logs-notification"
              class="absolute left-1/2 -translate-x-1/2 top-full mt-2 z-50 flex items-center gap-2 px-3 py-2 text-sm text-gray-100 bg-gray-900 rounded-lg shadow-lg border border-gray-800 whitespace-nowrap"
            >
              <CheckmarkOutline class="w-4 h-4" />
              <span>{notificationMessage}</span>
            </div>
          {/if}
        </div>
        <div class="relative">
          <button
            class="p-1.5 text-gray-100 hover:text-gray-100 rounded-lg hover:bg-gray-800"
            data-testid="download-logs-btn"
            on:click={handleDownload}
          >
            <Download class="w-4 h-4" />
          </button>
        </div>
        <div class="relative">
          <button
            type="button"
            class="text-gray-100 bg-transparent hover:bg-gray-200 hover:text-gray-900 rounded-lg text-sm w-8 h-8 top-2.5 end-2.5 inline-flex items-center justify-center dark:hover:bg-gray-600 dark:hover:text-gray-100"
            data-testid="close-logs-btn"
            on:click={() => goto('/workloads/pods')}
          >
            <Close class="w-5 h-5" />
          </button>
        </div>
      </div>
    </div>

    <!-- Dropdowns -->
    <div class="flex space-x-4 p-4">
      <Dropdown
        bind:selectedValue={$selectedNamespace}
        items={Array.from($namespaces)}
        id="ns-dropdown"
        width="w-72"
        placeholder="Select namespace"
        label="Namespace"
      />
      <Dropdown
        bind:selectedValue={$podName}
        items={$filteredPods.map((p) => p.metadata.name)}
        id="pod-dropdown"
        width="w-72"
        placeholder="Select pod"
        disabled={!$selectedNamespace}
        label="Pod"
      />
      <Dropdown
        bind:selectedValue={$containerName}
        items={$filteredContainers.map((c) => c.name)}
        id="container-dropdown"
        width="w-72"
        placeholder="Select container"
        disabled={!$podName}
        label="Container"
      />
    </div>
    <!-- Logs content area -->
    <div class="flex-1 m-1 min-h-12">
      <div class="h-full bg-gray-900 rounded-md p-1">
        <div class="h-full overflow-y-auto" data-testid="ansi-display">
          {#if !$containerName}
            <HintMessage primaryMsg="No container selected" secondaryMsg="Select a container to view its logs" />
          {:else if $containerName && $logElements.length === 0}
            <HintMessage primaryMsg="No logs found" secondaryMsg="The selected container hasn't generated any logs" />
          {/if}
          <div id="scroll-anchor" class="h-px" />
        </div>
      </div>
    </div>
  </div>
</section>
