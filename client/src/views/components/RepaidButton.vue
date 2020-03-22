<template>
  <div>
    <span v-for="user in repaidtotrapid" :key="user">
      <v-btn color="primary" dark v-on:click="postrepaid(user)">{{
        user
      }}</v-btn>
    </span>
    <span v-if="repaidtotrapid.length === 0">
      何かがおかしいです。一度リロードしなおしてみて下さい。
    </span>
  </div>
</template>
<script>
import axios from "axios";
export default {
  methods: {
    postrepaid(repaid_to_trap_id) {
      axios
        .put(
          "../api/applications/" +
            this.$store.state.application_detail_paper.application_id +
            "/states/repaid/" +
            repaid_to_trap_id
        )
        .then(response => console.log(response.status));
      alert(repaid_to_trap_id + "に払い戻ししました。");
    }
  },
  computed: {
    repaidtotrapid() {
      let trap_ids = new Array();
      let repayment_logs = this.$store.state.application_detail_paper
        .repayment_logs;
      let count = 0;
      for (let i = 0; i < repayment_logs.length - 1; i++) {
        if (
          repayment_logs[i].repaid_at === "" ||
          repayment_logs[i].repaid_at === null
        ) {
          trap_ids[count] = repayment_logs[i].repaid_to_user.trap_id;
          count++;
        }
      }
      return trap_ids;
    }
  }
};
</script>
