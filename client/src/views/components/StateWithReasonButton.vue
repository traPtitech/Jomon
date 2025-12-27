<template>
  <div :class="$style.container">
    <v-dialog v-model="isDialogOpen" max-width="600px">
      <template #activator="{ props: activatorProps }">
        <simple-button
          :label="
            getStateLabel(toState) + (toState === 'submitted' ? 'に戻す' : '')
          "
          :variant="getStateButtonColor(toState)"
          v-bind="activatorProps"
        />
      </template>

      <v-card>
        <v-card-title>
          <span v-if="toState === `submitted`" class="headline"
            >承認済み→{{ getStateLabel(toState) }} へ戻す理由</span
          >
          <span v-else class="headline"
            >承認待ち→{{ getStateLabel(toState) }} への変更理由</span
          >
        </v-card-title>
        <v-form ref="formRef" v-model="valid">
          <v-card-text>
            <v-container>
              <v-row>
                <v-col cols="12">
                  <v-text-field
                    v-model="reason"
                    :autofocus="isDialogOpen"
                    :rules="nullRules"
                    @blur="blur"
                  />
                </v-col>
              </v-row>
            </v-container>
          </v-card-text>
          <v-card-actions>
            <v-spacer />
            <simple-button :label="'戻る'" @click="isDialogOpen = false" />
            <simple-button
              :label="getStateLabel(toState) + 'にする'"
              :disabled="!valid"
              @click="postReason"
            />
          </v-card-actions>
        </v-form>
      </v-card>
    </v-dialog>
  </div>
</template>
<script setup lang="ts">
import { useApplicationDetailStore } from "@/stores/applicationDetail";
import { getStateButtonColor, getStateLabel } from "@/use/stateColor";
import SimpleButton from "@/views/shared/SimpleButton.vue";
import axios from "axios";
import { storeToRefs } from "pinia";
import { nextTick, ref, useTemplateRef, watch } from "vue";

const props = withDefaults(
  defineProps<{
    toState?: string;
  }>(),
  {
    toState: ""
  }
);

const applicationDetailStore = useApplicationDetailStore();
const { core: detailCore } = storeToRefs(applicationDetailStore);
const { fetchApplicationDetail } = applicationDetailStore;

import { VForm } from "vuetify/components";

const valid = ref(true);
const isDialogOpen = ref(false);
const reason = ref("");
const nullRules = [(v: unknown) => !!v || ""];
const formRef = useTemplateRef<VForm>("formRef");

watch(isDialogOpen, async () => {
  if (isDialogOpen.value) {
    await nextTick();
    reason.value = "";
    formRef.value?.resetValidation();
  }
});

const blur = () => {
  if (reason.value === "" || reason.value === undefined) {
    formRef.value?.resetValidation();
  }
};

const postReason = async () => {
  const validateResult = await formRef.value?.validate();
  if (validateResult?.valid) {
    try {
      await axios.put(
        "/api/applications/" + detailCore.value.application_id + "/states",
        {
          to_state: props.toState,
          reason: reason.value
        }
      );
    } catch (e) {
      alert(e);
      return;
    }
    reason.value = "";
    formRef.value?.resetValidation();
    isDialogOpen.value = false;
    await fetchApplicationDetail(detailCore.value.application_id);
  }
};
</script>

<style lang="scss" module>
.container {
  justify-content: center;
}
</style>
