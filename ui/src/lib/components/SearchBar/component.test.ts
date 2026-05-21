// Copyright 2026 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

import { cleanup, fireEvent, render, screen } from '@testing-library/svelte'
import { afterEach, describe, expect, it } from 'vitest'

import SearchBar from './component.svelte'

describe('SearchBar', () => {
  afterEach(cleanup)

  it('renders search input with placeholder', () => {
    render(SearchBar, { props: { value: '' } })
    expect(screen.getByPlaceholderText('Search')).toBeInTheDocument()
  })

  it('clear button hidden when value is empty', () => {
    render(SearchBar, { props: { value: '' } })
    const clearBtn = screen.getByRole('button', { name: /clear/i })
    expect(clearBtn).toHaveClass('hidden')
  })

  it('clear button visible when value is non-empty', () => {
    render(SearchBar, { props: { value: 'argo' } })
    const clearBtn = screen.getByRole('button', { name: /clear/i })
    expect(clearBtn).not.toHaveClass('hidden')
  })

  it('clicking clear button resets input value', async () => {
    render(SearchBar, { props: { value: 'argo' } })
    const clearBtn = screen.getByRole('button', { name: /clear/i })
    await fireEvent.click(clearBtn)
    const input = screen.getByPlaceholderText('Search') as HTMLInputElement
    expect(input.value).toBe('')
    expect(clearBtn).toHaveClass('hidden')
  })

  it('clear button appears after typing in input', async () => {
    render(SearchBar, { props: { value: '' } })
    const input = screen.getByPlaceholderText('Search') as HTMLInputElement
    const clearBtn = screen.getByRole('button', { name: /clear/i })
    expect(clearBtn).toHaveClass('hidden')
    await fireEvent.input(input, { target: { value: 'Grafana' } })
    expect(clearBtn).not.toHaveClass('hidden')
  })
})
