<template>
  <div v-if="displayAcceptBottom" :class="$style.button_container">
    <simple-button :label="'承認'" :variant="'success'" @click="accept()" />
    <with-reason-button to-state="fix_required" />
    <with-reason-button to-state="rejected" />
  </div>
  <div v-else-if="displayRepaidBottom" :class="$style.button_container">
    <with-reason-button to-state="submitted" />
    <repaid-button />
  </div>
  <div v-else-if="displayFixResubmitBottom" :class="$style.button_container">
    <simple-button :label="'修正'" @click="changeFix" />
    <simple-button :label="'再申請'" @click="reSubmit" />
  </div>
</template>
<script setup lang="ts">
import { useApplicationDetailStore } from "@/stores/applicationDetail";
import { useMeStore } from "@/stores/me";
import SimpleButton from "@/views/shared/SimpleButton.vue";
import { storeToRefs } from "pinia";
import { computed } from "vue";
import RepaidButton from "./RepaidButton.vue";
import WithReasonButton from "./StateWithReasonButton.vue";

const applicationDetailStore = useApplicationDetailStore();
const meStore = useMeStore();

const { core: detailCore } = storeToRefs(applicationDetailStore);
const { trapId, isAdmin } = storeToRefs(meStore);
const { changeFix, updateApplicationState } = applicationDetailStore;

const displayAcceptBottom = computed(() => {
  return detailCore.value.current_state === `submitted` && isAdmin.value;
});

const displayRepaidBottom = computed(() => {
  return detailCore.value.current_state === `accepted` && isAdmin.value;
});

const displayFixResubmitBottom = computed(() => {
  return (
    detailCore.value.current_state === `fix_required` &&
    (isAdmin.value || trapId.value === detailCore.value.applicant.trap_id)
  );
});

const accept = async () => {
  try {
    await updateApplicationState("accepted");
  } catch (e) {
    alert(e);
    return;
  }
  alert("承認しました");
};

const reSubmit = async () => {
  try {
    await updateApplicationState("submitted");
  } catch (e) {
    alert(e);
    return;
  }
  alert("再申請しました");
};
</script>

<style lang="scss" module>
.button_container {
  display: flex;
}
</style>
