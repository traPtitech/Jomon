<template>
  <div :class="$style.container">
    <h1>申請ログ</h1>
    <v-timeline dense clipped>
      <template v-for="(log, index) in logs" :key="index">
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
      </template>
    </v-timeline>
    <new-comment />
  </div>
</template>

<script>
import { mapGetters } from "vuex";
import ChangeLog from "./ChangeLog";
import CommentLog from "./CommentLog";
import RefundLog from "./RefundLog";
import StatusLog from "./StatusLog";
import NewComment from "./TimelineNewComment";

export default {
  components: {
    CommentLog,
    StatusLog,
    ChangeLog,
    RefundLog,
    NewComment
  },
  props: {},
  computed: {
    ...mapGetters(["logs"])
  }
};
</script>

<style lang="scss" module>
.container {
}
</style>
