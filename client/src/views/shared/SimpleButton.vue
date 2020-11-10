<template>
  <button :class="$style.container" :data-variant="variant" @click="onClick">
    <div>{{ label }}</div>
  </button>
</template>

<script>
export default {
  name: "SimpleButton",
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
    onClick() {
      this.$emit("click");
    }
  }
};
</script>

<style lang="scss" module scoped>
@import "../../styles/index.scss"; // TODO: なんかGlobal変数使えん
.container {
  display: flex;
  padding: 4px 12px;
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
  &:not(:disabled) {
    filter: drop-shadow(1px 2px 4px rgba($color-primary, 0.5));
  }
  &[data-variant="primary"] {
    color: $color-text-primary;
    background-color: $color-primary;
    &:hover:not(:disabled) {
      background-color: $color-primary-hover;
    }
  }
  &[data-variant="secondary"] {
    color: $color-text-secondary;
    background-color: $color-secondary;
    &:hover:not(:disabled) {
      background-color: $color-secondary-hover;
    }
  }
  &[data-variant="caution"] {
    color: $color-error;
    background-color: $color-background;
    &:hover:not(:disabled) {
      color: $color-background;
      background-color: $color-error;
    }
  }
}
.icon {
  margin-right: 4px;
}
</style>
