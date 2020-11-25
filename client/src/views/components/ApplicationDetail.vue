<!-- 受け取ったデータを基に申請詳細ページの下半分にログ、コメント等配置 -->
<template>
  <v-card :class="$style.card" outlined tile>
    <v-row :class="$style.title_container">
      <div :class="$style.title">
        <h1>{{ returnType(this.detail.core.current_detail.type) }}申請</h1>
        <state-chip :state="this.detail.core.current_state" />
      </div>
      <v-col>
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
      </v-col>
    </v-row>

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
      <v-row>
        <v-col cols="12" sm="5" class="pt-0 pb-0">
          <div
            :key="user.repaid_to_user.trap_id"
            v-for="user in this.detail.core.repayment_logs"
          >
            <Icon :user="user.repaid_to_user.trap_id" :size="25" />
            {{ user.repaid_to_user.trap_id }}
          </div>
          <v-divider></v-divider>
        </v-col>
      </v-row>
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
      <div
        :class="$style.image"
        :key="path"
        v-for="path in this.detail.core.images"
      >
        <v-img :src="`/api/images/${path}`" />
      </div>
      <div v-if="this.detail.core.images.length === 0">画像はありません</div>
    </div>
  </v-card>
</template>

<script>
import Icon from "@/views/shared/Icon";
import StateChip from "@/views/shared/StateChip";
import { mapState } from "vuex";
import { remarksTitle, applicationType } from "@/use/applicationDetail";
import { dayPrint } from "@/use/dataFormat";

export default {
  components: {
    Icon,
    StateChip
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
.card {
  margin: 24px;
  padding: 8px;
  max-width: 1200px;
}
.title_container {
  margin: 8px;
  justify-content: space-between;
}
.title {
  display: flex;
  align-items: center;
}
.image {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
}
</style>
