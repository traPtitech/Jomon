<template>
  <!-- 本番はこのファイルの
  this.adminはthis.$store.state.me.adminへ、
this.applicantはthis.detail.applicant.trap_idとすれば良い。この際は、dataを削除する-->
  <div v-if="this.detail.current_state === `submitted` && this.admin">
    <v-row>
      <v-btn v-on:click="accept()">承認 </v-btn>
      <with-reason-button class="ml-4 mr-5" to_state="fix_required" />
      <with-reason-button class="mr-4" to_state="rejected" />
    </v-row>
  </div>
  <div v-else-if="this.detail.current_state === `accepted`">
    払い戻し完了ボタン:
    <repaid-button />
    <!-- このリストの制御とリストのクリック先 -->
    <!-- 条件付きでsubmittedへ -->
  </div>
  <div
    v-else-if="
      this.detail.current_state === `fix_required` &&
        (this.admin || this.applicant)
    "
  >
    <v-btn>修正</v-btn>
    <!-- この修正ボタンの先 -->
  </div>
</template>
<script>
import axios from "axios";
import WithReasonButton from "./StateWithReasonButton";
import RepaidButton from "./RepaidButton";
import { mapState } from "vuex";
export default {
  data: function() {
    return {
      admin: true,
      applicant: true
    };
  },
  components: {
    WithReasonButton,
    RepaidButton
  },
  computed: {
    ...mapState({ detail: "application_detail_paper" })
  },
  methods: {
    accept() {
      axios
        .put(
          "../api/applications/" +
            this.$store.state.application_detail_paper.application_id +
            "/states",
          {
            to_state: "accepted"
          }
        )
        .then(response => console.log(response.status));
      alert("承認しました");
    }
  }
};
</script>
