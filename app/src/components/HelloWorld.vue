<template>
  <img alt="Vue logo" :src="require('../assets/logo.png')">
  <div class="hello">
    <p>Today's server message</p>
    <p>{{message}}</p>
    <a href="/logs">View Logs</a>
  </div>
  <v-btn depressed elevation="2" rounded>aa</v-btn>
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
    axios.get("http://localhost:8082/").then(message => {
      this.message = message.data;
    });
  }
}
</script>
