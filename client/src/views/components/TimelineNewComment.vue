<template>
  <div>
    <v-row justify="space-between">
      <v-col cols="1">
        <Icon :user="$store.state.me.trap_id" :size="25" />
      </v-col>
      <v-col cols="11">
        <v-card class="pa-2">
          <v-form ref="form" v-model="valid">
            <v-textarea
              v-model="comment"
              :rules="nullRules"
              outlined
              label="コメントを書いてください"
              @blur="blur"
            />
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
import Icon from "@/views/shared/Icon.vue";
import axios from "axios";
import { mapActions } from "vuex";

export default {
  components: {
    Icon
  },
  data: () => {
    return {
      valid: true,
      comment: "",
      nullRules: [v => !!v || ""]
    };
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
          .catch(e => {
            alert(e);
            return;
          });
        this.$refs.form.reset();
        this.getApplicationDetail(
          this.$store.state.application_detail_paper.core.application_id
        );
      }
    }
  }
};
</script>
