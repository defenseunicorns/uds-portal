// Copyright 2024 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

import type { KubernetesObject } from '@kubernetes/client-node'
import type { DeployedPackage } from '$features/k8s/applications/packages/types'

export type DrawerDetails = Array<{ label: string; value: string | Array<string>; isList?: boolean }>

export type DrawerOpts = {
  // include the events tab in the details panel
  includeEvents?: boolean

  // use custom details in details tab
  addDetails?: (resource: KubernetesObject | DeployedPackage) => DrawerDetails

  // transform underlying resource before processing (this will occur before addDetails)
  transformResource?: (resource: KubernetesObject) => unknown
}
