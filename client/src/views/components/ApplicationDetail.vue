<!-- 受け取ったデータを基に申請詳細ページの下半分にログ、コメント等配置 -->
<template>
  <div :class="$style.container">
    <div :class="$style.header">
      <div :class="$style.title">
        <h1>{{ applicationType(detailCore.current_detail.type) }}申請</h1>
        <state-chip :state="detailCore.current_state" />
      </div>
      <state-button-controller />
    </div>

    <div :class="$style.section">
      <div :class="$style.section_title">申請日</div>
      <div :class="$style.section_item">
        {{ dayPrint(detailCore.created_at) }}
      </div>
    </div>

    <div :class="$style.section">
      <div :class="$style.section_title">申請者</div>
      <div :class="$style.section_item">
        <Icon :user="detailCore.applicant.trap_id" :size="24" />
        {{ detailCore.applicant.trap_id }}
      </div>
    </div>

    <div :class="$style.section">
      <div :class="$style.section_title">概要</div>
      <div :class="$style.section_item">
        {{ detailCore.current_detail.title }}
      </div>
    </div>

    <div :class="$style.section">
      <div :class="$style.section_title">支払日</div>
      <div :class="$style.section_item">
        {{ dayPrint(detailCore.current_detail.paid_at) }}
      </div>
    </div>

    <div :class="$style.section">
      <div :class="$style.section_title">支払金額</div>
      <div :class="$style.section_item">
        {{ detailCore.current_detail.amount }}円
      </div>
    </div>

    <div :class="$style.section">
      <div :class="$style.section_title">払い戻し対象者</div>
      <div :class="$style.target_container">
        <div
          v-for="user in detailCore.repayment_logs"
          :key="user.repaid_to_user.trap_id"
        >
          <Icon :user="user.repaid_to_user.trap_id" :size="24" />
          {{ user.repaid_to_user.trap_id }}
        </div>
      </div>
    </div>

    <div :class="$style.section">
      <div :class="$style.section_title">
        {{ remarksTitle(detailCore.current_detail.type) }}
      </div>
      <!-- eslint-disable-next-line vue/no-v-html -->
      <div :class="$style.section_item" v-html="rendered" />
    </div>

    <!-- 最後のsectionなのでスタイルを当てなくてもOK -->
    <div>
      <div :class="$style.section_title">画像</div>
      <div :class="$style.image_container">
        <img
          v-for="path in detailCore.images"
          :key="path"
          :src="`/api/images/${path}`"
        />
      </div>
      <div v-if="detailCore.images.length === 0">画像はありません</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useApplicationDetailStore } from "@/stores/applicationDetail";
import { applicationType, remarksTitle } from "@/use/applicationDetail";
import { dayPrint } from "@/use/dataFormat";
import { render } from "@/use/markdown";
import StateButtonController from "@/views/components/StateButtonController.vue";
import Icon from "@/views/shared/Icon.vue";
import StateChip from "@/views/shared/StateChip.vue";
import { storeToRefs } from "pinia";
import { onMounted, ref, watch } from "vue";

const applicationDetailStore = useApplicationDetailStore();
const { core: detailCore } = storeToRefs(applicationDetailStore);

const rendered = ref("");

// Watch for changes in remarks to re-render markdown
watch(
  () => detailCore.value.current_detail.remarks,
  async newVal => {
    rendered.value = await render(newVal);
  }
);

onMounted(async () => {
  if (detailCore.value.current_detail.remarks) {
    rendered.value = await render(detailCore.value.current_detail.remarks);
  }
});
</script>

<style lang="scss" module>
.container {
  height: fit-content;
  margin: 12px;
  padding: 12px;
  box-shadow:
    0 3px 1px -2px rgb(0 0 0 / 20%),
    0 2px 2px 0 rgb(0 0 0 / 14%),
    0 1px 5px 0 rgb(0 0 0 / 12%);
}
.header {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
}
.title {
  display: flex;
  align-items: center;
}
.section {
  margin: 16px 0;
  border-bottom: 1px solid $color-grey;
}
.section_title {
  color: $color-text-primary-disabled;
}
.section_item {
  margin-left: 8px;
}
.target_container {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(min(100%, 240px), 1fr));
  gap: 16px;
  padding: 8px;
}
.image_container {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(min(100%, 360px), 1fr));
  gap: 16px;
  padding: 8px;
}
</style>
