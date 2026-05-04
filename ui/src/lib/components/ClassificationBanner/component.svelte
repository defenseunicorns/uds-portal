<script lang="ts">
  import { classColorMap, Classification, getClassification } from './helpers'

  export let text = ''
  export let enabled = false
  export let element: 'header' | 'footer' = 'header'

  let classification: Classification

  $: classification = enabled ? getClassification(text) : Classification.Unknown
</script>

<svelte:element
  this={element}
  data-testid="classification-{element}"
  tabindex="-1"
  class="z-50 flex h-6 w-full justify-between"
  class:hidden={!enabled}
  style="background-color: {classColorMap[classification][0]};"
>
  <div class="mx-auto flex items-center">
    <p class="flex items-center text-base font-semibold" style="color: {classColorMap[classification][1]}">
      {classification}
    </p>
  </div>
</svelte:element>
