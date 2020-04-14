<template>
  <v-app>
    <nav-bar />
    <core-login-dialog />
    <core-view />
  </v-app>
</template>

<script>
import NavBar from "./views/core/NavBar";
import CoreView from "./views/core/CoreView";
import CoreLoginDialog from "./views/core/LoginDialog";
export default {
  name: "App",
  components: {
    NavBar,
    CoreView,
    CoreLoginDialog
  },
  created() {
    this.$store.dispatch("getMe");
  },
  data: () => ({
    //
  })
};
</script>
