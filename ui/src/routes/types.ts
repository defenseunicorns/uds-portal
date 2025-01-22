// Copyright 2025 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

export type Metadata = {
  name: string
}

export type Status = {
  endpoints: string[]
}

export type App = {
  metadata: Metadata
  status: Status
  icon: ConstructorOfATypedSvelteComponent
}
