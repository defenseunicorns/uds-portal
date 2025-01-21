// Copyright 2024 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

import { derived, get, writable } from 'svelte/store'

import { addLog, cleanupLogs, debounceScroll, highlightMatches, resetHighlights } from '$components/k8s/Logs/helpers'
import type { Pod } from '$features/k8s/types'

// stores for pod data
export const pods = writable<Pod[]>([]) // raw pod data
export const selectedNamespace = writable<string | null>(null)
export const namespaces = writable<Set<string>>(new Set()) // derived from pods
export const podName = writable<string | null>(null)
export const selectedPod = writable<Pod | null>(null) // used to derive containers
export const containerName = writable<string | null>(null)

// scroll and search stores
export const autoScroll = writable(true)
export const scrollTimeout = writable<number | null>(null) // used for smooth scroll debouncing
export const searchTerm = writable('')

// log store
export const logElements = writable<HTMLDivElement[]>([])

// Derive filtered lists for pod and container dropdowns
export const filteredPods = derived([selectedNamespace, pods], ([$selectedNamespace, $pods]) =>
  $selectedNamespace ? Array.from($pods).filter((p) => p.metadata.namespace === $selectedNamespace) : [],
)
export const filteredContainers = derived([selectedPod], ($selectedPod) => $selectedPod[0]?.spec.containers ?? [])

// EventSource store for logs
export const logEventSource = writable<EventSource | null>()

// handle search term changes
export function handleSearch() {
  searchTerm.subscribe((term) => {
    if (!term) {
      resetHighlights()
      return
    }
    highlightMatches(term)
  })
}

// handle autoscroll button toggles
export function handleAutoscroll() {
  autoScroll.subscribe((scroll) => {
    // scroll down when the button is toggled on
    if (scroll) {
      debounceScroll()
    }
  })
}

// handle namespace updates
export function handleNamespace() {
  selectedNamespace.subscribe((newNamespace) => {
    if (newNamespace) {
      cleanupLogs()
      podName.set(null)
      selectedPod.set(null)
      containerName.set(null)
    }
  })
}

// handle pod updates
export function handlePod() {
  podName.subscribe((newPodName) => {
    if (newPodName) {
      cleanupLogs()
      containerName.set(null)
      const pods = get(filteredPods)
      const pod = pods.find((p) => p.metadata.name === newPodName)
      selectedPod.set(pod ?? null)

      // if only one container, select it
      if (pod?.spec.containers.length === 1) {
        containerName.set(pod.spec.containers[0].name)
      }
    }
  })
}

// handle container updates, sets up logs SSE stream
export function handleContainer() {
  containerName.subscribe((newContainerName) => {
    if (newContainerName) {
      cleanupLogs()
      const ns = get(selectedNamespace)
      const pod = get(podName)

      // create new SSE stream for logs
      const logsURL = `/api/v1/resources/workloads/pods/logs?namespace=${ns}&pod=${pod}&container=${newContainerName}`
      const es = new EventSource(logsURL)

      es.onmessage = (event) => {
        addLog(event.data)
      }

      es.onerror = (event) => {
        console.error('EventSource failed:', event)
      }
      logEventSource.set(es)
    }
  })
}
