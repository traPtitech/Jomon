<!-- 申請ページ描画画面制御 -->
<!-- コンポーネント配置とaxiosで受け取ったjsonを分割し各コンポーネントに渡す-->
<!-- todo 返金の操作ボタン -->
<template>
  <div v-if="this.loading">loading</div>
  <div v-else class="application-detail">
    <div :class="$style.banner_container">
      <div :class="$style.current_state">
        <h3>現在の状態:{{ returnState(this.detail.core.current_state) }}</h3>
      </div>
      <state-button-controller />
    </div>
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
  data() {
    return {
      loading: true
    };
  },
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
  async created() {
    await this.getApplicationDetail(this.$route.params.id);
    this.deleteFix();
    this.loading = false;
  },
  methods: {
    ...mapMutations(["deleteFix"]),
    ...mapActions(["getApplicationDetail"]),
    returnState: function (state) {
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

<style lang="scss" module>
.current_state {
  // TODO: 文字の高さを中央揃えに
  margin: 4px;
}
.banner_container {
  display: flex;
  justify-content: space-between;
  width: 100%;
}
</style>
