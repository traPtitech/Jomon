<template>
  <v-layout row justify-center>
    <v-dialog v-model="token.loginDialog" persistent max-width="290">
      <v-card>
        <v-card-title class="headline">
          ログインしてください
        </v-card-title>
        <v-card-text>
          OKを押すとtraQに飛びます。<br />
          「承認」を押すとログインできます。
        </v-card-text>
        <v-card-actions
          >q
          <v-spacer></v-spacer>
          <v-btn color="green darken-1" @click="login">OK</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-layout>
</template>

<script>
import { mapMutations, mapState } from "vuex";
import axios from "axios";
import { redirectAuthorizationEndpoint } from "../../utils/api";

export default {
  computed: {
    ...mapState(["token"])
  },
  methods: {
    ...mapMutations(["toggleLoginDialog"]),
    async login() {
      this.toggleLoginDialog();
      try {
        const response = await axios.get("/api/auth/genpkce");
        const client_id = response.data.client_id;
        const response_type = response.data.response_type;
        const code_challenge = response.data.code_challenge;
        const code_challenge_method = response.data.code_challenge_method;
        redirectAuthorizationEndpoint(
          client_id,
          response_type,
          code_challenge,
          code_challenge_method
        );
      } catch (err) {
        console.log(err);
      }
    }
  }
};
</script>
