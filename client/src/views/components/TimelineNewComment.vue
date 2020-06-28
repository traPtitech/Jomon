<template>
  <div>
    <v-row justify="space-between">
      <v-col cols="1">
        <Icon :user="this.$store.state.me.trap_id" :size="25" />
      </v-col>
      <v-col cols="11">
        <v-card class="pa-2">
          <v-form ref="form" v-model="valid">
            <v-textarea
              @blur="blur"
              v-model="comment"
              :rules="nullRules"
              outlined
              label="コメントを書いてください"
            ></v-textarea>
            <v-btn
              :disabled="!valid"
              color="primary"
              class="mr-4"
              @click="postcomment"
            >
              コメントする
            </v-btn>
          </v-form>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>
<script>
import axios from "axios";
import Icon from "../shered/Icon";
export default {
  data: () => {
    return {
      valid: true,
      comment: "",
      nullRules: [v => !!v || ""]
    };
  },
  components: {
    Icon
  },
  methods: {
    blur() {
      if (this.comment === "" || this.comment === undefined) {
        this.$refs.form.reset();
      }
    },
    ...mapActions(["getApplicationDetail"]),
    postcomment() {
      if (this.$refs.form.validate()) {
        axios
          .post(
            "/api/applications/" +
              this.$store.state.application_detail_paper.core.application_id +
              "/comments",
            {
              comment: this.comment
            }
          )
          .then(response => console.log(response.status));
        this.$refs.form.reset();
        this.getApplicationDetail(
          this.$store.state.application_detail_paper.core.application_id
        );
      }
    }
  }
};
</script>
