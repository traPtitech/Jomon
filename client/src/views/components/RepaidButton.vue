<template>
  <div>
    <v-dialog v-model="dialog" scrollable max-width="500px">
      <template v-slot:activator="{ on }">
        <simple-button :label="'払い戻し済みのユーザーを選択'" v-on="on" />
      </template>
      <v-card :class="$style.modal">
        <v-menu
          :class="$style.user_select"
          v-model="menu"
          :close-on-content-click="false"
          :nudge-right="40"
          transition="scale-transition"
          offset-y
        >
          <template v-slot:activator="{ on }">
            <v-text-field
              v-model="date"
              label="払い戻し完了日"
              readonly
              v-on="on"
            />
          </template>
          <v-date-picker v-model="date" @input="menu = false" />
        </v-menu>
        <v-autocomplete
          ref="traPID"
          v-model="traPID"
          :rules="[
            traPID =>
              traPID.length > 0 ||
              '払い戻し済みのユーザーが一人以上選ばれている必要があります'
          ]"
          :items="repaidToTraPId"
          label="払い戻し済みのユーザーを選択"
          required
          multiple
        />
        <simple-button
          :label="'OK'"
          :disabled="traPID.length === 0"
          :variant="'secondary'"
          @click="putRepaid(traPID, date)"
        />
      </v-card>
    </v-dialog>
    <span v-if="repaidToTraPId.length === 0">
      何かがおかしいです。一度リロードしなおしてみて下さい。
    </span>
  </div>
</template>
<script>
import axios from "axios";
import { mapActions } from "vuex";
import SimpleButton from "@/views/shared/SimpleButton";

export default {
  components: {
    SimpleButton
  },
  data: () => ({
    date: new Date().toISOString().substr(0, 10),
    menu: false,
    dialog: false,
    traPID: []
  }),
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
    repaidToTraPId() {
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

<style lang="scss" module>
.modal {
  padding: 8px;
}
.user_select {
  min-width: 280px;
}
</style>
