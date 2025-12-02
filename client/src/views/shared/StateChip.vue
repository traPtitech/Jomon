<template>
  <div :class="$style.container">
    <v-chip
      :size="size || 'default'"
      :color="chip_color"
      :class="text_color"
      label
    >
      {{ state_ja }}
    </v-chip>
  </div>
</template>

<script>
export default {
  props: {
    state: {
      type: String,
      default: ""
    },
    size: {
      type: String,
      default: ""
    }
  },
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
    chip_color() {
      switch (this.state) {
        case "submitted":
          return "yellow";
        case "rejected":
          return "grey";
        case "fix_required":
          return "red";
        case "accepted":
          return "blue";
        case "fully_repaid":
          return "green";
        default:
          return "white";
      }
    },
    text_color() {
      switch (this.state) {
        case "submitted":
        case "rejected":
        case "fully_repaid":
        default:
          return "text-black";
        case "fix_required":
        case "accepted":
          return "text-white";
      }
    }
  }
};
</script>

<style lang="scss" module>
.container {
  margin: 0 4px;
}
</style>
