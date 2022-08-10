<template>
  <alp-diff :cmp1="parseInt($route.params.cmp1)" :cmp2="parseInt($route.params.cmp2)"></alp-diff>
</template>

<script>
import AlpDiff from './AlpDiff.vue'

export default {
  name: "AlpDiffRouter",
  components: {AlpDiff},
}
</script>