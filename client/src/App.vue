<template>
  <v-app>
    <nav-bar />

    <v-main>
      <router-view />
    </v-main>
  </v-app>
</template>

<script>
import NavBar from "./views/core/NavBar";

export default {
  name: "App",
  components: {
    NavBar
  },
  created() {
    this.$store.dispatch("getMe");
  },
  data: () => ({
    //
  })
};
</script>
