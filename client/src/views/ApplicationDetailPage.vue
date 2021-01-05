<template>
  <div v-if="this.loading">loading...</div>
  <div v-else :class="$style.container">
    <!-- todo storeのfixで制御する、このページのcreatedでstoreのfixはfalseに。 -->
    <application-paper :class="$style.paper" v-if="!this.detail.fix" />
    <fix-application-paper :class="$style.paper" v-else />
    <application-logs />
  </div>
</template>

<script>
import ApplicationPaper from "@/views/components/ApplicationDetail";
import FixApplicationPaper from "@/views/components/FixApplicationDetail";
import ApplicationLogs from "@/views/components/ApplicationDetailLogs";
import { mapState, mapActions, mapMutations } from "vuex";

export default {
  components: {
    ApplicationPaper,
    FixApplicationPaper,
    ApplicationLogs
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

<style lang="scss" module>
.container {
  display: flex;
  justify-content: center;
  flex-wrap: wrap;
  max-width: 100%;
}
.paper {
  max-width: 60vw;
}
@media (max-width: $breakpoint) {
  .paper {
    max-width: 100%;
  }
}
</style>
