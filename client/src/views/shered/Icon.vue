<template>
  <v-avatar :size="size">
    <v-img :src="cPicture" width="100%" alt="" @error="errorHandler" />
  </v-avatar>
</template>

<script>
export default {
  name: "Icon",
  data: () => {
    return {
      failed_image: false
    };
  },
  computed: {
    cPicture: function() {
      return this.failed_image
        ? "https://q.trap.jp/api/v3/public/icon/traP" //userがtraPIDではない時traPIconを表記
        : "https://q.trap.jp/api/v3/public/icon/" + this.user;
    }
  },
  props: {
    user: String,
    size: {
      type: Number,
      default: 40
    }
  },
  methods: {
    errorHandler() {
      this.failed_image = true;
    }
  }
};
</script>
