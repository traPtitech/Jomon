<template>
  <v-avatar :size="size">
    <v-img :src="iconSrc" width="100%" alt="" />
  </v-avatar>
</template>

<script>
export default {
  name: "Icon",
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
