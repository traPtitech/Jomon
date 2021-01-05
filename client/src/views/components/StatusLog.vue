<template>
  <v-timeline-item color="red lighten-2" large>
    <icon :user="log.content.update_user.trap_id" :size="25" />
    {{ log.content.update_user.trap_id }}
    が 申請の状態を
    <state-chip :state="log.content.to_state" />に変更しました。
    {{ dayPrint(log.content.created_at) }}
    <v-card v-if="log.content.reason !== ''">
      <v-card-text> 理由: {{ log.content.reason }} </v-card-text>
    </v-card>
  </v-timeline-item>
</template>

<script>
import Icon from "@/views/shared/Icon";
import StateChip from "@/views/shared/StateChip";
import { mapActions } from "vuex";

export default {
  props: {
    log: Object
  },
  components: {
    Icon,
    StateChip
  },
  methods: {
    dayPrint(time) {
      let now = new Date();
      let d = new Date(time);
      let diff = (now.getTime() - d.getTime()) / 1000;
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
        let month = d.getMonth() + 1;
        let day = d.getDate();
        return month + "/" + day;
      } else {
        let year = d.getFullYear();
        let month = d.getMonth() + 1;
        let day = d.getDate();
        return year + "/" + month + "/" + day;
      }
    }
  }
};
</script>
