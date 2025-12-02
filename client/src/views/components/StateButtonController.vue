<template>
  <div v-if="displayAcceptBottom" :class="$style.button_container">
    <simple-button :label="'承認'" @click="accept()" />
    <with-reason-button to-state="fix_required" />
    <with-reason-button to-state="rejected" />
  </div>
  <div v-else-if="displayRepaidBottom" :class="$style.button_container">
    <with-reason-button to-state="submitted" />
    <repaid-button />
  </div>
  <div v-else-if="displayFixResubmitBottom" :class="$style.button_container">
    <simple-button :label="'修正'" @click="changeFix" />
    <simple-button :label="'再申請'" @click="reSubmit" />
  </div>
</template>
<script>
import SimpleButton from "@/views/shared/SimpleButton.vue";
import axios from "axios";
import { mapActions, mapMutations, mapState } from "vuex";
import RepaidButton from "./RepaidButton.vue";
import WithReasonButton from "./StateWithReasonButton.vue";

export default {
  components: {
    WithReasonButton,
    RepaidButton,
    SimpleButton
  },
  data: function () {
    return {
      dialog: false
    };
  },
  computed: {
    ...mapState({ detail: "application_detail_paper" }),
    displayAcceptBottom() {
      return (
        this.detail.core.current_state === `submitted` &&
        this.$store.state.me.is_admin
      );
    },
    displayRepaidBottom() {
      return (
        this.detail.core.current_state === `accepted` &&
        this.$store.state.me.is_admin
      );
    },
    displayFixResubmitBottom() {
      return (
        this.detail.core.current_state === `fix_required` &&
        (this.$store.state.me.is_admin ||
          this.$store.state.me.trap_id === this.detail.core.applicant.trap_id)
      );
    }
  },
  methods: {
    ...mapMutations(["changeFix"]),
    ...mapActions(["getApplicationDetail"]),
    async accept() {
      await axios
        .put(
          "../api/applications/" +
            this.$store.state.application_detail_paper.core.application_id +
            "/states",
          {
            to_state: "accepted"
          }
        )
        .catch(e => {
          alert(e);
          return;
        });
      alert("承認しました");
      this.getApplicationDetail(
        this.$store.state.application_detail_paper.core.application_id
      );
    },
    async reSubmit() {
      await axios
        .put(
          "../api/applications/" +
            this.$store.state.application_detail_paper.core.application_id +
            "/states",
          {
            to_state: "submitted"
          }
        )
        .catch(e => {
          alert(e);
          return;
        });
      alert("再申請しました");
      this.getApplicationDetail(
        this.$store.state.application_detail_paper.core.application_id
      );
    }
  }
};
</script>

<style lang="scss" module>
.button_container {
  display: flex;
}
</style>
