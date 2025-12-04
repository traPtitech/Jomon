<template>
  <div :class="$style.container">
    <h1>申請ログ</h1>
    <v-timeline density="compact" side="end">
      <template v-for="(log, index) in logs" :key="index">
        <v-timeline-item
          v-if="log.log_type === `comment`"
          dot-color="purple"
          size="small"
        >
          <comment-log :log="log" />
        </v-timeline-item>
        <v-timeline-item
          v-else-if="log.log_type === `state`"
          dot-color="red"
          size="small"
        >
          <status-log :log="log" />
        </v-timeline-item>
        <v-timeline-item
          v-else-if="log.log_type === `application`"
          dot-color="purple"
          size="small"
        >
          <change-log :log="log" />
        </v-timeline-item>
        <v-timeline-item
          v-else-if="
            log.log_type === `repayment` &&
            !(log.content.repaid_at === `` || log.content.repaid_at === null)
          "
          dot-color="grey"
          size="small"
        >
          <refund-log :log="log" />
        </v-timeline-item>
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
  :global {
    .v-timeline--density-compact .v-timeline-item {
      min-height: 0 !important;
    }
    .v-timeline-item__body {
      padding-block-start: 0 !important;
      padding-block-end: 0 !important;
      margin-bottom: 4px !important;
    }
  }
}
</style>
