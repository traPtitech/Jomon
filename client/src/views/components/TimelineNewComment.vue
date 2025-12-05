<template>
  <div>
    <v-row justify="space-between">
      <v-col cols="1">
        <Icon :user="trapId" :size="25" />
      </v-col>
      <v-col cols="11">
        <v-card class="pa-2">
          <v-form ref="formRef" v-model="valid">
            <v-textarea
              v-model="comment"
              :rules="nullRules"
              outlined
              label="コメントを書いてください"
              @blur="blur"
            />
            <v-btn
              :disabled="!valid"
              color="primary"
              class="mr-4"
              @click="postcomment"
            >
              コメントする
            </v-btn>
          </v-form>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>
<script setup lang="ts">
import { useApplicationDetailStore } from "@/stores/applicationDetail";
import { useMeStore } from "@/stores/me";
import Icon from "@/views/shared/Icon.vue";
import axios from "axios";
import { storeToRefs } from "pinia";
import { ref, useTemplateRef } from "vue";

const applicationDetailStore = useApplicationDetailStore();
const meStore = useMeStore();

const { core: detailCore } = storeToRefs(applicationDetailStore);
const { trapId } = storeToRefs(meStore);
const { fetchApplicationDetail } = applicationDetailStore;

import { VForm } from "vuetify/components";

const valid = ref(true);
const comment = ref("");
const nullRules = [(v: unknown) => !!v || ""];
const formRef = useTemplateRef<VForm>("formRef");

const blur = () => {
  if (comment.value === "" || comment.value === undefined) {
    formRef.value?.reset();
  }
};

const postcomment = async () => {
  const validateResult = await formRef.value?.validate();
  if (validateResult?.valid) {
    try {
      await axios.post(
        "/api/applications/" + detailCore.value.application_id + "/comments",
        {
          comment: comment.value
        }
      );
    } catch (e) {
      alert(e);
      return;
    }
    formRef.value?.reset();
    await fetchApplicationDetail(detailCore.value.application_id);
  }
};
</script>
