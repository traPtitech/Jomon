<!-- ここはログの最小単位の表示を制御-->
<template>
  <v-timeline-item v-if="log.log_type === `comment`" color="purple lighten-2">
    <template v-slot:icon>
      <span>Com</span>
    </template>
    <v-card>
      <v-card-title class="headline">
        <v-row>
          <v-col cols="8">
            <Icon :user="log.content.user.trap_id" :size="25" />{{
              log.content.user.trap_id
            }}がコメントしました。</v-col
          ><v-col cols="2">{{ dayPrint(log.content.created_at) }}</v-col
          ><v-col cols="2">
            <v-btn class="ma-2" tile outlined color="success">
              <v-icon left>mdi-pencil</v-icon> Edit</v-btn
            >
          </v-col>
        </v-row>
      </v-card-title>
      <v-card-text>
        {{ log.content.comment }}
      </v-card-text>
    </v-card>
  </v-timeline-item>
  <!-- <div v-else-if="log.log_type === `state`"> -->
  <v-timeline-item
    v-else-if="log.log_type === `state`"
    color="red lighten-2"
    large
  >
    <v-row justify="space-between">
      <v-col cols="11">
        <Icon :user="log.content.update_user.trap_id" :size="25" />{{
          log.content.update_user.trap_id
        }}が申請書の状態を{{ log.content.to_state }}に変更しました。</v-col
      >
      <v-col class="text-right" cols="1">{{
        dayPrint(log.content.created_at)
      }}</v-col>
    </v-row>
    <v-card v-if="log.content.reason !== ''">
      <v-card-text>
        理由:
      </v-card-text>
      {{ log.content.reason }}
    </v-card>
  </v-timeline-item>

  <div v-else-if="log.log_type === `application`">
    <v-timeline-item color="purple lighten-2">
      <template v-slot:icon>
        <span>Application </span>
      </template>
      <v-row justify="space-between">
        <v-col cols="11">
          <Icon :user="log.content.update_user.trap_id" :size="25" />{{
            log.content.update_user.trap_id
          }}が修正しました。</v-col
        >
        <v-col class="text-right" cols="1">{{
          dayPrint(log.content.updated_at)
        }}</v-col>
      </v-row>
    </v-timeline-item>
    <v-card>
      実際は修正された部分だけを示すので、(ex:-hoge
      +piyo)そのロジックをどっかに書きます。 type:{{
        log.content.type
      }}
      title:{{ log.content.title }} remarks: {{ log.content.remarks }} ammount
      {{ log.content.ammount }} paid_at{{ log.content.paid_at }}
    </v-card>
  </div>

  <v-timeline-item
    v-else-if="log.log_type === `repayment`"
    class="mb-4"
    color="grey"
    icon-color="grey lighten-2"
    small
  >
    <v-row justify="space-between">
      <v-col cols="11">
        <Icon :user="log.content.repaid_by_user.trap_id" :size="25" />{{
          log.content.repaid_by_user.trap_id
        }}が{{
          log.content.repaid_to_user.trap_id
        }}に払い戻しをしました。</v-col
      >
      <v-col class="text-right" cols="1">{{
        dayPrint(log.content.repaid_at)
      }}</v-col>
    </v-row>
  </v-timeline-item>
</template>

<script>
import Icon from "./Icon";
export default {
  props: {
    log: Object
  },
  components: {
    Icon
  },
  methods: {
    dayPrint(time) {
      //let d = this.log.content.created_at;
      let d = new Date(time);
      //let year = hiduke.getFullYear();
      let month = d.getMonth() + 1;
      //let week = hiduke.getDay();
      let day = d.getDate();
      let res = month + "/" + day;
      return res;
    }
  }
};
</script>
