<template>
  <v-timeline-item color="purple">
    <div :class="$style.text">
      <icon :user="log.content.log.update_user.trap_id" :size="25" />
      {{ `${log.content.log.update_user.trap_id}が修正しました。` }}
      <formatted-date :date="log.content.log.updated_at" />
    </div>
    <v-card>
      <defference
        v-if="log.content.pre_log.type !== log.content.log.type"
        item="type"
        :pre="log.content.pre_log.type"
        :now="log.content.log.type"
      />
      <defference
        v-if="log.content.pre_log.title !== log.content.log.title"
        item="title"
        :pre="log.content.pre_log.title"
        :now="log.content.log.title"
      />
      <defference
        v-if="log.content.pre_log.remarks !== log.content.log.remarks"
        item="remarks"
        :pre="log.content.pre_log.remarks"
        :now="log.content.log.remarks"
      />
      <defference
        v-if="log.content.pre_log.amount !== log.content.log.amount"
        item="amount"
        :pre="log.content.pre_log.amount"
        :now="log.content.log.amount"
      />
      <defference
        v-if="log.content.pre_log.paid_at !== log.content.log.paid_at"
        item="paid_at"
        :pre="log.content.pre_log.paid_at"
        :now="log.content.log.paid_at"
      />
    </v-card>
  </v-timeline-item>
</template>

<script>
import Icon from "@/views/shared/Icon";
import Defference from "./ApplicationDetailDefference";
import FormattedDate from "./FormattedDate";

export default {
  props: {
    log: Object
  },
  components: {
    Icon,
    Defference,
    FormattedDate
  }
};
</script>
<style lang="scss" module>
.text {
  display: flex;
  flex-wrap: wrap;
  max-width: 100%;
}
</style>
