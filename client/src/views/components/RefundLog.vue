<!-- ここはログの最小単位の表示を制御-->
<template>
  <v-timeline-item color="grey" :class="$style.text">
    <icon :user="log.content.repaid_by_user.trap_id" :size="25" />
    <strong>
      {{ log.content.repaid_by_user.trap_id }}
    </strong>
    が
    <strong>
      {{ log.content.repaid_to_user.trap_id }}
    </strong>
    に
    <strong>
      {{ datePrint(log.content.repaid_at) }}
    </strong>
    に払い戻しました。
    {{ dayPrint(log.content.created_at) }}
  </v-timeline-item>
</template>

<script setup lang="ts">
import { RefundLog } from "@/types/log";
import Icon from "@/views/shared/Icon.vue";

defineProps<{
  log: RefundLog;
}>();

const dayPrint = (time: string) => {
  const now = new Date();
  const d = new Date(time);
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
};

const datePrint = (date: string) => {
  const d = new Date(date);
  const year = d.getFullYear();
  const month = d.getMonth() + 1;
  const day = d.getDate();
  const res = year + "/" + month + "/" + day;
  return res;
};
</script>

<style lang="scss" module>
.text {
  color: $color-grey;
}
</style>
