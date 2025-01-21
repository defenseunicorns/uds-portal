// Copyright 2024 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

import '@testing-library/jest-dom'

import type { DrawerOpts } from '$components/k8s/Drawer/types'
import { testK8sTableWithCustomColumns, testK8sTableWithDefaults } from '$features/k8s/test-helper'

import Component from './component.svelte'
import { createStore } from './store'

suite('Application packages Component', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  const name = 'Application Packages'
  const description = 'Packages are either UDS or Zarf packages that are deployed in the cluster.'
  const minWidth = 1250
  const drawerOpts: DrawerOpts = {
    includeEvents: false,
    addDetails: () => [],
    transformResource: () => {},
  }

  testK8sTableWithDefaults(Component, {
    createStore,
    columns: [
      ['name', 'w-2/12'],
      ['version', ' w-2/12'],
      ['description', 'italic w-4/12 line-clamp-2'],
      ['arch', ' w-1/12'],
      ['flavor', ' w-2/12'],
      ['age', ' w-1/12'],
    ],
    name,
    description,
    drawerOpts,
    isNamespaced: false,
    minWidth,
  })

  testK8sTableWithCustomColumns(Component, {
    createStore,
    name,
    description,
    isNamespaced: false,
    drawerOpts,
    minWidth,
  })
})
