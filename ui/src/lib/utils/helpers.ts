// Copyright 2024 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

import { differenceInDays, differenceInHours, differenceInMinutes, differenceInSeconds } from 'date-fns'

export const stringToSnakeCase = (name: string) => name.split(' ').join('-').toLocaleLowerCase()

export function formatDetailedAge(timestamp: Date) {
  const now = new Date()
  const seconds = differenceInSeconds(now, timestamp)

  if (seconds < 60) {
    return `${seconds}s`
  }

  const minutes = differenceInMinutes(now, timestamp)
  if (minutes < 60) {
    return `${minutes}m`
  }

  const hours = differenceInHours(now, timestamp)
  if (hours < 24) {
    const remainingMinutes = minutes % 60
    return remainingMinutes > 0 ? `${hours}h${remainingMinutes}m` : `${hours}h`
  }

  const days = differenceInDays(now, timestamp)
  return `${days} days`
}
