<template>
  <v-avatar :size="size">
    <v-img
      :src="`https://q.trap.jp/api/1.0/public/icon/` + this.user"
      width="100%"
      alt=""
    />
  </v-avatar>
</template>

<script>
export default {
  name: "Icon",
  computed: {},
  props: {
    user: String,
    size: {
      type: Number,
      default: 40
    }
  }
};
</script>
