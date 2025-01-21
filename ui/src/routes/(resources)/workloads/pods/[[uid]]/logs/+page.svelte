<!-- Copyright 2024 Defense Unicorns -->
<!-- SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial -->

<script lang="ts">
  import { Logs } from '$components'
  import { PodTable } from '$features/k8s'

  let podTableHeight = 70 // Initial height percentage of the PodTable
  $: logsHeight = 90 - podTableHeight
  let isDragging = false
  let initialMouseY = 0
  let initialHeight = 0

  const startDrag = (e: MouseEvent) => {
    e.preventDefault()
    e.stopPropagation()
    isDragging = true
    initialMouseY = e.clientY
    initialHeight = podTableHeight

    const onMouseMove = (e: MouseEvent) => {
      if (isDragging) {
        const deltaY = e.clientY - initialMouseY
        const containerHeight = document.documentElement.clientHeight
        podTableHeight = Math.max(10, Math.min(90, initialHeight + (deltaY / containerHeight) * 100)) // Restrict between 10% and 90%
      }
    }

    const stopDrag = () => {
      isDragging = false
      window.removeEventListener('mousemove', onMouseMove)
      window.removeEventListener('mouseup', stopDrag)
    }

    window.addEventListener('mousemove', onMouseMove)
    window.addEventListener('mouseup', stopDrag)
  }
</script>

<div class="flex h-full w-full flex-col" class:is-dragging={isDragging}>
  <!-- PodTable Section -->
  <div class="pod-table" style="--pod-table-height: {podTableHeight}%; pointer-events: {isDragging ? 'none' : 'auto'};">
    <PodTable />
  </div>

  <!-- Resize Bar -->
  <button
    class="resize-bar rounded-md cursor-row-resize"
    on:mousedown={startDrag}
    class:is-dragging={isDragging}
    data-testid="resize-bar"
  ></button>

  <!-- Logs Section -->
  <div class="logs flex-shrink-0 flex-grow" style="--logs-height: {logsHeight}%">
    <Logs />
  </div>
</div>

<style>
  .pod-table {
    height: var(--pod-table-height, 70%);
    transition: height ease;
  }

  .logs {
    transition: height ease;
    height: var(--logs-height, 30%);
  }

  .resize-bar {
    height: 5px;
    background: transparent;
    margin-top: 5px;
    user-select: none;
  }

  .resize-bar:hover {
    background: rgb(118, 123, 131);
  }

  .resize-bar.is-dragging {
    pointer-events: none;
  }

  .is-dragging {
    cursor: row-resize;
  }
</style>
