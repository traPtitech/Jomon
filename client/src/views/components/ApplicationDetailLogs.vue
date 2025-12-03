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

<script setup lang="ts">
import { useApplicationDetailStore } from "@/stores/applicationDetail";
import { storeToRefs } from "pinia";
import ChangeLog from "./ChangeLog.vue";
import CommentLog from "./CommentLog.vue";
import RefundLog from "./RefundLog.vue";
import StatusLog from "./StatusLog.vue";
import NewComment from "./TimelineNewComment.vue";

const applicationDetailStore = useApplicationDetailStore();
const { logs } = storeToRefs(applicationDetailStore);
</script>

<style lang="scss" module>
.container {
}
</style>
