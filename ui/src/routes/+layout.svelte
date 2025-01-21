<!-- Copyright 2024 Defense Unicorns -->
<!-- SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial -->

<script lang="ts">
  import 'flowbite'

  import { onDestroy, onMount } from 'svelte'

  import { afterNavigate, goto } from '$app/navigation'
  import { authenticated } from '$features/auth/store'
  import { LoadingCluster, Navbar, Sidebar } from '$features/navigation'
  import { ToastPanel } from '$features/toast'
  import { initFlowbite } from 'flowbite'

  import '../app.postcss'

  import { ClassBanner, Footer } from '$components'
  import { loadingCluster } from '$features/navigation/navbar/clustermenu/store'
  import { ClusterCheck } from '$lib/utils/cluster-check/cluster-check'

  import { _bannerCfg } from './+layout'

  export let data

  let clusterCheck: ClusterCheck

  onMount(() => {
    // initFlowbite loads the js necessary to target components which use flowbite js
    // i.e. data-dropdown-toggle
    initFlowbite()
  })

  onDestroy(() => {
    if (clusterCheck) clusterCheck.close()
  })

  afterNavigate(initFlowbite)

  $: if ($authenticated) {
    clusterCheck = new ClusterCheck()
  } else {
    goto('/auth')
  }
</script>

<div class="flex h-screen flex-col">
  <ClassBanner enabled={$_bannerCfg.enabled} text={$_bannerCfg.text} element="header" />
  <div class="flex-grow overflow-auto">
    <Navbar userData={data.userData} />

    {#if $authenticated}
      <Sidebar />
    {/if}

    <main class="flex h-[95%] pt-16 transition-all duration-300 ease-in-out">
      <ToastPanel />
      {#if $loadingCluster.loading}
        <LoadingCluster cluster={$loadingCluster.cluster} error={$loadingCluster.err} />
      {:else}
        <div class="mx-auto w-full max-w-3xl p-4 pt-16">
          <slot />
        </div>
      {/if}
    </main>
    <Footer/>
  </div>
  <ClassBanner enabled={$_bannerCfg.footer} text={$_bannerCfg.text} element="footer" />
</div>
