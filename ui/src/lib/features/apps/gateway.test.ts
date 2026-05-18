// Copyright 2026 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

import { describe, expect, it } from 'vitest'

import { isAdminGateway } from './gateway'

describe('isAdminGateway', () => {
  it.each([
    ['admin', true],
    ['Admin', true],
    ['Admin-East', true],
    ['custom-admin-gw', true],
    ['tenant', false],
    ['passthrough', false],
    ['shared-services', false],
    ['', false],
  ])('classifies %s as admin=%s', (gw, expected) => {
    expect(isAdminGateway(gw)).toBe(expected)
  })

  it('treats undefined as non-admin', () => {
    expect(isAdminGateway(undefined)).toBe(false)
  })
})
