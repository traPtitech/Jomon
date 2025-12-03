<!-- ここはログの最小単位の表示を制御-->
<template>
  <div class="$style.conteiner">
    {{ formatted }}
  </div>
</template>

<script lang="ts">
export default {
  props: {
    date: {
      type: String,
      required: true
    },
    simple: {
      type: Boolean,
      default: false
    }
  },
  computed: {
    formatted: function () {
      const now = new Date();
      const d = new Date(this.date);
      if (this.simple) {
        return `${d.getFullYear()}/${d.getMonth() + 1}/${d.getDate()}`;
      }
      const diff = (now.getTime() - d.getTime()) / 1000;
      if (diff < 60) {
        //1分以内
        return Math.round(diff) + "秒前";
      } else if (diff < 60 * 60) {
        //一時間以内
        return Math.round(diff / 60) + "分前";
      } else if (diff < 60 * 60 * 24) {
        //一日以内
        return Math.round(diff / 60 / 60) + "時間前";
      } else if (diff < 60 * 60 * 24 * 28) {
        //一か月以内
        const month = d.getMonth() + 1;
        const day = d.getDate();
        return month + "/" + day;
      } else {
        const year = d.getFullYear();
        const month = d.getMonth() + 1;
        const day = d.getDate();
        return year + "/" + month + "/" + day;
      }
    }
  }
};
</script>
<style lang="scss" module>
.container {
}
</style>
