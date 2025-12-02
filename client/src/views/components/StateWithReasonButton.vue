<template>
  <div :class="$style.container">
    <v-dialog v-model="isDialogOpen" max-width="600px">
      <template #activator="{ props }">
        <simple-button
          :label="
            toStateName(toState) + (toState === 'submitted' ? 'に戻す' : '')
          "
          :variant="toState === 'submitted' ? 'warning' : 'error'"
          v-bind="props"
        />
      </template>

      <v-card>
        <v-card-title>
          <span v-if="toState === `submitted`" class="headline"
            >承認済み→{{ toStateName(toState) }} へ戻す理由</span
          >
          <span v-else class="headline"
            >承認待ち→{{ toStateName(toState) }} への変更理由</span
          >
        </v-card-title>
        <v-form ref="form" v-model="valid">
          <v-card-text>
            <v-container>
              <v-row>
                <v-col cols="12">
                  <v-text-field
                    ref="reason"
                    v-model="reason"
                    :autofocus="isDialogOpen"
                    :rules="nullRules"
                    @blur="blur"
                  />
                </v-col>
              </v-row>
            </v-container>
          </v-card-text>
          <v-card-actions>
            <v-spacer />
            <simple-button
              :label="'戻る'"
              :variant="'secondary'"
              @click="isDialogOpen = false"
            />
            <simple-button
              :label="toStateName(toState) + 'にする'"
              :variant="'secondary'"
              :disabled="!valid"
              @click="postReason"
            />
          </v-card-actions>
        </v-form>
      </v-card>
    </v-dialog>
  </div>
</template>
<script>
import SimpleButton from "@/views/shared/SimpleButton.vue";
import axios from "axios";
import { mapActions } from "vuex";

export default {
  components: {
    SimpleButton
  },
  props: {
    toState: {
      type: String,
      default: ""
    }
  },
  data: () => ({
    valid: true,
    isDialogOpen: false,
    reason: "",
    nullRules: [v => !!v || ""]
  }),
  watch: {
    isDialogOpen: function () {
      // @ts-ignore
      if (this.isDialogOpen) {
        let self = this;
        this.$nextTick().then(function () {
          self.$refs.form.reset();
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
    async postReason() {
      if (this.$refs.form.validate()) {
        await axios
          .put(
            "../api/applications/" +
              this.$store.state.application_detail_paper.core.application_id +
              "/states",
            {
              to_state: this.toState,
              reason: this.reason
            }
          )
          .catch(e => {
            alert(e);
            return;
          });
        this.$refs.form.reset();
        this.isDialogOpen = false;
        this.getApplicationDetail(
          this.$store.state.application_detail_paper.core.application_id
        );
      }
    },
    toStateName: function (toState) {
      switch (toState) {
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

<style lang="scss" module>
.container {
  justify-content: center;
}
</style>
