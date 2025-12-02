<template>
  <v-app>
    <nav-bar />
    <v-main>
      <router-view />
    </v-main>
  </v-app>
</template>

<script>
import NavBar from "./views/core/NavBar.vue";

export default {
  name: "App",
  components: {
    NavBar
  },
  created() {
    this.$store.dispatch("getMe");
  }
};
</script>
