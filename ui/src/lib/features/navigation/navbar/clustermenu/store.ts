import { writable, type Writable } from 'svelte/store'

export type ClusterInfo = {
  name: string
  context: string
  selected: boolean
}

type Loading = {
  loading: boolean
  cluster: ClusterInfo
  err?: string
}

export const clusters: Writable<ClusterInfo[]> = writable([])

export const loadingCluster: Writable<Loading> = writable({
  loading: false,
  cluster: { name: '', context: '', selected: false },
})
