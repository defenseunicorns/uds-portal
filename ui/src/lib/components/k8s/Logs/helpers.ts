// Copyright 2024 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

import { get } from 'svelte/store'

import { autoScroll, logElements, logEventSource, scrollTimeout, searchTerm } from '$components/k8s/Logs/store'
import Convert from 'ansi-to-html'

const convert = new Convert({
  newline: true,
  stream: true,
  escapeXML: true,
  colors: {
    0: '#374151',
    1: '#F3F4F6',
    2: '#D1D5DB',
  },
})

// removes highlights from logs and resets to plain text
export const resetHighlights = () => {
  const logEls = get(logElements)
  logEls.forEach((element) => {
    // Remove all <mark> tags and restore their content
    element.innerHTML = element.innerHTML.replace(/<mark class="[^"]*">(.*?)<\/mark>/g, '$1')
  })
  logElements.set(logEls)
}

// removes log lines from DOM and resets logs store
export const cleanupLogs = () => {
  const logES = get(logEventSource)
  if (logES) {
    logES.close()
    logEventSource.set(null)
    const logLines = document.querySelectorAll('.log-line')
    logLines.forEach((line) => line.remove())
    logElements.set([])
  }
  autoScroll.set(true)
}

// highlights search term matches in logs
export const highlightMatches = (term: string) => {
  if (!term.trim()) return

  // First remove any existing highlights
  resetHighlights()

  // Then add new highlights
  const logEls = get(logElements)
  logEls.forEach((element) => {
    element.innerHTML = element.innerHTML.replace(
      new RegExp(term, 'gi'), // global and case-insensitive
      (match) => `<mark class="bg-yellow-500 text-gray-900 px-0.5 rounded">${match}</mark>`,
    )
  })
  logElements.set(logEls)
}

// adds logs to the DOM
export const addLog = (message: string) => {
  const scrollAnchor = document.getElementById('scroll-anchor')
  if (!message || !scrollAnchor?.parentElement) return

  const html = convert.toHtml(message)
  const lineElement = document.createElement('div')
  lineElement.className = 'log-line text-gray-100 text-sm font-mono py-0.5'
  lineElement.innerHTML = html

  const logEls = get(logElements)
  window.requestAnimationFrame(() => {
    scrollAnchor.parentElement?.insertBefore(lineElement, scrollAnchor)
    logEls.push(lineElement)
    // Keep only last 1000 lines in memory to prevent memory leaks
    if (logEls.length > 1000) {
      logEls.shift()
    }
    logElements.set(logEls)

    // get search term and highlight matches in new log
    const search = get(searchTerm)
    if (search) {
      highlightMatches(search)
    }

    debounceScroll()
  })
}

// Debounce auto-scrolling, if many logs come in at once, this creates a smooth scroll UX
export const debounceScroll = () => {
  if (!get(autoScroll)) return

  // if a new log arrives during the timeout, clear it and start over
  // this prevents the scroll from jumping around
  let timeout = get(scrollTimeout)
  if (timeout) {
    window.clearTimeout(timeout)
  }
  timeout = window.setTimeout(() => {
    const scrollAnchor = document.getElementById('scroll-anchor')
    scrollAnchor?.scrollIntoView({ behavior: 'auto' })
  }, 100)
  scrollTimeout.set(timeout)
}

// Get last N lines of logs (used for getting logs in plaintext for copying)
export const getLastNLines = (n: number): string => {
  const logEls = get(logElements)
  const lastN = logEls.slice(-n)
  return lastN
    .map((el) => {
      // Convert HTML back to plain text and remove ANSI codes
      const text = el.textContent || ''
      return text.replace(/\n$/, '') // Remove trailing newline if present
    })
    .join('\n')
}

// Function to get all logs for downloading
export const getAllLogs = (): string => {
  return get(logElements)
    .map((el) => {
      const text = el.textContent || ''
      return text.replace(/\n$/, '')
    })
    .join('\n')
}
