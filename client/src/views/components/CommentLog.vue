<template>
  <v-timeline-item color="purple">
    <v-card>
      <icon :user="log.content.user.trap_id" :size="24" />
      {{ log.content.user.trap_id }}
      がコメントしました。

      <formatted-date :date="log.content.created_at" :simple="true" />

      <div v-if="log.content.user.trap_id === this.$store.state.me.trap_id">
        <v-btn icon color="success" :disabled="!comment_readonly">
          <v-icon left @click="commentChange()">mdi-pencil</v-icon>
        </v-btn>
        <v-btn
          icon
          color="error"
          :disabled="!comment_readonly"
          @click="deleteComment()"
        >
          <v-icon>mdi-delete</v-icon>
        </v-btn>
      </div>

      <v-form v-model="comment_valid">
        <v-textarea
          v-model="comment_change"
          label="変更のコメント"
          :readonly="comment_readonly"
          :solo="comment_readonly"
          :autofocus="!comment_readonly"
          flat
          dense
          hide-details
          :rules="changeRules"
          rows="1"
          auto-grow
        >
        </v-textarea>

        <v-btn @click="cancelChange()" v-if="!comment_readonly"
          >変更を取消
        </v-btn>
        <v-btn
          @click="putComment"
          :disabled="!comment_valid"
          v-if="!comment_readonly"
          >変更を送信
        </v-btn>

        <span
          :class="grey_text"
          v-if="log.content.created_at !== log.content.updated_at"
          >編集済</span
        >
      </v-form>
    </v-card>
  </v-timeline-item>
</template>

<script>
import Icon from "@/views/shared/Icon";
import FormattedDate from "./FormattedDate";
import { nextTick } from "vue";
import { mapActions } from "vuex";
import axios from "axios";

export default {
  data: function () {
    return {
      comment_readonly: true,
      comment_change: this.log.content.comment,
      comment_valid: true,
      changeRules: [v => v !== this.log.content.comment && !!v]
    };
  },
  props: {
    log: Object
  },
  components: {
    Icon,
    FormattedDate
  },
  watch: {
    comment_readonly: function () {
      if (!this.comment_readonly) {
        let self = this;
        nextTick().then(function () {
          self.$refs.comment.focus();
        });
      }
    }
  },
  methods: {
    ...mapActions(["getApplicationDetail"]),
    commentChange() {
      this.comment_readonly = false;
    },
    async deleteComment() {
      await axios
        .delete(
          "../api/applications/" +
            this.$store.state.application_detail_paper.core.application_id +
            "/comments/" +
            this.log.content.comment_id
        )
        .catch(e => alert(e));
      alert("コメントを削除しました。");
      this.getApplicationDetail(
        this.$store.state.application_detail_paper.core.application_id
      );
    },
    async putComment() {
      await axios
        .put(
          "../api/applications/" +
            this.$store.state.application_detail_paper.core.application_id +
            "/comments/" +
            this.log.content.comment_id,
          {
            comment: this.comment_change
          }
        )
        .catch(e => {
          alert(e);
          return;
        });
      alert("コメントを変更しました");
      this.comment_readonly = true;
      this.getApplicationDetail(
        this.$store.state.application_detail_paper.core.application_id
      );
    },
    cancelChange() {
      this.comment_readonly = true;
      this.comment_change = this.log.content.comment;
    }
  }
};
</script>
