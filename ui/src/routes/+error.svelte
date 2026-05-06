<!-- Copyright 2025 Defense Unicorns -->
<!-- SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial -->

<script lang="ts">
  import { resolve } from '$app/paths'
</script>

<section class="bg-white dark:bg-gray-900">
  <div class="mx-auto max-w-screen-xl px-4 py-8 lg:px-6 lg:py-16">
    <div class="mx-auto max-w-screen-sm text-center">
      <h1 class="mb-4 text-7xl font-extrabold tracking-tight text-gray-900 lg:text-9xl dark:text-gray-100">404</h1>
      <p class="mb-4 text-3xl font-bold tracking-tight text-gray-900 md:text-4xl dark:text-gray-100">
        Something's missing.
      </p>
      <p class="mb-4 text-lg font-light text-gray-500 dark:text-gray-100">
        Sorry, we can't find that page. You'll find lots to explore on the home page.
      </p>
      <a
        href={resolve('/')}
        class="bg-primary-600 hover:bg-primary-800 focus:ring-primary-300 dark:focus:ring-primary-900 my-4 inline-flex rounded-lg px-5 py-2.5 text-center text-sm font-medium text-gray-100 focus:ring-4 focus:outline-none"
      >
        Back to Homepage
      </a>
    </div>
  </div>
</section>
