<template>
  <div>
    <div :class="$style.text">
      <icon :user="log.content.update_user.trap_id" :size="25" />
      {{ log.content.update_user.trap_id }}
      が申請の状態を
      <state-chip :state="log.content.to_state" />に変更しました。
      {{ dayPrint(log.content.created_at) }}
    </div>
    <v-card v-if="log.content.reason !== ''">
      <v-card-text> 理由: {{ log.content.reason }} </v-card-text>
    </v-card>
  </div>
</template>

<script setup lang="ts">
import { StatusLog } from "@/types/log";
import Icon from "@/views/shared/Icon.vue";
import StateChip from "@/views/shared/StateChip.vue";

defineProps<{
  log: StatusLog;
}>();

const dayPrint = (time: string | Date) => {
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
</script>

<style lang="scss" module>
.text {
  display: flex;
  flex-wrap: wrap;
  max-width: 100%;
  align-items: center;
  color: rgb(var(--v-theme-grey));
}
</style>
