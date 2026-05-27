// Copyright 2025-2026 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

export type ApiApp = {
  name: string
  icon?: string
  url: string
  gateway?: string
}

export type ClassBannerCfg = {
  enabled: boolean
  text: string
  footer: boolean
}
