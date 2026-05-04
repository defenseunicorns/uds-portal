// Copyright 2026 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

import { render, screen } from '@testing-library/svelte'
import { describe, expect, it } from 'vitest'

import ClassificationBanner from './component.svelte'
import { Classification as c, classColorMap } from './helpers'

const classifications = [
  { text: c.Unclassified, color: classColorMap[c.Unclassified] },
  { text: c.CUI, color: classColorMap[c.CUI] },
  { text: c.Confidential, color: classColorMap[c.Confidential] },
  { text: c.Secret, color: classColorMap[c.Secret] },
  { text: c.TopSecret, color: classColorMap[c.TopSecret] },
  { text: c.TopSecretSci, color: classColorMap[c.TopSecretSci] },
  { text: c.Unknown, color: classColorMap[c.Unknown] },
]

describe('ClassificationBanner', () => {
  classifications.forEach((c) => {
    it(`renders as header with ${c.text} classification with ${c.color} color`, () => {
      render(ClassificationBanner, { props: { enabled: true, text: c.text, element: 'header' } })
      const header = screen.getByTestId('classification-header')
      expect(header).toHaveStyle(`background-color: ${c.color[0]}`)
      const text = screen.getByText(c.text)
      expect(text).toHaveStyle(`color: ${c.color[1]}`)
    })
  })

  it('renders as footer', () => {
    render(ClassificationBanner, { props: { enabled: true, text: c.Unclassified, element: 'footer' } })
    const footer = screen.queryByTestId('classification-footer')
    expect(footer).toHaveStyle(`background-color: ${classColorMap[c.Unclassified][0]}`)

    const text = screen.getByText(c.Unclassified)
    expect(text).toHaveStyle(`color: ${classColorMap[c.Unclassified][1]}`)
  })

  it('renders unknown if classification not derived from text', () => {
    render(ClassificationBanner, { props: { enabled: true, text: 'unclass' } })
    const header = screen.getByTestId('classification-header')
    expect(header).toHaveStyle(`background-color: ${classColorMap[c.Unknown][0]}`)

    const text = screen.getByText(c.Unknown)
    expect(text).toHaveStyle(`color: ${classColorMap[c.Unknown][1]}`)
  })
})
