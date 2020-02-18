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
        <v-card-text>
          <v-container>
            <v-row>
              <v-col cols="12">
                <v-text-field label="変更理由*" required></v-text-field>
              </v-col>
            </v-row>
          </v-container>
          <small>*indicates required field</small>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="blue darken-1" text @click="dialog = false">戻る</v-btn>
          <v-btn color="blue darken-1" text @click="dialog = false"
            >{{ this.toStateName(to_state) }}にする</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-row>
</template>
<script>
export default {
  data: () => ({
    dialog: false
  }),
  props: {
    to_state: String
  },
  methods: {
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
