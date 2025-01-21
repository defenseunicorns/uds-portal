import type { ClusterInfo } from './store'

export const noClustersMsg = 'No other clusters found in KUBECONFIG'

export function displayClusterName(cluster: ClusterInfo) {
  return cluster.context === cluster.name ? cluster.name : `${cluster.context}.${cluster.name}`
}

export function calcClusterDropdownWidth(availableClusters: ClusterInfo[]) {
  const paddingAdjustment = 50
  if (availableClusters.length === 0) {
    return getWidthOfOpt(noClustersMsg) + paddingAdjustment
  }

  const longestOption = availableClusters.reduce(
    (a: ClusterInfo, b: ClusterInfo) => (displayClusterName(a).length > displayClusterName(b).length ? a : b),
    availableClusters[0],
  )
  return getWidthOfOpt(displayClusterName(longestOption))
}

function getWidthOfOpt(opt: string) {
  const span = document.createElement('span')
  span.style.font = 'inherit' // Ensure it uses the same font as the dropdown
  span.style.visibility = 'hidden'
  span.style.whiteSpace = 'nowrap'
  span.textContent = opt
  document.body.appendChild(span)
  const width = span.offsetWidth
  document.body.removeChild(span)
  return width
}
