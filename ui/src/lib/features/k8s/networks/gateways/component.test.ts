// Copyright 2024 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

import '@testing-library/jest-dom'

import { testK8sTableWithCustomColumns, testK8sTableWithDefaults } from '$features/k8s/test-helper'
import { resourceDescriptions } from '$lib/utils/descriptions'

import Component from './component.svelte'
import { createStore } from './store'

suite('Gateways Component', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  const name = 'Gateways'
  const description = resourceDescriptions[name]
  const minWidth = 800

  testK8sTableWithDefaults(Component, {
    createStore,
    columns: [
      ['name', 'w-3/12'],
      ['namespace', 'w-4/12'],
      ['hosts', 'w-4/12'],
      ['age', 'w-1/12'],
    ],
    name,
    description,
    minWidth,
  })

  testK8sTableWithCustomColumns(Component, { createStore, name, description, minWidth })
})
