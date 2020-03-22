<template>
  <v-row justify="center">
    <v-dialog v-model="dialog" max-width="600px">
      <template v-slot:activator="{ on }">
        <v-btn color="primary" dark v-on="on">{{
          toStateName(to_state)
        }}</v-btn>
      </template>

      <v-card>
        <v-card-title>
          <span class="headline"
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
                    label="変更理由"
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
            <v-btn color="blue darken-1" text @click="dialog = false"
              >戻る</v-btn
            >
            <v-btn
              :disabled="!valid"
              color="blue darken-1"
              text
              @click="postreason"
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
    blur() {
      if (this.reason === "" || this.reason === undefined) {
        this.$refs.form.reset();
      }
    },
    postreason() {
      if (this.$refs.form.validate()) {
        axios
          .post("api/applications​/{applicationId}​/states", {
            to_state: this.props.to_state,
            reason: this.reason
          })
          .then(response => console.log(response.status));
        this.$refs.form.reset();
        this.dialog = false;
      }
    },
    toStateName: function(to_state) {
      switch (to_state) {
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
