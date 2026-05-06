// Copyright 2026 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

export enum Classification {
  Unclassified = 'UNCLASSIFIED',
  CUI = 'CUI',
  Confidential = 'CONFIDENTIAL',
  Secret = 'SECRET',
  TopSecret = 'TOP SECRET',
  TopSecretSci = 'TOP SECRET//SCI',
  Unknown = 'UNKNOWN',
}

export const classColorMap: Record<Classification, string[]> = {
  UNCLASSIFIED: ['#007a33', '#ffffff'],
  CUI: ['#502b85', '#ffffff'],
  CONFIDENTIAL: ['#0033a0', '#ffffff'],
  SECRET: ['#c8102e', '#ffffff'],
  'TOP SECRET': ['#ff8c00', '#000000'],
  'TOP SECRET//SCI': ['#fce83a', '#000000'],
  UNKNOWN: ['#000000', '#ffffff'],
}

export function getClassification(classification: string) {
  if (Object.values(Classification).includes(classification.toUpperCase() as Classification)) {
    return classification.toUpperCase() as Classification
  }

  console.error(`Invalid classification: ${classification}; using ${Classification.Unknown} instead`)
  return Classification.Unknown
}
