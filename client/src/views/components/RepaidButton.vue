<template>
  <div>
    <with-reason-button
      v-if="
        repaidtotrapid.length ===
          this.$store.state.application_detail_paper.core.repayment_logs.length
      "
      class="mr-4"
      to_state="submitted"
    />
    払い戻し完了ボタン:
    <span v-for="user in repaidtotrapid" :key="user">
      <v-btn color="primary" dark v-on:click="putRepaid(user)">{{
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
import WithReasonButton from "./StateWithReasonButton";
export default {
  components: {
    WithReasonButton
  },
  methods: {
    putRepaid(repaid_to_trap_id) {
      axios
        .put(
          "../api/applications/" +
            this.$store.state.application_detail_paper.core.application_id +
            "/states/repaid/" +
            repaid_to_trap_id
        )
        .then(response => console.log(response.status));
      alert(repaid_to_trap_id + "に払い戻ししました。");
    }
  },
  computed: {
    repaidtotrapid() {
      let trap_ids = [];
      this.$store.state.application_detail_paper.core.repayment_logs.forEach(
        log => {
          if (log.repaid_at === "" || log.repaid_at === null) {
            trap_ids.push(log.repaid_to_user.trap_id);
          }
        }
      );
      return trap_ids;
    }
  }
};
</script>
