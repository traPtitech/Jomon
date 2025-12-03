<template>
  <v-app>
    <nav-bar />
    <v-main>
      <router-view />
    </v-main>
  </v-app>
</template>

<script setup lang="ts">
import { useMeStore } from "@/stores/me";
import NavBar from "./views/core/NavBar.vue";

const store = useMeStore();
store.fetchMe();
</script>
