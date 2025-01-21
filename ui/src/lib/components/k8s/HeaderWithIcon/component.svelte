<script lang="ts">
  import { Information } from 'carbon-icons-svelte'

  export let title: string
  export let size: 'sm' | 'md' | 'lg' | 'xl' = 'lg'
  export let height: string = '6'
  export let results: string = ''
  export let description: string
</script>

<div class="flex h-{height} items-center space-x-1">
  <h2 class="text-{size} dark:text-gray-100 font-semibold justify-items-start">{title}</h2>

  {#if results.length > 0}
    <span class="pl-1 dark:text-gray-500 text-blue-500">({results})</span>
  {/if}

  <div class="relative group">
    <Information class="ml-1 w-4 h-4 dark:text-gray-100 text-blue-500" />

    <div class="tooltip tooltip-right min-w-72">
      <div class="whitespace-normal">{description}</div>
    </div>
  </div>
</div>
