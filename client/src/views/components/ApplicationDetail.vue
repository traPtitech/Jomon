<!-- 受け取ったデータを基に申請書詳細ページの下半分にログ、コメント等配置 -->
<template>
  <div class="app-detail-paper">
    <v-card class="ml-2 mr-2 mt-2 pa-3" tile>
      <v-row class="ml-4 mr-4" :justify="`space-between`">
        <h1>{{ returnType(this.detail.core.current_detail.type) }}申請書</h1>
        <div>
          <div>申請書ID: {{ this.detail.core.application_id }}</div>
          <v-divider></v-divider>
        </div>
      </v-row>

      <template>
        <v-divider></v-divider>
      </template>
      <h1>タイトル:{{ this.detail.core.current_detail.title }}</h1>

      <div>
        <v-container class="pa-0">
          <v-row>
            <!-- 以下は左列 -->
            <v-col cols="12" md="6">
              <v-row no-gutters>
                <v-col cols="4" md="6">
                  <v-card height="100%" class="pa-2" outlined tile>
                    申請者trapid
                  </v-card>
                </v-col>
                <v-col cols="8" md="6">
                  <v-card height="100%" class="pa-2" outlined tile>
                    <Icon
                      :user="this.detail.core.applicant.trap_id"
                      :size="20"
                    />{{ this.detail.core.applicant.trap_id }}
                  </v-card>
                </v-col>
              </v-row>
              <v-row no-gutters>
                <v-col cols="4" md="6">
                  <v-card height="100%" class="pa-2" outlined tile>
                    申請金額
                  </v-card>
                </v-col>
                <v-col>
                  <v-card height="100%" class="pa-2" outlined tile>
                    {{ this.detail.core.current_detail.amount }}円
                  </v-card>
                </v-col>
              </v-row>
            </v-col>
            <!-- 以上左列以下右列 -->
            <v-col cols="12" md="6">
              <v-row no-gutters>
                <v-col cols="4" md="6">
                  <v-card height="100%" class="pa-2" outlined tile>
                    申請書作成日
                  </v-card>
                </v-col>
                <v-col height="100%" cols="8" md="6">
                  <v-card height="100%" class="pa-2" outlined tile>
                    {{ returnDate(this.detail.core.created_at) }}
                  </v-card>
                </v-col>
              </v-row>
              <v-row no-gutters>
                <v-col cols="4" md="6">
                  <v-card height="100%" class="pa-2" outlined tile>
                    支払った日
                  </v-card>
                </v-col>
                <v-col cols="8" md="6">
                  <v-card height="100%" class="pa-2" outlined tile>
                    {{ returnDate(this.detail.core.current_detail.paid_at) }}
                  </v-card>
                </v-col>
              </v-row>
            </v-col>
            <!-- 以上右列 -->
          </v-row>
        </v-container>
      </div>

      <h3>
        {{ returnRemarkTitle(this.detail.core.current_detail.type) }}:{{
          this.detail.core.current_detail.remarks
        }}
      </h3>
      <h3>払い戻し対象者</h3>
      <!-- 以下のkeyは多分よろしくない -->
      <li
        :key="user.repaid_to_user.trap_id"
        v-for="user in this.detail.core.repayment_logs"
      >
        <Icon :user="user.repaid_to_user.trap_id" :size="25" />
        {{ user.repaid_to_user.trap_id }}
      </li>
      <h3>申請書画像リスト</h3>
      <div :key="path.ID" v-for="path in this.detail.core.images">
        <v-img :src="'/api/images/' + path.ID" max-width="80%" />
      </div>
    </v-card>
  </div>
</template>

<script>
import Icon from "./Icon";
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
          return "購入したものの概要";
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
