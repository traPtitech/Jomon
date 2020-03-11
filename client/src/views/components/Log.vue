<!-- ここはログの最小単位の表示を制御-->
<template>
  <!-- 以下はコメントログ -->
  <v-timeline-item v-if="log.log_type === `comment`" color="purple lighten-2">
    <template v-slot:icon>
      <Icon :user="log.content.user.trap_id" :size="25" />
    </template>
    <v-card class="pa-2">
      <v-row>
        <v-col class="pb-0 pt-0" :class="grey_text" cols="8">
          <strong :class="strong_text"> {{ log.content.user.trap_id }}</strong>
          <span :class="larger_size">がコメントしました。</span
          ><span :class="smaller_size">のコメント</span>
        </v-col>
        <v-col class="pa-0">{{ dayPrint(log.content.created_at) }}</v-col>
        <v-col class="pa-0" cols="1">
          <v-btn icon color="success">
            <v-icon left>mdi-pencil</v-icon>
          </v-btn>
        </v-col>
        <v-col class="pa-0" cols="1">
          <v-btn icon color="error">
            <v-icon>mdi-delete</v-icon>
          </v-btn>
        </v-col>
      </v-row>

      <v-divider></v-divider>
      <v-card-text class="pa-0 black--text">
        {{ log.content.comment }}
      </v-card-text>
    </v-card>
  </v-timeline-item>
  <!-- 以下は申請書の状態変化ログ -->
  <v-timeline-item
    v-else-if="log.log_type === `state`"
    color="red lighten-2"
    large
  >
    <v-row>
      <v-col cols="10" :class="grey_text">
        <Icon :user="log.content.update_user.trap_id" :size="25" /><strong
          :class="strong_text"
          >{{ log.content.update_user.trap_id }}</strong
        >が
        <span :class="larger_size"
          >申請書の状態を<StateChip
            :state="log.content.to_state"
          />に変更しました。</span
        >
        <span :class="smaller_size"
          ><StateChip :state="log.content.to_state" size="small" />に変更</span
        >
      </v-col>
      <v-col class="text-right" cols="2">{{
        dayPrint(log.content.created_at)
      }}</v-col>
    </v-row>
    <v-card v-if="log.content.reason !== ''">
      <v-card-text>
        理由:
        <span class="black--text">{{ log.content.reason }}</span>
      </v-card-text>
    </v-card>
  </v-timeline-item>
  <!-- 以下は申請書の変更ログ -->
  <div v-else-if="log.log_type === `application`">
    <v-timeline-item color="purple lighten-2">
      <template v-slot:icon>
        <span>Application </span>
      </template>
      <v-row justify="space-between">
        <v-col cols="10" :class="grey_text">
          <Icon :user="log.content.update_user.trap_id" :size="25" /><strong
            :class="strong_text"
            >{{ log.content.update_user.trap_id }}</strong
          >が修正<span :class="larger_size">しました。</span></v-col
        >
        <v-col class="text-right" cols="2">{{
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
  <!-- 以下は払い戻しログ -->
  <v-timeline-item
    v-else-if="log.log_type === `repayment`"
    class="mb-4"
    color="grey"
    icon-color="grey lighten-2"
    small
  >
    <v-row justify="space-between">
      <v-col cols="10" :class="grey_text">
        <Icon :user="log.content.repaid_by_user.trap_id" :size="25" />
        <strong :class="strong_text">{{
          log.content.repaid_by_user.trap_id
        }}</strong
        >が
        <strong :class="strong_text">{{
          log.content.repaid_by_user.trap_id
        }}</strong
        >に 払い戻し<span :class="larger_size">をしました。</span>
      </v-col>
      <v-col class="text-right" cols="2">{{
        dayPrint(log.content.repaid_at)
      }}</v-col>
    </v-row>
  </v-timeline-item>
</template>

<script>
import Icon from "./Icon";
import StateChip from "./StateChip";
export default {
  data: function() {
    return {
      smaller_size: "hidden-lg-and-up",
      larger_size: "hidden-md-and-down",
      grey_text: "grey--text text--darken-1 body-2",
      strong_text: "black--text subtitle-1"
    };
  },
  props: {
    log: Object
  },
  components: {
    Icon,
    StateChip
  },
  methods: {
    dayPrint(time) {
      let d = new Date(time);
      let month = d.getMonth() + 1;
      let day = d.getDate();
      let res = month + "/" + day;
      return res;
    }
  }
};
</script>
