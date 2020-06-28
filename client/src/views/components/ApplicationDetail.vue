<!-- 受け取ったデータを基に申請書詳細ページの下半分にログ、コメント等配置 -->
<template>
  <v-container>
    <v-card class="ml-2 mr-2 mt-2 pa-3" tile>
      <v-row class="ml-4 mr-4" :justify="`space-between`">
        <v-col cols="12" sm="8" class="pt-0 pb-0">
          <h1>{{ returnType(this.detail.core.current_detail.type) }}申請</h1>
        </v-col>

        <v-col cols="12" sm="4" class="pt-0 pb-0">
          <div>申請日: {{ returnDate(this.detail.core.created_at) }}</div>
          <v-divider></v-divider>
          <div>
            申請者:<Icon
              :user="this.detail.core.applicant.trap_id"
              :size="20"
            />{{ this.detail.core.applicant.trap_id }}
          </div>
          <div>
            <v-divider></v-divider>
          </div>
        </v-col>
      </v-row>

      <template>
        <v-divider class="mt-1"></v-divider>
      </template>

      <div>
        <div class="grey--text">
          概要
        </div>
        <div class="headline">
          {{ this.detail.core.current_detail.title }}
        </div>
        <v-divider></v-divider>
      </div>

      <div>
        <div class="grey--text">
          支払日
        </div>
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
        <div class="grey--text">
          支払金額
        </div>
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
        <div class="grey--text">
          払い戻し対象者
        </div>
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
          {{ returnRemarkTitle(this.detail.core.current_detail.type) }}
        </div>
        <div class="headline">
          <p
            style="white-space:pre-wrap;"
            v-text="this.detail.core.current_detail.remarks"
          ></p>
        </div>
        <v-divider></v-divider>
      </div>

      <div>
        <div class="grey--text">
          画像
        </div>
        <div :key="path" v-for="path in this.detail.core.images">
          <v-img :src="'/api/images/' + path" max-width="80%" />
        </div>
        <div v-if="this.detail.core.images.length == 0">
          画像はありません
        </div>
      </div>
    </v-card>
  </v-container>
</template>

<script>
import Icon from "../shered/Icon";
import { mapState } from "vuex";
export default {
  computed: {
    ...mapState({ detail: "application_detail_paper" })
  },
  methods: {
    returnDate: function(date) {
      const normalizedDate = new Date(date);
      return (
        normalizedDate.getFullYear() +
        "年" +
        (normalizedDate.getMonth() + 1) +
        "月" +
        normalizedDate.getDate() +
        "日"
      );
    },
    returnType: function(type) {
      switch (type) {
        case "club":
          return "部費利用";
        case "contest":
          return "大会等旅費補助";
        case "event":
          return "イベント交通費補助";
        case "public":
          return "渉外交通費補助";
        default:
          return "タイプが間違っています";
      }
    },
    returnRemarkTitle: function(type) {
      switch (type) {
        case "club":
          return "購入物の詳細";
        case "contest":
          return "旅程";
        case "event":
          return "乗車区間";
        case "public":
          return "乗車区間";
        default:
          return "タイプが間違っています";
      }
    }
  },
  props: {},
  components: {
    Icon
  }
};
</script>
