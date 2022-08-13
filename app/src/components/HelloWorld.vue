<template>
  <v-container>
    <v-row class="text-center">
      <v-col cols="12">
        <v-img
          :src="require('../assets/logo.svg')"
          class="my-3"
          contain
          height="200"
        />
      </v-col>

      <v-col class="mb-4">
        <h1 class="display-2 font-weight-bold mb-3">
          Welcome to ISULOGGER
        </h1>

        <div class="text-center d-flex align-center justify-space-around">
          <v-tooltip bottom color="primary">
            <template v-slot:activator="{ on, attrs }">
              <v-btn
                  color="primary"
                  dark
                  v-bind="attrs"
                  v-on="on"
                  to="/logs"
              >
                ログを見る
              </v-btn>
            </template>
            <span>アップロードされたログを確認できます</span>
          </v-tooltip>
        </div>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
  export default {
    name: 'HelloWorld',

    data: () => ({
      whatsNext: [
        {
          text: 'logs',
          href: '/logs',
        },
        {
          text: 'create contest',
          href: '/create_contest',
        },
      ],
    }),
  }
</script>
