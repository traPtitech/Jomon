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
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="green darken-1" @click="login">OK</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-layout>
</template>

<script>
import { mapMutations, mapState } from "vuex";
import { redirectAuthorizationEndpoint } from "./../../utils/api";
import { createCodeChallenge } from "./../../utils/hash";
export default {
  computed: {
    ...mapState(["token"])
  },
  methods: {
    ...mapMutations(["toggleLoginDialog"]),
    async login() {
      //   console.log(`in login:` + this.token.loginDialog);
      this.toggleLoginDialog();
      //   console.log(`in login:` + this.token.loginDialog);
      const client_id = "client_id";
      const response_type = "code";
      const code_challenge = await createCodeChallenge();
      const code_challenge_method = "S256";
      redirectAuthorizationEndpoint(
        client_id,
        response_type,
        code_challenge,
        code_challenge_method
      );
    }
  }
};
</script>
