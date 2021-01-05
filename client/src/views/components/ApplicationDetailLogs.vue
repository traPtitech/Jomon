<template>
  <div :class="$style.container">
    <h1>申請ログ</h1>
    <v-timeline dense clipped>
      <div v-for="(log, index) in this.logs" :key="index">
        <comment-log v-if="log.log_type === `comment`" :log="log" />
        <status-log v-else-if="log.log_type === `state`" :log="log" />
        <change-log v-else-if="log.log_type === `application`" :log="log" />
        <refund-log
          v-else-if="
            log.log_type === `repayment` &&
            !(log.content.repaid_at === `` || log.content.repaid_at === null)
          "
          :log="log"
        />
      </div>
    </v-timeline>
    <new-comment />
  </div>
</template>

<script>
import CommentLog from "./CommentLog";
import StatusLog from "./StatusLog";
import ChangeLog from "./ChangeLog";
import RefundLog from "./RefundLog";
import NewComment from "./TimelineNewComment";
import { mapGetters } from "vuex";

export default {
  props: {},
  computed: {
    ...mapGetters(["logs"])
  },
  components: {
    CommentLog,
    StatusLog,
    ChangeLog,
    RefundLog,
    NewComment
  }
};
</script>

<style lang="scss" module>
.container {
}
</style>
