<template>
  <v-container>
  <form>
    <v-text-field
        v-model="name"
        :error-messages="nameErrors"
        :counter="10"
        label="Name"
        required
    ></v-text-field>
    <v-text-field
        v-model="email"
        :error-messages="emailErrors"
        label="E-mail"
        required
    ></v-text-field>
    <v-select
        v-model="select"
        :items="items"
        :error-messages="selectErrors"
        label="Item"
        required
    ></v-select>

    <v-btn
        class="mr-4"
        @click="submit"
    >
      submit
    </v-btn>
    <v-btn @click="clear">
      clear
    </v-btn>
  </form>
  </v-container>
</template>

<script>
export default {
  data: () => ({
    name: '',
    email: '',
    select: null,
    items: [
      'Item 1',
      'Item 2',
      'Item 3',
      'Item 4',
    ],
    checkbox: false,
  }),
  computed: {
    checkboxErrors () {
      return []
    },
    selectErrors () {
      return []
    },
    nameErrors () {
      return []
    },
    emailErrors () {
      return []
    },
  },
  methods: {
    submit () {

    },
    clear () {
      this.name = ''
      this.email = ''
      this.select = null
      this.checkbox = false
    },
  },
}
</script>