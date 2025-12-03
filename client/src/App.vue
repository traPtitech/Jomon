<template>
  <v-app>
    <nav-bar />
    <v-main>
      <router-view />
    </v-main>
  </v-app>
</template>

<script lang="ts">
import { useMeStore } from "@/stores/me";
import { mapActions } from "pinia";
import NavBar from "./views/core/NavBar.vue";
export default {
  name: "App",
  components: {
    NavBar
  },
  created() {
    this.fetchMe();
  },
  methods: {
    ...mapActions(useMeStore, ["fetchMe"])
  }
};
</script>
