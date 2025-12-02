<template>
  <v-timeline-item color="purple">
    <div :class="$style.text">
      <icon :user="log.content.user.trap_id" :size="24" />
      {{ log.content.user.trap_id }}
      がコメントしました。
      <formatted-date :date="log.content.created_at" :simple="true" />
    </div>
    <div v-if="log.content.user.trap_id === $store.state.me.trap_id">
      <v-btn icon color="success" :disabled="!comment_readonly">
        <v-icon @click="commentChange()"> mdi-pencil </v-icon>
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
        ref="comment"
        v-model="comment_change"
        label="変更後のコメント"
        :readonly="comment_readonly"
        :solo="comment_readonly"
        :autofocus="!comment_readonly"
        flat
        dense
        hide-details
        :rules="changeRules"
        rows="1"
        auto-grow
      />

      <div>
        <v-btn
          v-if="!comment_readonly"
          :class="$style.button"
          @click="cancelChange"
        >
          変更を取消
        </v-btn>
        <v-btn
          v-if="!comment_readonly"
          :class="$style.button"
          :disabled="!comment_valid"
          @click="putComment"
        >
          変更を送信
        </v-btn>
      </div>

      <span
        v-if="log.content.created_at !== log.content.updated_at"
        :class="grey_text"
        >編集済</span
      >
    </v-form>
  </v-timeline-item>
</template>

<script>
import Icon from "@/views/shared/Icon";
import axios from "axios";
import { nextTick } from "vue";
import { mapActions } from "vuex";
import FormattedDate from "./FormattedDate";

export default {
  components: {
    Icon,
    FormattedDate
  },
  props: {
    log: {
      type: Object,
      default: () => ({})
    }
  },
  data: function () {
    return {
      comment_readonly: true,
      comment_change: this.log.content.comment,
      comment_valid: true,
      changeRules: [v => v !== this.log.content.comment && !!v]
    };
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

<style lang="scss" module>
.text {
  display: flex;
  color: $color-grey;
}
.button {
  margin: 8px 8px 0 0;
}
</style>
