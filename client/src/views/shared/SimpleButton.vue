<template>
  <button :class="$style.container" :data-variant="variant" @click="onClick">
    <div>{{ label }}</div>
  </button>
</template>

<script>
export default {
  props: {
    label: {
      type: String,
      required: true
    },
    variant: {
      type: String,
      default: "primary"
    }
  },
  methods: {
    onClick(e) {
      this.$emit("click", e);
    }
  }
};
</script>

<style lang="scss" module>
.container {
  display: flex;
  padding: 6px 16px;
  margin: 4px;
  border-radius: 4px;
  justify-content: center;
  align-items: center;
  cursor: pointer;
  transition: all 0.2s;
  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
  &[data-variant="primary"] {
    color: $color-text-primary;
    background-color: $color-primary;
    &:not(:disabled) {
      filter: drop-shadow(1px 2px 4px rgba($color-primary, 0.5));
    }
    &:hover:not(:disabled) {
      background-color: rgba($color-primary, 0.8);
    }
  }
  &[data-variant="secondary"] {
    color: $color-text-primary;
    background-color: $color-secondary;
    &:not(:disabled) {
      filter: drop-shadow(1px 2px 4px rgba($color-secondary, 0.5));
    }
    &:hover:not(:disabled) {
      background-color: rgba($color-secondary, 0.8);
    }
  }
  &[data-variant="warning"] {
    color: $color-text-primary;
    background-color: $color-warning;
    &:not(:disabled) {
      filter: drop-shadow(1px 2px 4px rgba($color-warning, 0.5));
    }
    &:hover:not(:disabled) {
      background-color: rgba($color-warning, 0.8);
    }
  }
  &[data-variant="error"] {
    color: $color-text-primary;
    background-color: $color-error;
    &:not(:disabled) {
      filter: drop-shadow(1px 2px 4px rgba($color-error, 0.5));
    }
    &:hover:not(:disabled) {
      background-color: rgba($color-error, 0.8);
    }
  }
}
.icon {
  margin-right: 4px;
}
</style>
