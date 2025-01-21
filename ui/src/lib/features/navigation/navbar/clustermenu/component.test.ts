import { fireEvent, render, waitFor } from '@testing-library/svelte'
import { vi } from 'vitest'

import '@testing-library/jest-dom'

import { goto } from '$app/navigation'

import ClustersComponent from './component.svelte'
import { clusters, loadingCluster, type ClusterInfo } from './store'

vi.mock('./store', async (importActual) => {
  const actual: Record<string, unknown> = await importActual()
  return {
    ...actual,
    loadingCluster: {
      set: vi.fn(),
      subscribe: vi.fn(),
      update: vi.fn(),
    },
  }
})

vi.mock('$app/navigation', () => ({
  goto: vi.fn(),
}))

vi.mock('$app/stores', () => {
  const unsubscribe = vi.fn()
  return {
    page: {
      subscribe: () => vi.fn().mockReturnValue(unsubscribe),
    },
  }
})

const mockClusters: ClusterInfo[] = [
  { name: 'Cluster 1', context: 'ctx-1', selected: true },
  { name: 'Cluster 2', context: 'ctx-2', selected: false },
  { name: 'Cluster 3', context: 'ctx-3', selected: false },
]

describe('ClustersComponent', () => {
  const fetchMock = vi.fn()
  global.fetch = fetchMock

  clusters.set(mockClusters)

  const cluster1 = `${mockClusters[0].context}.${mockClusters[0].name}`
  const cluster2 = `${mockClusters[1].context}.${mockClusters[1].name}`
  const cluster3 = `${mockClusters[2].context}.${mockClusters[2].name}`

  beforeEach(() => {
    vi.resetAllMocks()
    vi.useFakeTimers()

    fetchMock.mockResolvedValueOnce({
      json: async () => mockClusters,
      ok: true,
    })
  })

  afterAll(() => {
    vi.restoreAllMocks()
    vi.useRealTimers()
  })

  test('renders the selected cluster', async () => {
    const { getByText, getByTestId } = render(ClustersComponent)
    expect(getByText(cluster1)).toBeInTheDocument()

    const dropdown = getByTestId('clusterDropdown')
    expect(dropdown.classList.contains('hidden')).toBe(true)
  })

  test('renders the dropdown with all clusters', async () => {
    const { getByRole } = render(ClustersComponent)

    const button = getByRole('button', { name: cluster1 })
    await fireEvent.click(button)

    expect(fetchMock).toHaveBeenCalledWith('/api/v1/clusters', {
      method: 'GET',
    })
  })

  test('search in dropdown', async () => {
    const { getByRole, queryByRole, getByPlaceholderText } = render(ClustersComponent)

    const button = getByRole('button', { name: cluster1 })
    await fireEvent.click(button)

    const search = getByPlaceholderText('Search')
    await fireEvent.input(search, { target: { value: cluster3 } })

    expect((search as HTMLInputElement).value).toBe(cluster3)

    expect(getByRole('button', { name: cluster3 })).toBeInTheDocument()
    expect(queryByRole('button', { name: cluster2 })).not.toBeInTheDocument()
  })

  test('changes the selected cluster on click', async () => {
    const { getByText, getByRole } = render(ClustersComponent)

    const button = getByRole('button', { name: cluster1 })
    await fireEvent.click(button)

    fetchMock.mockResolvedValueOnce({ ok: true }).mockResolvedValueOnce({
      json: async () => mockClusters,
      ok: true,
    })

    await fireEvent.click(getByText(cluster2))
    vi.advanceTimersByTime(3000)

    expect(loadingCluster.set).toHaveBeenCalledTimes(2)

    expect(loadingCluster.set).toHaveBeenNthCalledWith(1, { loading: true, cluster: mockClusters[1] })
    expect(loadingCluster.set).toHaveBeenCalledWith({ loading: false, cluster: mockClusters[1] })

    // We made POST request to change the cluster
    expect(fetchMock).toHaveBeenNthCalledWith(2, '/api/v1/cluster', {
      method: 'POST',
      body: JSON.stringify({ cluster: mockClusters[1] }),
    })

    // Then got updated cluster list since selected has changed
    expect(fetchMock).toHaveBeenLastCalledWith('/api/v1/clusters', {
      method: 'GET',
    })

    await waitFor(() => expect(goto).toHaveBeenCalledWith('/'))
  })

  test('timeout is reached and fetch response has not returned', async () => {
    const { getByText, getByRole } = render(ClustersComponent)

    const button = getByRole('button', { name: cluster1 })
    await fireEvent.click(button)

    vi.clearAllMocks()

    fetchMock
      .mockImplementationOnce(() => {
        return new Promise((resolve) => {
          setTimeout(() => {
            resolve({ ok: true })
          }, 5000)
        })
      })
      .mockImplementationOnce(() => {
        return Promise.resolve({
          json: async () => mockClusters,
          ok: true,
        })
      })

    await fireEvent.click(getByText(cluster2))
    expect(loadingCluster.set).toHaveBeenNthCalledWith(1, { loading: true, cluster: mockClusters[1] })

    vi.advanceTimersByTime(5000)
    vi.runOnlyPendingTimersAsync()
    await waitFor(() => {
      expect(fetchMock).toHaveBeenLastCalledWith('/api/v1/clusters', {
        method: 'GET',
      })
      expect(loadingCluster.set).toHaveBeenCalledWith({ loading: false, cluster: mockClusters[1] })
      expect(goto).toHaveBeenCalledWith('/')
    })
  })

  test('handles fetch error', async () => {
    const { getByText, getByRole } = render(ClustersComponent)
    vi.useRealTimers()
    const button = getByRole('button', { name: cluster1 })
    await fireEvent.click(button)

    fetchMock
      .mockResolvedValueOnce({ ok: false, text: () => Promise.resolve('Error') })
      .mockResolvedValueOnce({ ok: true, json: async () => mockClusters })

    await fireEvent.click(getByText(cluster2))
    await waitFor(() => {
      expect(loadingCluster.set).toHaveBeenCalledTimes(1)
      expect(goto).not.toHaveBeenCalled()
      expect(loadingCluster.update).toHaveBeenCalled()
    })
  })
})
