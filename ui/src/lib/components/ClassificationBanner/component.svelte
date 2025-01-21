<script lang="ts">
  import { onMount } from 'svelte'

  import { classColorMap, Classification, getClassification } from './helpers'

  export let text = ''
  export let enabled = false
  export let element: 'header' | 'footer' = 'header'

  let classification: Classification = Classification.Unknown

  onMount(() => {
    if (enabled) {
      classification = getClassification(text)
    }
  })
</script>

<svelte:element
  this={element}
  data-testid="classification-{element}"
  tabindex="-1"
  class="flex justify-between w-full h-6 z-50"
  class:hidden={!enabled}
  style="background-color: {classColorMap[classification][0]};"
>
  <div class="flex items-center mx-auto">
    <p class="flex items-center text-base font-semibold" style="color: {classColorMap[classification][1]}">
      {classification}
    </p>
  </div>
</svelte:element>
