<template>
  <div v-if="loading">loading...</div>
  <div v-else-if="error">
    <v-alert type="error"> 申請情報の取得に失敗しました。 </v-alert>
  </div>
  <div v-else :class="$style.container">
    <application-paper v-if="!fix" :class="$style.paper" />
    <fix-application-paper v-else :class="$style.paper" />
    <application-logs :class="$style.log" />
  </div>
</template>

<script setup lang="ts">
import { useApplicationDetailStore } from "@/stores/applicationDetail";
import ApplicationPaper from "@/views/components/ApplicationDetail.vue";
import ApplicationLogs from "@/views/components/ApplicationDetailLogs.vue";
import FixApplicationPaper from "@/views/components/FixApplicationDetail.vue";
import { storeToRefs } from "pinia";
import { onMounted, ref } from "vue";
import { useRoute } from "vue-router";

const route = useRoute();
const applicationDetailStore = useApplicationDetailStore();
const { fix } = storeToRefs(applicationDetailStore);
const { fetchApplicationDetail, deleteFix } = applicationDetailStore;

const loading = ref(true);

const error = ref(false);

onMounted(async () => {
  try {
    await fetchApplicationDetail(route.params.id as string);
    deleteFix();
  } catch (e) {
    console.error(e);
    error.value = true;
  } finally {
    loading.value = false;
  }
});
</script>

<style lang="scss" module>
.container {
  display: flex;
  justify-content: center;
  flex-wrap: wrap;
  max-width: 100%;
  padding: 16px;
}
.paper {
  width: 720px;
}
.log {
  min-width: 360px;
  max-width: 480px;
}
@media (max-width: $breakpoint) {
  .paper {
    max-width: 100%;
  }
}
</style>
