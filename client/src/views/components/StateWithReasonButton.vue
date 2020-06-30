<template>
  <v-row justify="center">
    <v-dialog v-model="dialog" max-width="600px">
      <template v-slot:activator="{ on }">
        <v-btn class="primary_accent--text" color="secondary" dark v-on="on">{{
          toStateName(to_state)
        }}</v-btn>
      </template>

      <v-card>
        <v-card-title>
          <span v-if="to_state === `submitted`" class="headline"
            >承認済み→{{ this.toStateName(to_state) }} へ戻す理由</span
          >
          <span v-else class="headline"
            >承認待ち→{{ this.toStateName(to_state) }} への変更理由</span
          >
        </v-card-title>
        <v-form ref="form" v-model="valid">
          <v-card-text>
            <v-container>
              <v-row>
                <v-col cols="12">
                  <v-text-field
                    @blur="blur"
                    ref="reason"
                    :autofocus="dialog"
                    v-model="reason"
                    :rules="nullRules"
                  ></v-text-field>
                </v-col>
              </v-row>
            </v-container>
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="primary" text @click="dialog = false">戻る</v-btn>
            <v-btn :disabled="!valid" color="primary" text @click="postreason"
              >{{ this.toStateName(to_state) }}にする</v-btn
            >
          </v-card-actions>
        </v-form>
      </v-card>
    </v-dialog>
  </v-row>
</template>
<script>
import axios from "axios";
import Vue from "vue";
import { mapActions } from "vuex";
export default {
  data: () => ({
    valid: true,
    dialog: false,
    reason: "",
    nullRules: [v => !!v || ""]
  }),
  props: {
    to_state: String
  },
  watch: {
    dialog: function() {
      if (this.dialog) {
        let self = this;
        Vue.nextTick().then(function() {
          self.$refs.reason.focus();
        });
      }
    }
  },
  methods: {
    ...mapActions(["getApplicationDetail"]),
    blur() {
      if (this.reason === "" || this.reason === undefined) {
        this.$refs.form.reset();
      }
    },
    async postreason() {
      if (this.$refs.form.validate()) {
        await axios
          .put(
            "../api/applications/" +
              this.$store.state.application_detail_paper.core.application_id +
              "/states",
            {
              to_state: this.to_state,
              reason: this.reason
            }
          )
          .then(response => console.log(response.status));
        this.$refs.form.reset();
        this.dialog = false;
        this.getApplicationDetail(
          this.$store.state.application_detail_paper.core.application_id
        );
      }
    },
    toStateName: function(to_state) {
      switch (to_state) {
        case "submitted":
          return "提出済み";
        case "fix_required":
          return "要修正";
        case "rejected":
          return "取り下げ";
        default:
          return "状態が間違っています";
      }
    }
  }
};
</script>
