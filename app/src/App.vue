<template>
  <v-app>
    <v-app-bar
      app
      color="primary"
      dark
    >
      <div class="d-flex align-center">
        <router-link to="/">
          <v-img
            alt="Vuetify Logo"
            class="shrink mr-2"
            contain
            src="./assets/logo.jpg"
            transition="scale-transition"
            width="40"
          />
        </router-link>

        <router-link to="/">
          <v-img
            alt="Vuetify Name"
            class="shrink mt-1 hidden-sm-and-down"
            contain
            min-width="100"
            src="./assets/icon.png"
            width="160"
          />
        </router-link>
      </div>

      <v-spacer></v-spacer>
      <v-app-bar-nav-icon @click="drawer = true"></v-app-bar-nav-icon>
    </v-app-bar>
    <v-navigation-drawer
        v-model="drawer"
        fixed
        right
        temporary
    >
      <v-list-item
          v-for="item in items"
          :key="item.title"
          :to="item.link"
      >
        <v-list-item-icon>
          <v-icon>{{ item.icon }}</v-icon>
        </v-list-item-icon>

        <v-list-item-content>
          <v-list-item-title>{{ item.title }}</v-list-item-title>
        </v-list-item-content>
      </v-list-item>
    </v-navigation-drawer>

    <v-main>
      <router-view></router-view>
    </v-main>
  </v-app>
</template>

<script>
export default {
  name: 'App',

  components: {
  },

  data: () => ({
    drawer: false,
    items: [
      {title: 'ホーム', icon: 'mdi-home', link: {name: 'HelloWorld'}},
      {title: 'ログを見る', icon: 'mdi-text-box-multiple', link: {name: 'Logs'}},
      {title: 'コンテスト作成', icon: 'mdi-plus-circle', link: {name: 'Create Contest'}},
    ]
  }),
};
</script>
