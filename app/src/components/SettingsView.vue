<template>
  <div class="settings-view">
    <v-text-field
      label="Secret key"
      v-model="secretKey"
    ></v-text-field>
    <v-btn elevation="2" @click="onSave">Save</v-btn>
  </div>
</template>

<script>
import {LS_SECRET_KEY} from '../ls.js';

export default {
  name: 'SettingsView',
  data: () => {
    return {
      secretKey: ""
    }
  },
  methods: {
    onSave() {
      localStorage.setItem(LS_SECRET_KEY, this.secretKey);
    }
  },
  mounted() {
    const secretKey = localStorage.getItem(LS_SECRET_KEY);
    if (secretKey !== null && typeof(secretKey) === "string") {
      this.secretKey = secretKey;
    }
  }
}
</script>

<style scoped>
.settings-view {
  margin: 10px;
}
</style>
