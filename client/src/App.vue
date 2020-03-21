<template>
  <v-app>
    <nav-bar />
    <core-view />
  </v-app>
</template>

<script>
import NavBar from "./views/core/NavBar";
import CoreView from "./views/core/CoreView";
export default {
  name: "App",
  components: {
    NavBar,
    CoreView
  },
  created() {
    this.$store.dispatch("getMe");
  },
  data: () => ({
    //
  })
};
</script>
