<template>
  <v-avatar :size="size" color="secondary">
    <v-img :src="iconSrc" width="100%" alt="" />
  </v-avatar>
</template>

<script setup lang="ts">
import { computed } from "vue";

const props = withDefaults(
  defineProps<{
    user?: string;
    size?: number;
  }>(),
  {
    user: "",
    size: 40
  }
);

const iconSrc = computed(() => {
  return props.user !== ""
    ? "https://q.trap.jp/api/v3/public/icon/" + props.user
    : "https://q.trap.jp/api/v3/public/icon/traP";
});
</script>
