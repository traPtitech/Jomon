<template>
  <div>
    <v-dialog v-model="dialog" max-width="500px">
      <template #activator="{ props }">
        <simple-button
          :label="'払い戻し済みのユーザーを選択'"
          :variant="'done'"
          v-bind="props"
        />
      </template>
      <v-card :class="$style.container">
        <h3>払い戻し日</h3>
        <v-date-picker
          v-model="date"
          full-width
          flat
          @update:model-value="menu = false"
        />
        <v-autocomplete
          ref="traPID"
          v-model="traPID"
          :rules="[
            traPID =>
              traPID.length > 0 ||
              '払い戻し済みのユーザーが一人以上選ばれている必要があります'
          ]"
          :items="repaidToTraPId"
          label="払い戻し済みのユーザーを選択"
          :menu-props="{ zIndex: 3000 }"
          required
          multiple
        />
        <simple-button
          :label="'OK'"
          :disabled="traPID.length === 0"
          :variant="'done'"
          @click="putRepaid(traPID, date)"
        />
      </v-card>
    </v-dialog>
    <span v-if="repaidToTraPId.length === 0">
      何かがおかしいです。一度リロードしなおしてみて下さい。
    </span>
  </div>
</template>
<script setup lang="ts">
import { useApplicationDetailStore } from "@/stores/applicationDetail";
import SimpleButton from "@/views/shared/SimpleButton.vue";
import axios from "axios";
import { storeToRefs } from "pinia";
import { computed, ref } from "vue";

const applicationDetailStore = useApplicationDetailStore();
const { core: detailCore } = storeToRefs(applicationDetailStore);
const { fetchApplicationDetail } = applicationDetailStore;

const date = ref(new Date().toISOString().substr(0, 10));
const menu = ref(false);
const dialog = ref(false);
const traPID = ref<string[]>([]);

const repaidToTraPId = computed(() => {
  const trap_ids: string[] = [];
  detailCore.value.repayment_logs.forEach(log => {
    if (log.repaid_at === "" || log.repaid_at === null) {
      trap_ids.push(log.repaid_to_user.trap_id);
    }
  });
  return trap_ids;
});

const putRepaid = async (traPIDs: string[], date: string) => {
  await Promise.all(
    traPIDs.map(async traPID => {
      await axios
        .put(
          "/api/applications/" +
            detailCore.value.application_id +
            "/states/repaid/" +
            traPID,
          {
            repaid_at: date
          }
        )
        .catch(e => alert(e));
    })
  ).then(async () => {
    traPID.value = [];
    dialog.value = false;
    await fetchApplicationDetail(detailCore.value.application_id);
  });
};
</script>

<style lang="scss" module>
.container {
  min-width: 280px;
  padding: 8px;
}
</style>
