<!-- 受け取ったデータを基に申請詳細ページの下半分にログ、コメント等配置 -->
<template>
  <div :class="$style.container">
    <div :class="$style.header">
      <div :class="$style.title">
        <h1>{{ returnType(detail.core.current_detail.type) }}申請</h1>
        <state-chip :state="detail.core.current_state" />
      </div>
      <state-button-controller />
    </div>

    <div :class="$style.section">
      <div :class="$style.section_title">申請日</div>
      <div :class="$style.section_item">
        {{ returnDate(detail.core.created_at) }}
      </div>
    </div>

    <div :class="$style.section">
      <div :class="$style.section_title">申請者</div>
      <div :class="$style.section_item">
        <Icon :user="detail.core.applicant.trap_id" :size="24" />
        {{ detail.core.applicant.trap_id }}
      </div>
    </div>

    <div :class="$style.section">
      <div :class="$style.section_title">概要</div>
      <div :class="$style.section_item">
        {{ detail.core.current_detail.title }}
      </div>
    </div>

    <div :class="$style.section">
      <div :class="$style.section_title">支払日</div>
      <div :class="$style.section_item">
        {{ returnDate(detail.core.current_detail.paid_at) }}
      </div>
    </div>

    <div :class="$style.section">
      <div :class="$style.section_title">支払金額</div>
      <div :class="$style.section_item">
        {{ detail.core.current_detail.amount }}円
      </div>
    </div>

    <div :class="$style.section">
      <div :class="$style.section_title">払い戻し対象者</div>
      <div :class="$style.target_container">
        <div
          v-for="user in detail.core.repayment_logs"
          :key="user.repaid_to_user.trap_id"
        >
          <Icon :user="user.repaid_to_user.trap_id" :size="24" />
          {{ user.repaid_to_user.trap_id }}
        </div>
      </div>
    </div>

    <div :class="$style.section">
      <div :class="$style.section_title">
        {{ returnRemarksTitle(detail.core.current_detail.type) }}
      </div>
      <!-- eslint-disable-next-line vue/no-v-html -->
      <div :class="$style.section_item" v-html="rendered" />
    </div>

    <!-- 最後のsectionなのでスタイルを当てなくてもOK -->
    <div>
      <div :class="$style.section_title">画像</div>
      <div :class="$style.image_container">
        <img
          v-for="path in detail.core.images"
          :key="path"
          :src="`/api/images/${path}`"
        />
      </div>
      <div v-if="detail.core.images.length === 0">画像はありません</div>
    </div>
  </div>
</template>

<script>
import { applicationType, remarksTitle } from "@/use/applicationDetail";
import { dayPrint } from "@/use/dataFormat";
import { render } from "@/use/markdown";
import StateButtonController from "@/views/components/StateButtonController.vue";
import Icon from "@/views/shared/Icon.vue";
import StateChip from "@/views/shared/StateChip.vue";
import { mapState } from "vuex";

export default {
  components: {
    Icon,
    StateChip,
    StateButtonController
  },
  data: function () {
    return {
      rendered: ""
    };
  },
  computed: {
    ...mapState({ detail: "application_detail_paper" })
  },
  async mounted() {
    this.rendered = await render(this.detail.core.current_detail.remarks);
  },
  methods: {
    returnDate: function (date) {
      return dayPrint(date);
    },
    returnType: function (type) {
      return applicationType(type);
    },
    returnRemarksTitle: function (type) {
      return remarksTitle(type);
    }
  }
};
</script>

<style lang="scss" module>
.container {
  height: fit-content;
  margin: 12px;
  padding: 12px;
  box-shadow: 0 3px 1px -2px rgb(0 0 0 / 20%), 0 2px 2px 0 rgb(0 0 0 / 14%),
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
