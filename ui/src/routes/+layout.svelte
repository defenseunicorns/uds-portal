<!-- Copyright 2025 Defense Unicorns -->
<!-- SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial -->

<script lang="ts">
  import 'flowbite'

  import { afterNavigate } from '$app/navigation'
  import { authenticated } from '$features/auth/store'
  import { Navbar } from '$features/navigation'
  import { initFlowbite } from 'flowbite'

  import '../app.postcss'

  import { ClassBanner, Footer } from '$components'

  import { _bannerCfg } from './+layout'

  export let data

  afterNavigate(initFlowbite)

  // initFlowbite loads the js necessary to target components which use flowbite js
  // i.e. data-dropdown-toggle
  initFlowbite()

  $: if (!$authenticated) {
    if (window.location.pathname !== '/auth') {
      window.location.assign('/auth')
    }
  }

  $: bannerOffset = ($_bannerCfg.enabled ? 24 : 0) + ($_bannerCfg.footer ? 24 : 0)
  $: mainMinHeight = `calc(100dvh - ${32 + bannerOffset}px)`
</script>

<div class="flex h-screen flex-col">
  <ClassBanner enabled={$_bannerCfg.enabled} text={$_bannerCfg.text} element="header" />
  <div class="flex-grow overflow-auto">
    <Navbar userData={data.userData} />

    <main
      class="flex bg-[linear-gradient(180deg,_#030712_66.35%,_#213E68_100%)] pt-16 transition-all duration-300 ease-in-out"
      style="min-height: {mainMinHeight};"
    >
      <div class="mx-auto w-full max-w-6xl p-4 pt-16">
        <slot />
      </div>
    </main>
    <Footer />
  </div>
  <ClassBanner enabled={$_bannerCfg.footer} text={$_bannerCfg.text} element="footer" />
</div>
