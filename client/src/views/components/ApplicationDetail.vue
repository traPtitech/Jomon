<!-- 受け取ったデータを基に申請書詳細ページの下半分にログ、コメント等配置 -->
<template>
  <div class="app-detail-paper">
    <h1>ここから</h1>
    <h1>{{ this.$store.state.application_detail_test }}</h1>
    <h1>{{ test }}</h1>
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
    <h3>{{ detail.images[0] }}</h3>
    <h1>ここまで</h1>
  </div>
</template>

<script>
export default {
  // name: "ApplicationDetail"
  //.Vueファイルではnameはファイル名なので不要

  //props: ['test',detail]
  //propsはhttps://jp.vuejs.org/v2/style-guide/index.htmlによると下の書き方ほうが良い。

  props: {
    test: String,
    detail: Object
  }
};
</script>
