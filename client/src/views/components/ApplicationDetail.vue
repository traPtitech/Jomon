<!-- 受け取ったデータを基に申請詳細ページの下半分にログ、コメント等配置 -->
<template>
  <div :class="$style.container">
    <div :class="$style.title_container">
      <div :class="$style.header">
        <div :class="$style.title">
          <h1>{{ returnType(this.detail.core.current_detail.type) }}申請</h1>
          <state-chip :state="this.detail.core.current_state" />
        </div>
        <state-button-controller />
      </div>
    </div>

    <div>
      <div>申請日: {{ returnDate(this.detail.core.created_at) }}</div>
      <v-divider></v-divider>
      <div>
        申請者:
        <Icon :user="this.detail.core.applicant.trap_id" :size="24" />
        {{ this.detail.core.applicant.trap_id }}
      </div>
      <div>
        <v-divider></v-divider>
      </div>
    </div>

    <div>
      <div class="grey--text">概要</div>
      <div class="headline">
        {{ this.detail.core.current_detail.title }}
      </div>
      <v-divider></v-divider>
    </div>

    <div>
      <div class="grey--text">支払日</div>
      <v-row>
        <v-col cols="12" sm="5" class="pt-0 pb-0">
          <div class="headline">
            {{ returnDate(this.detail.core.current_detail.paid_at) }}
          </div>
          <v-divider></v-divider>
        </v-col>
      </v-row>
    </div>

    <div>
      <div class="grey--text">支払金額</div>
      <v-row>
        <v-col cols="12" sm="5" class="pt-0 pb-0">
          <div class="headline">
            {{ this.detail.core.current_detail.amount }}円
          </div>
          <v-divider></v-divider>
        </v-col>
      </v-row>
    </div>

    <div>
      <div class="grey--text">払い戻し対象者</div>
      <div :class="$style.target_container">
        <div
          :key="user.repaid_to_user.trap_id"
          v-for="user in this.detail.core.repayment_logs"
        >
          <Icon :user="user.repaid_to_user.trap_id" :size="25" />
          {{ user.repaid_to_user.trap_id }}
        </div>
      </div>
      <v-divider></v-divider>
    </div>

    <div>
      <div class="grey--text">
        {{ returnRemarksTitle(this.detail.core.current_detail.type) }}
      </div>
      <div class="headline">
        <p
          style="white-space: pre-wrap"
          v-text="this.detail.core.current_detail.remarks"
        ></p>
      </div>
      <v-divider></v-divider>
    </div>

    <div>
      <div class="grey--text">画像</div>
      <div :class="$style.image_container">
        <img
          :key="path"
          v-for="path in this.detail.core.images"
          :src="`/api/images/${path}`"
        />
      </div>
      <div v-if="this.detail.core.images.length === 0">画像はありません</div>
    </div>
  </div>
</template>

<script>
import Icon from "@/views/shared/Icon";
import StateChip from "@/views/shared/StateChip";
import StateButtonController from "@/views/components/StateButtonController";
import { mapState } from "vuex";
import { remarksTitle, applicationType } from "@/use/applicationDetail";
import { dayPrint } from "@/use/dataFormat";

export default {
  components: {
    Icon,
    StateChip,
    StateButtonController
  },
  computed: {
    ...mapState({ detail: "application_detail_paper" })
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
  padding: 8px;
  border: 1px solid #cccccc;
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
.target_container {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(min(100%, 240px), 1fr));
  gap: 16px;
}
.image_container {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(min(100%, 360px), 1fr));
  gap: 16px;
}
</style>
