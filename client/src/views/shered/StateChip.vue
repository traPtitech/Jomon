<template>
  <span class="text-center">
    <v-chip
      :x-small="size === `x-small`"
      :small="size === `small`"
      :large="size === `large`"
      :x-large="size === `x-large`"
      :class="chip"
      label
    >
      {{ state_text }}
    </v-chip>
  </span>
</template>

<script>
export default {
  name: "StateChip",
  computed: {
    state_text() {
      switch (this.state) {
        case "submitted":
          return "提出済";
        case "rejected":
          return "却下";
        case "fix_required":
          return "要修正";
        case "accepted":
          return "許可済";
        case "fully_repaid":
          return "払戻完";
        default:
          return "error";
      }
    },
    chip() {
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
      type: String,
      default: "fail"
    },
    size: {
      type: String,
      default: ""
    }
  }
};
</script>
