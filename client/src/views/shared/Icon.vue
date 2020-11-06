<template>
  <v-avatar :size="size">
    <v-img :src="iconSrc" width="100%" alt="" />
  </v-avatar>
</template>

<script>
import store from "@store";

export default {
  name: "Icon",
  data: () => {
    return {
      failed_image: false
    };
  },
  computed: {
    iconSrc: function () {
      return this.user !== ""
        ? "https://q.trap.jp/api/v3/public/icon/" + this.user
        : "https://q.trap.jp/api/v3/public/icon/traP";
    }
  },
  props: {
    user: String,
    size: {
      type: Number,
      default: 40
    }
  }
};
</script>
