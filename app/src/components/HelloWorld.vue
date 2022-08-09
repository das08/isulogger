<template>
  <div class="hello">
    <p>Today's server message</p>
    <p>{{message}}</p>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'HelloWorld',
  data() {
    return {
      message: 'Loading...'
    }
  },
  mounted() {
    axios.get("/").then(message => {
      this.message = message.data;
    });
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
</style>
