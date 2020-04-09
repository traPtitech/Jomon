<!-- ここはログの最小単位の表示を制御-->
<template>
  <!-- 以下はコメントログ -->
  <v-timeline-item v-if="log.log_type === `comment`" color="purple lighten-2">
    <template v-slot:icon>
      <Icon :user="log.content.user.trap_id" :size="25" />
    </template>
    <v-card class="pa-2">
      <v-row>
        <!-- todo update_atと比較して薄い灰色で編集済みの記述 -->
        <v-col class="pb-0 pt-0" :class="grey_text" cols="7">
          <strong :class="strong_text"> {{ log.content.user.trap_id }}</strong>
          <span :class="larger_size">がコメントしました。</span>
          <span :class="smaller_size">のコメント</span>
        </v-col>
        <v-col class="pa-0" cols="2">{{
          dayPrint(log.content.created_at)
        }}</v-col>
        <v-col
          v-if="log.content.user.trap_id === this.$store.state.me.trap_id"
          class="pa-0"
          cols="3"
        >
          <v-btn icon color="success" :disabled="!comment_readonly">
            <v-icon left @click="commentChange()">mdi-pencil</v-icon>
          </v-btn>

          <v-btn
            icon
            color="error"
            :disabled="!comment_readonly"
            v-on:click="deleteComment()"
          >
            <v-icon>mdi-delete</v-icon>
          </v-btn>
        </v-col>
      </v-row>

      <v-divider></v-divider>
      <!-- <v-card-text class="pa-0 black--text">
        {{ log.content.comment }}
      </v-card-text> -->
      <v-form ref="form" v-model="comment_valid">
        <!-- v-text-fieldから変更 -->
        <v-textarea
          ref="comment"
          v-model="comment_change"
          class="ma-0 black--text"
          :class="comment_readonly === false ? `pt-3 pl-3` : `pa-0`"
          label="変更のコメント"
          :readonly="comment_readonly"
          :solo="comment_readonly"
          :autofocus="!comment_readonly"
          flat
          dense
          hide-details
          :rules="changeRules"
        >
        </v-textarea>
        <span v-if="log.content.created_at !== log.content.updated_at"
          >編集済</span
        >
        <div v-if="!comment_readonly" class="pt-2">
          <v-btn
            @click="
              (comment_readonly = true), (comment_change = log.content.comment)
            "
            >変更を取消</v-btn
          ><v-btn @click="putComment" :disabled="!comment_valid"
            >変更を送信</v-btn
          >
        </div>
      </v-form>
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
        <Icon :user="log.content.update_user.trap_id" :size="25" />
        <strong :class="strong_text">
          {{ log.content.update_user.trap_id }}
        </strong>
        が
        <span :class="larger_size">
          申請書の状態を
          <StateChip :state="log.content.to_state" />に変更しました。
        </span>
        <span :class="smaller_size">
          <StateChip :state="log.content.to_state" size="small" />に変更
        </span>
      </v-col>
      <v-col class="text-right" cols="2">
        {{ dayPrint(log.content.created_at) }}
      </v-col>
    </v-row>
    <v-card v-if="log.content.reason !== ''">
      <v-card-text>
        理由:
        <span class="black--text">{{ log.content.reason }}</span>
      </v-card-text>
    </v-card>
  </v-timeline-item>
  <!-- 以下は申請書の変更ログ -->
  <div v-else-if="log.log_type === `application`" class="ma-2">
    <v-timeline-item color="purple lighten-2">
      <template v-slot:icon>
        <span>Application </span>
      </template>
      <v-row justify="space-between">
        <v-col cols="10" :class="grey_text">
          <Icon :user="log.content.log.update_user.trap_id" :size="25" />
          <strong :class="strong_text">
            {{ log.content.log.update_user.trap_id }}
          </strong>
          が修正
          <span :class="larger_size">しました。</span>
        </v-col>
        <v-col class="text-right" cols="2">{{
          dayPrint(log.content.log.updated_at)
        }}</v-col>
      </v-row>
    </v-timeline-item>
    <v-card>
      <div style="white-space: pre-line" class="pa-2 ma-0">
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
      </div>
    </v-card>
  </div>
  <!-- 以下は払い戻しログ -->
  <v-timeline-item
    v-else-if="
      log.log_type === `repayment` &&
        !(log.content.repaid_at === `` || log.content.repaid_at === null)
    "
    class="mb-4"
    color="grey"
    icon-color="grey lighten-2"
    small
  >
    <v-row justify="space-between">
      <v-col cols="10" :class="grey_text">
        <Icon :user="log.content.repaid_by_user.trap_id" :size="25" />
        <strong :class="strong_text">
          {{ log.content.repaid_by_user.trap_id }}
        </strong>
        が
        <strong :class="strong_text">
          {{ log.content.repaid_to_user.trap_id }}
        </strong>
        に 払い戻し
        <span :class="larger_size">をしました。</span>
      </v-col>
      <v-col class="text-right" cols="2">
        {{ dayPrint(log.content.repaid_at) }}
      </v-col>
    </v-row>
  </v-timeline-item>
</template>

<script>
import Icon from "./Icon";
import StateChip from "./StateChip";
import Vue from "vue";
import axios from "axios";
import Defference from "./ApplicationDetailDefference";
export default {
  data: function() {
    return {
      smaller_size: "hidden-lg-and-up",
      larger_size: "hidden-md-and-down",
      grey_text: "grey--text text--darken-1 body-2",
      strong_text: "black--text subtitle-1",
      comment_readonly: true,
      comment_change: this.log.content.comment,
      comment_valid: true,
      changeRules: [v => (v !== this.log.content.comment && !!v) || ""]
    };
  },
  props: {
    log: Object
  },
  components: {
    Icon,
    StateChip,
    Defference
  },
  watch: {
    comment_readonly: function() {
      if (!this.comment_readonly) {
        let self = this;
        Vue.nextTick().then(function() {
          self.$refs.comment.focus();
        });
      }
    }
  },
  methods: {
    defferenceRecord(pre, now) {
      let defference_record =
        pre.type !== now.type
          ? "申請書:" + pre.type + "から" + now.type + "\n"
          : "";
      defference_record +=
        pre.title !== now.title
          ? "タイトル:" + pre.title + "から" + now.title + "\n"
          : "";
      defference_record +=
        pre.remarks !== now.remarks
          ? "詳細:" + pre.remarks + "から" + now.remarks + "\n"
          : "";
      defference_record +=
        pre.amount !== now.amount
          ? "金額:" + pre.amount + "から" + now.amount + "\n"
          : "";
      defference_record +=
        pre.paid_at !== now.paid_at
          ? "支払日:" +
            this.dayPrint(pre.paid_at) +
            "から" +
            this.dayPrint(now.paid_at) +
            "\n"
          : "";
      return defference_record;
    },
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
        let res = month + "/" + day;
        return res;
      } else {
        let year = d.getFullYear();
        let month = d.getMonth() + 1;
        let day = d.getDate();
        let res = year + "/" + month + "/" + day;
        return res;
      }
    },
    commentChange() {
      this.comment_readonly = false;
    },
    deleteComment() {
      axios
        .delete(
          "../api/applications/" +
            this.$store.state.application_detail_paper.application_id +
            "/comments/" +
            this.log.content.comment_id
        )
        .then(response => console.log(response.status));
      alert("コメントを削除しました。");
    },
    putComment() {
      axios
        .put(
          "../api/applications/" +
            this.$store.state.application_detail_paper.application_id +
            "/comments/" +
            this.log.content.comment_id,
          {
            comment: this.comment_change
          }
        )
        .then(response => console.log(response.status));
      alert("コメントを変更しました");
      this.comment_readonly = true;
      this.comment_change = this.log.content.comment;
    }
  }
};
</script>
