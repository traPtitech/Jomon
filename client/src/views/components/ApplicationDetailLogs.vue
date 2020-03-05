<!-- 親コンポーネントから受け取ったjsonを基に申請書logを描画 -->
<!-- ここへのlogのjsonは日付の古い方を前にして突っ込んでください。--->
<template>
  <div class="app-detail-log">
    <template>
      <v-container style="max-width: 600px;">
        <h1>申請書ログ</h1>
        <v-timeline dense clipped>
          <Log v-for="log in this.logs" :key="log.log_type" :log="log" />
          <v-divider class="m pb-1" />
        </v-timeline>
        <NewComment />
      </v-container>
    </template>
  </div>
</template>

<script>
import Log from "./Log";
import NewComment from "./TimelineNewComment";
import { mapGetters } from "vuex";
export default {
  props: {},
  computed: {
    ...mapGetters(["logs"])
  },
  components: {
    Log,
    NewComment
  }
};
</script>
