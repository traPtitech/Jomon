<template>
  <v-chip
    :x-small="size === `x-small`"
    :small="size === `small`"
    :large="size === `large`"
    :x-large="size === `x-large`"
    :class="chip"
    label
    :width="width"
  >
    {{ this.state }}
  </v-chip>
</template>

<script>
export default {
  name: "StateChip",
  computed: {
    chip: function() {
      switch (this.state) {
        case "submitted":
          return "black white--text";
        case "rejected":
          return "yellow red--text";
        case "fix_required":
          return "orange white--text";
        case "accepted":
          return "red white--text";
        case "fully_repaid":
          return "gray black--text";
        default:
          return "white red--text";
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
    },
    width: {
      type: String,
      default: ""
    }
  }
};
</script>
