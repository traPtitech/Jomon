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
      <v-menu>
        <template v-slot:activator="{ on }">
          <v-btn color="primary" dark @click="repay.name = user" v-on="on">{{
            user
          }}</v-btn>
        </template>
        <v-date-picker
          v-model="repay.date"
          no-title
          @input="dialog = true"
        ></v-date-picker>
      </v-menu>
    </span>
    <span v-if="repaidtotrapid.length === 0">
      何かがおかしいです。一度リロードしなおしてみて下さい。
    </span>
    <v-dialog v-model="dialog">
      <v-card>
        <v-card-title>
          <span>{{ repay.name }}へ払い戻ししますか？</span>
        </v-card-title>
        <br />
        支払日:{{ repay.date }}
        <br />
        <v-btn color="primary" @click="dialog = false">no</v-btn
        ><v-btn color="primary" @click="putRepaid(repay.name)">yes</v-btn>
      </v-card>
    </v-dialog>
  </div>
</template>
<script>
import axios from "axios";
import WithReasonButton from "./StateWithReasonButton";
export default {
  components: {
    WithReasonButton
  },
  data: () => ({
    dialog: false,
    repay: {
      date: "",
      name: ""
    }
  }),
  methods: {
    putRepaid(repaid_to_trap_id) {
      axios
        .put(
          "../api/applications/" +
            this.$store.state.application_detail_paper.core.application_id +
            "/states/repaid/" +
            repaid_to_trap_id,
          { repaid_at: this.getDate() }
        )
        .then(response => console.log(response.status));
      alert(repaid_to_trap_id + "に払い戻ししました。");
    },
    getDate() {
      let date = new Date(this.repay.date);
      let y = date.getFullYear();
      let m = ("00" + (date.getMonth() + 1)).slice(-2);
      let d = ("00" + date.getDate()).slice(-2);
      return y + "-" + m + "-" + d;
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
