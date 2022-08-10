<template>
  <v-container>
    <v-alert
        v-model="success"
        dense
        text
        type="success"
    >
      コンテストを作成しました!
    </v-alert>
    <v-alert
        v-model="failure"
        dense
        text
        type="error"
    >
      コンテストを作成できませんでした
    </v-alert>
  <form>
    <v-text-field
        v-model="contest_name"
        :error-messages="nameErrors"
        :counter="20"
        label="コンテスト名"
        required
    ></v-text-field>

    <v-btn
        class="mr-4"
        @click="submit"
    >
      作成する
    </v-btn>
    <v-btn @click="clear">
      リセット
    </v-btn>
  </form>
  </v-container>
</template>

<script>
import axios from "axios";

export default {
  data: () => ({
    contest_name: '',
    success: false,
    failure: false,
  }),
  computed: {
    nameErrors () {
      if (this.contest_name.length < 3) {
        return ['コンテスト名は3文字以上です(適当)']
      }
      return []
    },
  },
  methods: {
    submit () {
      console.log("submit", this.contest_name);
      axios.post("http://localhost:8082/new_contest", {
        contest_name: this.contest_name,
      })
          .then((response) => {
            console.log(response.data);
            this.contest_name = '';
            this.success = true;
            this.failure = false;
          })
          .catch((error) => {
            console.log(error);
            this.success = false;
            this.failure = true;
          });
    },
    clear () {
      this.contest_name = ''
    },
  },
}
</script>