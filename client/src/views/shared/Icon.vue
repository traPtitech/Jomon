<template>
  <v-avatar :size="size" color="secondary">
    <v-img :src="iconSrc" width="100%" alt="" />
  </v-avatar>
</template>

<script>
export default {
  name: "UserIcon",
  props: {
    user: {
      type: String,
      default: ""
    },
    size: {
      type: Number,
      default: 40
    }
  },
  computed: {
    iconSrc: function () {
      return this.user !== ""
        ? "https://q.trap.jp/api/v3/public/icon/" + this.user
        : "https://q.trap.jp/api/v3/public/icon/traP";
    }
  }
};
</script>
