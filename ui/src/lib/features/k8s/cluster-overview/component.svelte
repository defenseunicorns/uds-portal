<!-- Copyright 2024 Defense Unicorns -->
<!-- SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial -->

<script lang="ts">
  import { onMount } from 'svelte'

  import type { V1Pod } from '@kubernetes/client-node'
  import { goto } from '$app/navigation'
  import {
    ApplicationsWidget,
    CoreServicesWidget,
    CVEReportWidget,
    HeaderWithIcon,
    InactiveBadge,
    ProgressBarWidget,
    WithRightIconWidget,
  } from '$components'
  import EventsOverviewWidget from '$components/k8s/Event/EventsOverviewWidget.svelte'
  import { createStore } from '$lib/features/k8s/events/store'
  import { type ClusterOverviewUDSPackageType } from '$lib/types'
  import { resourceDescriptions } from '$lib/utils/descriptions'
  import { Analytics, ChevronRight, DataVis_1 } from 'carbon-icons-svelte'
  import Chart from 'chart.js/auto'

  import { calculatePercentage, formatTime, mebibytesToGigabytes, millicoresToCores } from '../helpers'
  import type { ClusterData } from '../types'
  import { getChartData, getChartOptions, renderCustomLegend } from './chart'

  import './styles.postcss'

  let clusterData: ClusterData = {
    totalPods: 0,
    totalNodes: 0,
    cpuCapacity: 0,
    memoryCapacity: 0,
    currentUsage: {
      CPU: 0,
      Memory: 0,
      Timestamp: new Date().toISOString(),
    },
    historicalUsage: [],
    packagesData: [],
    podsData: [],
  }

  let udsPackages: ClusterOverviewUDSPackageType[] = []
  let pods: V1Pod[] = []
  const description = resourceDescriptions['Events']

  let myChart: Chart | null = null
  let canvasElement: HTMLCanvasElement

  let metricsServerAvailable = true
  let metricsServerNewlyAvailable = false
  let previousMetricsServerState = metricsServerAvailable

  $: metricsServerAvailable = !(clusterData.currentUsage.CPU === -1 && clusterData.currentUsage.Memory === -1)
  $: {
    metricsServerNewlyAvailable = !previousMetricsServerState && metricsServerAvailable
    previousMetricsServerState = metricsServerAvailable
  }

  $: cpuPercentage = calculatePercentage(clusterData.currentUsage.CPU, clusterData.cpuCapacity)
  $: memoryPercentage = calculatePercentage(clusterData.currentUsage.Memory, clusterData.memoryCapacity)
  $: gbUsed = mebibytesToGigabytes(clusterData.currentUsage.Memory)
  $: gbCapacity = mebibytesToGigabytes(clusterData.memoryCapacity)
  $: cpuUsed = millicoresToCores(clusterData.currentUsage.CPU)
  $: formattedCpuCapacity = millicoresToCores(clusterData.cpuCapacity)

  $: if (canvasElement && !myChart) {
    const options = getChartOptions(metricsServerAvailable)

    myChart = new Chart(canvasElement, {
      type: 'line',
      data: getChartData(metricsServerAvailable),
      options: {
        ...options,
        scales: {
          ...options.scales,
          y: { display: false },
          y1: { display: false },
          x: { ticks: { color: 'white', maxTicksLimit: 7 } },
        },
        plugins: {
          legend: {
            display: false,
          },
        },
      },
      plugins: [
        {
          id: 'customLegend',
          afterRender: renderCustomLegend,
        },
      ],
    })
  }

  $: if (myChart && clusterData.historicalUsage) {
    if (metricsServerNewlyAvailable) {
      myChart.data = getChartData(metricsServerAvailable)
      myChart.options = getChartOptions(metricsServerAvailable)
    }

    myChart.data.labels = clusterData.historicalUsage.map((point) => [formatTime(point.Timestamp)])

    if (metricsServerAvailable) {
      myChart.data.datasets[0].data = clusterData.historicalUsage.map((point) => point.Memory / (1024 * 1024 * 1024))
      myChart.data.datasets[1].data = clusterData.historicalUsage.map((point) => point.CPU / 1000)
    } else {
      myChart.data = getChartData(false)
      myChart.options = getChartOptions(false)
    }

    myChart.update()
  }

  function updateClusterData(newData: Partial<ClusterData>) {
    clusterData = {
      ...clusterData,
      ...newData,
    }
  }

  onMount(() => {
    const overviewPath: string = '/api/v1/monitor/cluster-overview'
    const overview = new EventSource(overviewPath)

    overview.addEventListener('usageDataUpdate', function (event) {
      const newData = JSON.parse(event.data) as ClusterData
      if (newData && Object.keys(newData).length > 0) {
        updateClusterData({
          totalNodes: newData.totalNodes,
          totalPods: newData.totalPods,
          cpuCapacity: newData.cpuCapacity,
          memoryCapacity: newData.memoryCapacity,
          currentUsage: newData.currentUsage,
          historicalUsage: newData.historicalUsage,
        })
      }
    })

    overview.addEventListener('packagesDataUpdate', function (event) {
      const newData = JSON.parse(event.data) as ClusterData
      udsPackages = [...newData.packagesData]
      updateClusterData({ packagesData: newData.packagesData })
    })

    overview.addEventListener('podsDataUpdate', function (event) {
      const newData = JSON.parse(event.data) as ClusterData
      pods = [...newData.podsData]
      updateClusterData({ podsData: newData.podsData })
    })

    Chart.register({})

    return () => {
      overview.close()
      if (myChart) {
        myChart.destroy()
        myChart = null
      }
    }
  })

  Chart.defaults.datasets.line.tension = 0.4
</script>

<div class="p-4 dark:text-gray-100 pt-0 space-y-8">
  <h1 class="text-xl font-bold mb-4">Cluster Overview</h1>

  <div class="grid grid-cols-2 min-[1510px]:grid-cols-4 gap-4">
    <WithRightIconWidget
      statText={`${clusterData.totalPods}`}
      helperText="Pods Running"
      icon={Analytics}
      link="/workloads/pods"
    />

    <WithRightIconWidget
      statText={`${udsPackages.length}`}
      helperText="Packages Deployed"
      icon={DataVis_1}
      link="/configs/uds-packages"
    />

    <ProgressBarWidget
      capacity={formattedCpuCapacity}
      progress={cpuUsed}
      statText="CPU Usage"
      unit="Cores"
      value={metricsServerAvailable ? `${cpuPercentage.toFixed(2)}%` : 'Unavailable'}
      deactivated={!metricsServerAvailable}
    />

    <ProgressBarWidget
      capacity={gbCapacity}
      progress={gbUsed}
      statText="Memory Usage"
      unit="GB"
      value={metricsServerAvailable ? `${memoryPercentage.toFixed(2)}%` : 'Unavailable'}
      deactivated={!metricsServerAvailable}
    />
  </div>

  <div class="grid grid-cols-1 gap-8 xl:grid-cols-2 xl:gap-4">
    <div class="px-2 py-5 bg-gray-900 rounded-lg overflow-hidden shadow h-full">
      <ApplicationsWidget {udsPackages} />
    </div>

    <div class="px-2 py-5 bg-gray-900 rounded-lg overflow-hidden shadow h-full">
      <CoreServicesWidget coreServices={udsPackages} {pods} />
    </div>
  </div>

  <div class="grid grid-cols-1 xl:grid-cols-2 gap-8 xl:gap-4">
    <div class="px-2 py-5 bg-gray-900 rounded-lg overflow-hidden shadow h-full">
      <CVEReportWidget />
    </div>

    <div class="py-5 px-8 bg-gray-900 rounded-lg overflow-hidden shadow h-full">
      <div class="">
        <HeaderWithIcon title="Resource Usage" description="Resource Usage report" />
      </div>

      <div style:position="relative" style:margin="auto">
        <canvas bind:this={canvasElement} height={335} />

        {#if !metricsServerAvailable}
          <div class="flex items-center">
            <div class="absolute z-40 group ml-2 flex items-center -top-6 right-5">
              <InactiveBadge
                tooltipDirection="tooltip-left"
                tooltipText="Metrics Server is unavailable.
            Ensure Metrics Server is running in the cluster."
              />
            </div>
          </div>
        {:else}
          <div id="custom-legend" style="position: absolute; top: -20px; right: 0; display: flex" />
        {/if}
      </div>

      <div class="dark:bg-gray-900">
        <div class="bg-white dark:bg-gray-900 flex items-center justify-end pt-3 dark:border-gray-700">
          <button
            disabled={!metricsServerAvailable}
            class="text-sm {!metricsServerAvailable
              ? 'dark:text-gray-400'
              : 'dark:text-blue-300'} text-blue-500 flex items-center space-x-1"
            on:click={() => goto('/resource-usage')}
          >
            <span>VIEW REPORT</span>
            <ChevronRight />
          </button>
        </div>
      </div>
    </div>
  </div>

  <EventsOverviewWidget {createStore} {description} />
</div>
