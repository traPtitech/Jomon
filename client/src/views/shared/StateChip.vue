<template>
  <div :class="$style.container">
    <v-chip
      :size="size || 'default'"
      :color="chip_color"
      :class="text_color"
      variant="flat"
      label
    >
      {{ state_ja }}
    </v-chip>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";

const props = withDefaults(
  defineProps<{
    state?: string;
    size?: string;
  }>(),
  {
    state: "",
    size: ""
  }
);

const state_ja = computed(() => {
  switch (props.state) {
    case "submitted":
      return "承認待ち";
    case "fix_required":
      return "要修正";
    case "accepted":
      return "承認済み";
    case "fully_repaid":
      return "返済完了";
    case "rejected":
      return "却下";
    default:
      return "ERROR";
  }
});

const chip_color = computed(() => {
  switch (props.state) {
    case "submitted":
      return "warning";
    case "rejected":
      return "grey";
    case "fix_required":
      return "error";
    case "accepted":
      return "info";
    case "fully_repaid":
      return "success";
    default:
      return "white";
  }
});

const text_color = computed(() => {
  switch (props.state) {
    case "submitted":
      return "text-black";
    case "rejected":
    case "fix_required":
    case "accepted":
    case "fully_repaid":
    default:
      return "text-white";
  }
});
</script>

<style lang="scss" module>
.container {
  margin: 0 4px;
}
</style>
