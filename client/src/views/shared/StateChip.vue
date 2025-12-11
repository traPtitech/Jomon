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
import {
  getStateColor,
  getStateLabel,
  getStateTextColor
} from "@/use/stateColor";
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

const state_ja = computed(() => getStateLabel(props.state));
const chip_color = computed(() => getStateColor(props.state));
const text_color = computed(() => getStateTextColor(props.state));
</script>

<style lang="scss" module>
.container {
  margin: 0 4px;
}
</style>
