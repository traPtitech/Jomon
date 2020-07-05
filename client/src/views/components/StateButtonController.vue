<template>
  <div v-if="displayAcceptBottom">
    <v-row>
      <v-btn
        class="primary_accent--text"
        color="secondary"
        v-on:click="accept()"
        >承認
      </v-btn>
      <with-reason-button class="ml-4 mr-5" to_state="fix_required" />
      <with-reason-button class="mr-4" to_state="rejected" />
    </v-row>
  </div>
  <div v-else-if="displayRepaidBottom">
    <repaid-button />
  </div>
  <div v-else-if="displayFixResubmitBottom">
    <v-btn
      :disabled="this.detail.fix"
      class="primary_accent--text"
      color="secondary"
      @click="changeFix"
      >修正</v-btn
    >
    <v-btn
      :disabled="this.detail.fix"
      class="primary_accent--text"
      color="secondary"
      @click="reSubmit"
      >再申請</v-btn
    >
  </div>
</template>
<script>
import axios from "axios";
import WithReasonButton from "./StateWithReasonButton";
import RepaidButton from "./RepaidButton";
import { mapState, mapMutations, mapActions } from "vuex";
export default {
  data: function() {
    return {
      dialog: false
    };
  },
  components: {
    WithReasonButton,
    RepaidButton
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
        .then(response => console.log(response.status));
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
        .then(response => console.log(response.status));
      alert("再申請しました");
      this.getApplicationDetail(
        this.$store.state.application_detail_paper.core.application_id
      );
    }
  }
};
</script>
