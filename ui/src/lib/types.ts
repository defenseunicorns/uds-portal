// Copyright 2025 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

import type { UDSPackageStatus } from '$features/k8s/types'

export type ClusterOverviewUDSPackageType = {
  metadata: {
    name: string
  }
  status: {
    phase: UDSPackageStatus
    endpoints: string[]
  }
}

export type ClassBannerCfg = {
  enabled: boolean
  text: string
  footer: boolean
}
