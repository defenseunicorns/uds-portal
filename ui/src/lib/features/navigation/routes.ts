// Copyright 2024 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

import { CRDIcon } from '$lib/icons'
import {
  AudioConsole,
  ChartAverage,
  ChartCombo,
  DataTableReference,
  Db2Database,
  KubernetesPod,
  Layers,
  Network_2,
  SearchLocate,
  Security,
  Share,
  WorkflowAutomation,
} from 'carbon-icons-svelte'

import type { BaseRoute, Route } from './types'

const baseRoutes: BaseRoute[] = [
  {
    name: 'Overview',
    icon: ChartCombo,
    path: '/',
  },
  {
    name: 'Applications',
    icon: DataTableReference,
    children: ['Packages', 'Endpoints'],
  },
  {
    name: 'Monitor',
    icon: SearchLocate,
    children: ['Pepr', 'Events'],
  },
  {
    name: 'Security',
    icon: Security,
    children: ['CVE Report'],
  },
  {
    name: 'Workloads',
    icon: KubernetesPod,
    class: 'top-border',
    children: ['Pods', 'Deployments', 'DaemonSets', 'StatefulSets', 'Jobs', 'CronJobs'],
  },
  {
    name: 'Configs',
    icon: AudioConsole,
    class: 'icon-rotate-90',
    children: ['UDS Packages', 'UDS Exemptions', 'ConfigMaps', 'Secrets'],
  },
  {
    name: 'Cluster Ops',
    icon: WorkflowAutomation,
    children: [
      'Mutating Webhooks',
      'Validating Webhooks',
      'HPA',
      'Pod Disruption Budgets',
      'Resource Quotas',
      'Limit Ranges',
      'Priority Classes',
      'Runtime Classes',
    ],
  },
  {
    name: 'Networks',
    icon: Network_2,
    children: ['Services', 'Virtual Services', 'Network Policies', 'Endpoints', 'Gateways'],
  },
  {
    name: 'Storage',
    icon: Db2Database,
    children: ['Persistent Volumes', 'Persistent Volume Claims', 'Storage Classes'],
  },
  {
    name: 'Namespaces',
    icon: Layers,
  },
  {
    name: 'Nodes',
    icon: Share,
  },
  {
    name: 'Custom Resources',
    icon: CRDIcon,
  },
  {
    name: 'Resource Usage',
    icon: ChartAverage,
  },
]

// Convert the path to a URL-friendly format
const createPath = (name: string) => `/${name.replace(/\s+/g, '-').toLowerCase()}`

// Convert the base routes to routes
export const routes: Route[] = baseRoutes.map(({ name, children, path, ...rest }) => ({
  ...rest,
  name,
  path: path || createPath(name),
  children: children?.map((name) => ({ name, path: createPath(name) })),
}))
