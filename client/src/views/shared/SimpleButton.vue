<template>
  <v-btn :color="variant" :class="$style.button" @click="onClick">
    {{ label }}
  </v-btn>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    label: string;
    variant?: string;
  }>(),
  {
    variant: "primary"
  }
);

const emit = defineEmits<{
  (e: "click", event: Event): void;
}>();

const onClick = (e: Event) => {
  emit("click", e);
};
</script>

<style lang="scss" module>
.button {
  margin: 4px;
}
</style>
