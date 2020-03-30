<!-- 申請書ページ描画画面制御 -->
<!-- コンポーネント配置とaxiosで受け取ったjsonを分割し各コンポーネントに渡す-->
<!-- todo 返金の操作ボタン -->
<template>
  <div class="application-detail">
    <!-- <div>ID: {{ $route.params.id }}</div> -->
    <v-row :justify="`space-between`">
      <h3 class="ml-4">
        現在の状態:{{ returnState(this.detail.core.current_state) }}
      </h3>
      <state-button-controller class="mr-4 ml-10" />
    </v-row>
    <!-- todo storeのfixで制御する、このページのcreatedでstoreのfixはfalseに。 -->
    <application-paper v-if="!this.detail.fix" />
    <fix-application-paper v-else />
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
  name: "AppDetail",
  components: {
    ApplicationPaper,
    FixApplicationPaper,
    ApplicationLogs,
    StateButtonController
  },
  computed: {
    ...mapState({ detail: "application_detail_paper" })
  },
  created() {
    this.getApplicationDetail(this.$route.params.id);
    this.deleteFix();
  },
  methods: {
    ...mapMutations(["deleteFix"]),
    ...mapActions(["getApplicationDetail"]),
    returnState: function(state) {
      switch (state) {
        case "submitted":
          return "承認待ち";
        case "fix_required":
          return "要修正";
        case "accepted":
          return "承認済み";
        case "fully_repaid":
          return "返済完了";
        case "rejected":
          return "取り下げ";
        default:
          return "状態が間違っています";
      }
    }
  }
};
</script>
