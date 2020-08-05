<template>
  <div>
    <v-row>
      <v-col>
        <with-reason-button
          v-if="
            repaidtotrapid.length ===
            this.$store.state.application_detail_paper.core.repayment_logs
              .length
          "
          class="mr-4"
          to_state="submitted"
        />
      </v-col>
      <v-col>
        <v-dialog v-model="dialog" scrollable max-width="500px">
          <template v-slot:activator="{ on }">
            <v-btn color="primary" dark v-on="on">払い戻し完了者の選択</v-btn>
          </template>
          <v-card>
            <v-menu
              v-model="menu"
              :close-on-content-click="false"
              :nudge-right="40"
              lazy
              transition="scale-transition"
              offset-y
              full-width
              min-width="290px"
            >
              <template v-slot:activator="{ on }">
                <v-text-field
                  v-model="date"
                  label="払い戻し日"
                  readonly
                  v-on="on"
                ></v-text-field>
              </template>
              <v-date-picker
                v-model="date"
                @input="menu = false"
              ></v-date-picker>
            </v-menu>
            <v-autocomplete
              ref="traPID"
              v-model="traPID"
              :rules="[() => !!traPID || '払い戻し完了者は一人以上必要です']"
              :items="repaidtotrapid"
              label="払い戻し完了者のtraPidを入力..."
              required
              multiple
            >
            </v-autocomplete>
            <v-btn color="primary" @click="putRepaid(traPID, date)">OK</v-btn>
          </v-card>
        </v-dialog>
      </v-col>
    </v-row>
    <span v-if="repaidtotrapid.length === 0">
      何かがおかしいです。一度リロードしなおしてみて下さい。
    </span>
  </div>
</template>
<script>
import axios from "axios";
import WithReasonButton from "./StateWithReasonButton";
import { mapActions } from "vuex";

export default {
  data: () => ({
    date: new Date().toISOString().substr(0, 10),
    menu: false,
    dialog: false,
    traPID: []
  }),
  components: {
    WithReasonButton
  },
  methods: {
    ...mapActions(["getApplicationDetail"]),
    async putRepaid(traPIDs, date) {
      await Promise.all(
        traPIDs.map(async traPID => {
          await axios
            .put(
              "../api/applications/" +
                this.$store.state.application_detail_paper.core.application_id +
                "/states/repaid/" +
                traPID,
              {
                repaid_at: date
              }
            )
            .then(response => console.log(traPID, response.status));
        })
      ).then(() => {
        this.traPID = [];
        this.dialog = false;
        this.getApplicationDetail(
          this.$store.state.application_detail_paper.core.application_id
        );
      });
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
