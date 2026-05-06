// Copyright 2025 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

import '@testing-library/jest-dom'

if (!Element.prototype.animate) {
  Element.prototype.animate = () =>
    ({
      cancel: () => {},
      finish: () => {},
      play: () => {},
      pause: () => {},
      reverse: () => {},
      addEventListener: () => {},
      removeEventListener: () => {},
      dispatchEvent: () => true,
      finished: Promise.resolve(),
      onfinish: null,
      oncancel: null,
      playState: 'finished',
      currentTime: 0,
      effect: null,
      id: '',
      pending: false,
      playbackRate: 1,
      ready: Promise.resolve(),
      replaceState: 'active',
      startTime: 0,
      timeline: null,
      commitStyles: () => {},
      persist: () => {},
      updatePlaybackRate: () => {},
    }) as unknown as Animation
}
