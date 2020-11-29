<template>
  <div :class="$style.container">
    <v-chip
      :x-small="size === `x-small`"
      :small="size === `small`"
      :large="size === `large`"
      :x-large="size === `x-large`"
      :class="chip_style"
      label
    >
      {{ state_ja }}
    </v-chip>
  </div>
</template>

<script>
export default {
  computed: {
    state_ja() {
      switch (this.state) {
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
    },
    chip_style() {
      switch (this.state) {
        case "submitted":
          return "yellow black--text";
        case "rejected":
          return "gray black--text";
        case "fix_required":
          return "red white--text";
        case "accepted":
          return "blue white--text";
        case "fully_repaid":
          return "green black--text";
        default:
          return "white black--text";
      }
    }
  },

  props: {
    state: {
      type: String
    },
    size: {
      type: String,
      default: ""
    }
  }
};
</script>

<style lang="scss" module>
.container {
  margin: 0 4px;
}
</style>
