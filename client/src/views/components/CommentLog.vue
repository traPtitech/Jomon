<template>
  <v-timeline-item color="purple">
    <div :class="$style.text">
      <icon :user="log.content.user.trap_id" :size="24" />
      {{ log.content.user.trap_id }}
      がコメントしました。
      <formatted-date :date="log.content.created_at" :simple="true" />
    </div>
    <div v-if="log.content.user.trap_id === trapId">
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
        ref="commentRef"
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
        :class="$style.grey_text"
        >編集済</span
      >
    </v-form>
  </v-timeline-item>
</template>

<script setup lang="ts">
import { useApplicationDetailStore } from "@/stores/applicationDetail";
import { useMeStore } from "@/stores/me";
import { CommentLog } from "@/types/log";
import Icon from "@/views/shared/Icon.vue";
import axios from "axios";
import { storeToRefs } from "pinia";
import { nextTick, ref, watch } from "vue";
import FormattedDate from "./FormattedDate.vue";

const props = defineProps<{
  log: CommentLog;
}>();

const applicationDetailStore = useApplicationDetailStore();
const meStore = useMeStore();

const { core: detailCore } = storeToRefs(applicationDetailStore);
const { trapId } = storeToRefs(meStore);
const { fetchApplicationDetail } = applicationDetailStore;

const comment_readonly = ref(true);
const comment_change = ref(props.log.content.comment);
const comment_valid = ref(true);
const changeRules = [
  (v: unknown) => (v !== props.log.content.comment && !!v) || ""
];

import { VTextarea } from "vuetify/components";

const commentRef = ref<VTextarea | null>(null);

watch(comment_readonly, async () => {
  if (!comment_readonly.value) {
    await nextTick();
    commentRef.value.focus();
  }
});

const commentChange = () => {
  comment_readonly.value = false;
};

const deleteComment = async () => {
  try {
    await axios.delete(
      "/api/applications/" +
        detailCore.value.application_id +
        "/comments/" +
        props.log.content.comment_id
    );
  } catch (e) {
    alert(e);
    return;
  }
  alert("コメントを削除しました。");
  await fetchApplicationDetail(detailCore.value.application_id);
};

const putComment = async () => {
  try {
    await axios.put(
      "/api/applications/" +
        detailCore.value.application_id +
        "/comments/" +
        props.log.content.comment_id,
      {
        comment: comment_change.value
      }
    );
  } catch (e) {
    alert(e);
    return;
  }
  alert("コメントを変更しました");
  comment_readonly.value = true;
  await fetchApplicationDetail(detailCore.value.application_id);
};

const cancelChange = () => {
  comment_readonly.value = true;
  comment_change.value = props.log.content.comment;
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
.grey_text {
  color: $color-grey;
}
</style>
