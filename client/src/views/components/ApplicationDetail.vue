<!-- 受け取ったデータを基に申請書詳細ページの下半分にログ、コメント等配置 -->
<template>
  <div class="app-detail-paper">
    <!-- <h1>ここから</h1> -->
    <v-card>
      <v-row class="ml-4 mr-4" :justify="`space-between`">
        <h1>{{ returnType(this.detail.type) }}申請書</h1>
        <div>
          <div>申請書ID: {{ this.detail.application_id }}</div>
          <v-divider></v-divider>
        </div>
      </v-row>

      <template>
        <v-divider></v-divider>
      </template>
      <h1>タイトル:{{ this.detail.title }}</h1>

      <div>
        <v-container>
          <v-row :justify="`space-around`">
            <v-row class="ml-4 mr-4" no-gutters>
              <v-col>
                <v-card class="pa-2" outlined tile>
                  申請者trapid
                </v-card>
                <v-card class="pa-2" outlined tile>
                  申請金額
                </v-card>
              </v-col>
              <v-col>
                <v-card class="pa-2" outlined tile>
                  <Icon :user="this.detail.applicant.trap_id" :size="20" />{{
                    this.detail.applicant.trap_id
                  }}
                </v-card>
                <v-card class="pa-2" outlined tile>
                  {{ this.detail.ammount }}円
                </v-card>
              </v-col>
            </v-row>
            <v-row class="ml-4 mr-4" no-gutters>
              <v-col>
                <v-card class="pa-2" outlined tile>
                  申請書作成日
                </v-card>
                <v-card class="pa-2" outlined tile>
                  支払った日
                </v-card>
              </v-col>
              <v-col>
                <v-card class="pa-2" outlined tile>
                  {{ returnDate(this.detail.created_at) }}
                </v-card>
                <v-card class="pa-2" outlined tile>
                  {{ returnDate(this.detail.paid_at) }}
                </v-card>
              </v-col>
            </v-row>
          </v-row>
        </v-container>
      </div>

      <h3>
        {{ returnRemarkTitle(this.detail.type) }}:{{ this.detail.remarks }}
      </h3>
      <h3>払い戻し対象者</h3>
      <li :key="user" v-for="user in this.detail.repaid_to_id">
        <Icon :user="user" :size="25" />
        {{ user }}
      </li>
      <h3>申請書画像リスト</h3>
      <li :key="path" v-for="path in this.detail.images">
        {{ path }}
      </li>
    </v-card>
    <!-- <h1>ここまで</h1> -->
  </div>
</template>

<script>
import Icon from "./Icon";
export default {
  data: function() {
    return {
      detail: this.$store.state.application_detail_paper
    };
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
  // name: "ApplicationDetail"
  //.Vueファイルではnameはファイル名なので不要
  //.Vueファイルについてはhttps://jp.vuejs.org/v2/guide/single-file-components.html

  //props: ['test',detail]
  //propsはhttps://jp.vuejs.org/v2/style-guide/index.htmlによると下の書き方ほうが良い。
  //型が異なるとjavascriptのコンソール画面で警告が出るhttps://jp.vuejs.org/v2/guide/components-props.html
  props: {
    //申請書詳細はpropsで管理しない
    // test: String,
    // detail: Object
  },
  components: {
    Icon
  }
};
</script>
