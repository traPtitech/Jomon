<template>
  <v-app>
    <nav-bar />
    <v-main>
      <router-view />
    </v-main>
    <v-snackbar v-model="toastStore.isVisible" :color="toastStore.color">
      {{ toastStore.message }}
      <template #actions>
        <v-btn
          color="white"
          variant="text"
          @click="toastStore.isVisible = false"
        >
          閉じる
        </v-btn>
      </template>
    </v-snackbar>
  </v-app>
</template>

<script setup lang="ts">
import { useMeStore } from "@/stores/me";
import { useToastStore } from "@/stores/toast";
import NavBar from "./views/core/NavBar.vue";

const store = useMeStore();
const toastStore = useToastStore();
store.fetchMe();
</script>
