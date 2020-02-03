<!-- 受け取ったデータを基に申請書詳細ページの下半分にログ、コメント等配置 -->
<template>
  <div class="app-detail-paper">
    <h1>ここから</h1>
    <!-- <h1>{{ this.$store.state.application_detail_test }}</h1> -->
    <v-card>
      <h1>{{ this.$store.state.application_detail_paper.type }}申請書</h1>
      <template>
        <v-divider></v-divider>
      </template>
      <h1>
        申請書タイトル:{{ this.$store.state.application_detail_paper.title }}
      </h1>
      <h3>
        申請書作成日:{{ this.$store.state.application_detail_paper.created_at }}
      </h3>
      <h3>
        申請者trapid:
        <Icon
          :user="this.$store.state.application_detail_paper.applicant.trap_id"
          :size="25"
        />{{ this.$store.state.application_detail_paper.applicant.trap_id }}
      </h3>
      <h3>
        申請者displayname:
        {{ this.$store.state.application_detail_paper.applicant.display_name }}
      </h3>
      <h3>
        現在の状態:{{
          this.$store.state.application_detail_paper.current_state
        }}
      </h3>
      <h3>
        申請書のタイプ:{{ this.$store.state.application_detail_paper.type }}
      </h3>
      <h3>
        申請書備考:{{ this.$store.state.application_detail_paper.remarks }}
      </h3>
      <h3>
        支払った日:{{ this.$store.state.application_detail_paper.paid_at }}
      </h3>
      <h3>
        申請金額:{{ this.$store.state.application_detail_paper.ammount }}円
      </h3>
      <h3>申請者リスト</h3>
      <li
        :key="user"
        v-for="user in this.$store.state.application_detail_paper.repaid_to_id"
      >
        <Icon :user="user" :size="25" />
        {{ user }}
      </li>
      <h3>申請書画像リスト</h3>
      <li
        :key="path"
        v-for="path in this.$store.state.application_detail_paper.images"
      >
        {{ path }}
      </li>
      <!-- <h1>{{ test }}</h1>
    <h1>{{ detail.title }}</h1>
    <h3>{{ detail.created_at }}</h3>
    <h3>{{ detail.applicant.trap_id }}</h3>
    <h3>{{ detail.applicant_display_name }}</h3>
    <h3>{{ detail.current_state }}</h3>
    <h3>{{ detail.type }}</h3>
    <h3>{{ detail.remarks }}</h3>
    <h3>{{ detail.paid_at }}</h3>
    <h3>{{ detail.amount }}</h3>
    <h3>{{ detail.repaid_to_id[0] }}</h3>
    <h3>{{ detail.images[0] }}</h3> -->
    </v-card>
    <h1>ここまで</h1>
  </div>
</template>

<script>
import Icon from "./Icon";
export default {
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
