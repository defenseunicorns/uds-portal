<script lang="ts">
  import { onMount } from 'svelte'

  import { InactiveBadge } from '$components'
  import { getChartData, getChartOptions } from '$features/k8s/cluster-overview/chart'
  import { formatTime } from '$features/k8s/helpers'
  import type { ClusterData } from '$features/k8s/types'
  import Chart from 'chart.js/auto'

  let canvasElement: HTMLCanvasElement
  let myChart: Chart | null = null
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

  let metricsServerAvailable = true
  let metricsServerNewlyAvailable = false
  let previousMetricsServerState = metricsServerAvailable

  $: metricsServerAvailable = !(clusterData.currentUsage.CPU === -1 && clusterData.currentUsage.Memory === -1)
  $: {
    metricsServerNewlyAvailable = !previousMetricsServerState && metricsServerAvailable
    previousMetricsServerState = metricsServerAvailable
  }

  $: if (canvasElement && !myChart) {
    myChart = new Chart(canvasElement, {
      type: 'line',
      data: getChartData(metricsServerAvailable),
      options: getChartOptions(metricsServerAvailable),
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

  onMount(() => {
    const overviewPath: string = '/api/v1/monitor/cluster-overview'
    const overview = new EventSource(overviewPath)

    overview.addEventListener('usageDataUpdate', function (event) {
      const newData = JSON.parse(event.data) as ClusterData
      if (newData && Object.keys(newData).length > 0) {
        clusterData = {
          ...clusterData,
          totalNodes: newData.totalNodes,
          totalPods: newData.totalPods,
          cpuCapacity: newData.cpuCapacity,
          memoryCapacity: newData.memoryCapacity,
          currentUsage: newData.currentUsage,
          historicalUsage: newData.historicalUsage,
        }
      }
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

<div class="relative">
  {#if !metricsServerAvailable}
    <div class="absolute z-50 group ml-2 flex items-center top-20 right-5">
      <InactiveBadge
        tooltipDirection="tooltip-left"
        tooltipText="Metrics Server is unavailable.
            Ensure Metrics Server is running in the cluster."
      />
    </div>
  {/if}

  <h1 class="text-xl font-bold mb-6">Resource Usage</h1>

  <div
    class="p-5 {metricsServerAvailable ? 'pt-10' : 'pt-20'} bg-gray-900 rounded-lg overflow-hidden shadow"
    style:position="relative"
  >
    <canvas bind:this={canvasElement} height={500} />
  </div>
</div>
