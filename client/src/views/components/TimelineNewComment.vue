<template>
  <div>
    <v-row justify="space-between">
      <v-col cols="1">
        <Icon :user="trapId" :size="25" />
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
<script lang="ts">
import { useMeStore } from "@/stores/me";
import Icon from "@/views/shared/Icon.vue";
import axios from "axios";
import { mapState } from "pinia";
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
  computed: {
    ...mapState(useMeStore, ["trapId"])
  },
  methods: {
    blur() {
      if (this.comment === "" || this.comment === undefined) {
        (this.$refs.form as any).reset();
      }
    },
    ...mapActions(["getApplicationDetail"]),
    postcomment() {
      if ((this.$refs.form as any).validate()) {
        axios
          .post(
            "/api/applications/" +
              this.$store.state.application_detail_paper.core.application_id +
              "/comments",
            {
              comment: this.comment
            }
          )
          .catch((e: any) => {
            alert(e);
            return;
          });
        (this.$refs.form as any).reset();
        this.getApplicationDetail(
          this.$store.state.application_detail_paper.core.application_id
        );
      }
    }
  }
};
</script>
