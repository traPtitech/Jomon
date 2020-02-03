<template>
  <v-app-bar app color="primary" dark>
    <div class="d-flex align-center">
      <v-img
        alt="Vuetify Logo"
        class="shrink-mr-2"
        contain
        src="https://cdn.vuetifyjs.com/images/logos/vuetify-logo-dark.png"
        transition="scale-transition"
        width="40"
      />

      <router-link to="/">
        <v-toolbar-title class="white--text">Jomon</v-toolbar-title>
      </router-link>
    </div>
    <v-spacer></v-spacer>

    <v-btn to="/admin" text>
      管理ページ
    </v-btn>
    <v-btn to="/applications/new" text>
      新規作成
    </v-btn>
  </v-app-bar>
</template>
