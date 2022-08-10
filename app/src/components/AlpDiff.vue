<template>
  <p>{{alp1}}</p>
</template>

<script>
import axios from 'axios'

export default {
  name: "AlpDiff",
  props: {
    "cmp1": Number,
    "cmp2": Number,
  },
  data() {
    return {
      "alp1": null,
      "alp2": null,
    }
  },
  mounted() {
    axios.get("/parsed/alp/" + this.cmp1)
    .then((response) => {
      this.alp1 = response.data;
    })
  }
}
</script>