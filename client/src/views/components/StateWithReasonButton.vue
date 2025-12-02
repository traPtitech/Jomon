<template>
  <div :class="$style.container">
    <v-dialog v-model="isDialogOpen" max-width="600px">
      <template #activator="{ props: activatorProps }">
        <simple-button
          :label="
            toStateName(toState) + (toState === 'submitted' ? 'に戻す' : '')
          "
          :variant="toState === 'submitted' ? 'warning' : 'error'"
          v-bind="activatorProps"
        />
      </template>

      <v-card>
        <v-card-title>
          <span v-if="toState === `submitted`" class="headline"
            >承認済み→{{ toStateName(toState) }} へ戻す理由</span
          >
          <span v-else class="headline"
            >承認待ち→{{ toStateName(toState) }} への変更理由</span
          >
        </v-card-title>
        <v-form ref="formRef" v-model="valid">
          <v-card-text>
            <v-container>
              <v-row>
                <v-col cols="12">
                  <v-text-field
                    ref="reason"
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
            <simple-button
              :label="'戻る'"
              :variant="'secondary'"
              @click="isDialogOpen = false"
            />
            <simple-button
              :label="toStateName(toState) + 'にする'"
              :variant="'secondary'"
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
import SimpleButton from "@/views/shared/SimpleButton.vue";
import axios from "axios";
import { storeToRefs } from "pinia";
import { nextTick, ref, watch } from "vue";

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

const valid = ref(true);
const isDialogOpen = ref(false);
const reason = ref("");
const nullRules = [(v: unknown) => !!v || ""];
const formRef = ref<any>(null); // eslint-disable-line @typescript-eslint/no-explicit-any

watch(isDialogOpen, async () => {
  if (isDialogOpen.value) {
    await nextTick();
    formRef.value.reset();
  }
});

const blur = () => {
  if (reason.value === "" || reason.value === undefined) {
    formRef.value.reset();
  }
};

const postReason = async () => {
  if (formRef.value.validate()) {
    try {
      await axios.put(
        "../api/applications/" + detailCore.value.application_id + "/states",
        {
          to_state: props.toState,
          reason: reason.value
        }
      );
    } catch (e) {
      alert(e);
      return;
    }
    formRef.value.reset();
    isDialogOpen.value = false;
    await fetchApplicationDetail(detailCore.value.application_id);
  }
};

const toStateName = (toState: string) => {
  switch (toState) {
    case "submitted":
      return "提出済み";
    case "fix_required":
      return "要修正";
    case "rejected":
      return "取り下げ";
    default:
      return "状態が間違っています";
  }
};
</script>

<style lang="scss" module>
.container {
  justify-content: center;
}
</style>
