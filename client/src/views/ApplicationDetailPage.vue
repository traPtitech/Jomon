<!-- 申請ページ描画画面制御 -->
<!-- コンポーネント配置とaxiosで受け取ったjsonを分割し各コンポーネントに渡す-->
<!-- todo 返金の操作ボタン -->
<template>
  <div v-if="this.loading">loading</div>
  <div v-else class="application-detail">
    <!-- todo storeのfixで制御する、このページのcreatedでstoreのfixはfalseに。 -->
    <div>
      <application-paper v-if="!this.detail.fix" />
      <fix-application-paper v-else />
      <state-button-controller />
    </div>
    <application-logs />
  </div>
</template>

<script>
import ApplicationPaper from "./components/ApplicationDetail";
import FixApplicationPaper from "./components/FixApplicationDetail";
import ApplicationLogs from "./components/ApplicationDetailLogs";
import StateButtonController from "./components/StateButtonController";
import { mapState, mapActions, mapMutations } from "vuex";

export default {
  components: {
    ApplicationPaper,
    FixApplicationPaper,
    ApplicationLogs,
    StateButtonController
  },
  data() {
    return {
      loading: true
    };
  },
  computed: {
    ...mapState({ detail: "application_detail_paper" })
  },
  async created() {
    await this.getApplicationDetail(this.$route.params.id);
    this.deleteFix();
    this.loading = false;
  },
  methods: {
    ...mapMutations(["deleteFix"]),
    ...mapActions(["getApplicationDetail"])
  }
};
</script>
