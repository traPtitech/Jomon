<!-- 親コンポーネントから受け取ったjsonを基に申請書logを描画 -->
<!-- ここへのlogのjsonは日付の古い方を前にして突っ込んでください。--->
<template>
  <div class="app-detail-log">
    <template>
      <v-container style="max-width: 600px;">
        <h1>ここからログ</h1>
        <v-timeline dense clipped>
          <!-- <Log :log="this.$store.state.logs[0]" /> -->
          <Log
            v-for="log in this.$store.state.logs"
            :key="log.log_type"
            :log="log"
          />
          <v-divider class="m pb-1" />
        </v-timeline>
        <NewComment />

        <h1>ここまでログ</h1>
      </v-container>
    </template>
  </div>
</template>

<script>
import Log from "./Log";
import NewComment from "./TimelineNewComment";
export default {
  //nameは.vueファイルはファイル名がnameになるので不要
  //props: ["log"]
  props: {
    //申請書詳細はpropsで管理しない
    // log: Object
  },
  components: {
    Log,
    NewComment
  }
};
</script>
