<template>
  <v-container>
    <v-form ref="form" v-model="valid" lazy-validation>
      <v-card class="ml-2 mr-2 mt-2 pa-3" tile>
        <v-row class="ml-4 mr-4" :justify="`space-between`">
          <v-col cols="12" sm="8" class="pt-0 pb-0">
            <h1>{{ returnType($route.params.type) }}申請</h1>
          </v-col>

          <v-col cols="12" sm="4" class="pt-0 pb-0">
            <div>申請日: {{ returnToday() }}</div>
            <v-divider></v-divider>
            <div>
              申請者:<Icon :user="this.$store.state.me.trap_id" :size="20" />{{
                this.$store.state.me.trap_id
              }}
            </div>
            <div>
              <v-divider></v-divider>
            </div>
          </v-col>
        </v-row>

        <template>
          <v-divider class="mt-1 mb-2"></v-divider>
        </template>

        <div>
          <v-text-field
            v-model="title"
            :rules="nullRules"
            label="概要"
            filled
            :placeholder="returnTitlePlaceholder($route.params.type)"
            ref="firstfocus"
          ></v-text-field>
        </div>

        <div>
          <v-row>
            <v-col cols="10" sm="5" class="pb-0 pt-0">
              <v-menu
                v-model="menu"
                :close-on-content-click="false"
                transition="scale-transition"
                offset-y
                max-width="290px"
                min-width="290px"
              >
                <template v-slot:activator="{ on }">
                  <v-text-field
                    v-model="computedDateFormatted"
                    :rules="nullRules"
                    label="支払日"
                    filled
                    readonly
                    placeholder="2020年5月2日"
                    v-on="on"
                    height="10"
                  ></v-text-field>
                </template>
                <v-date-picker
                  v-model="date"
                  no-title
                  color="primary"
                  @input="menu = false"
                ></v-date-picker>
              </v-menu>
            </v-col>
          </v-row>
        </div>

        <div>
          <v-row align="center">
            <v-col cols="10" sm="5" class="pb-0 pt-0">
              <v-text-field
                v-model="amount"
                :rules="amountRules"
                label="支払金額"
                filled
                type="number"
                placeholder="100"
                class="pa-0"
                height="25"
                suffix="円"
              ></v-text-field
            ></v-col>
          </v-row>
        </div>

        <div>
          <v-autocomplete
            ref="traPID"
            v-model="traPID"
            :rules="[
              () => !(traPID.length == 0) || '返金対象者は一人以上必要です'
            ]"
            label="返金対象者"
            filled
            :items="traPIDs"
            placeholder="traQIDs"
            hint="traQ IDの一部入力で候補が表示されます"
            required
            multiple
          >
          </v-autocomplete>
        </div>

        <div>
          <v-textarea
            v-model="remarks"
            :rules="nullRules"
            filled
            :label="returnRemarksTitle($route.params.type)"
            :placeholder="returnRemarksPlaceholder($route.params.type)"
            :hint="returnRemarksHint($route.params.type)"
            auto-grow
          ></v-textarea>
        </div>

        <div>
          <image-uploader v-model="imageBlobs" />
        </div>
      </v-card>

      <!-- todo focusしていないところのvalidateが機能していない -->
      <v-btn :disabled="!valid" @click.stop="submit" class="ma-3"
        >作成する</v-btn
      >
    </v-form>
    <!-- ここ作成したらokを押しても押さなくても自動遷移 -->
    <v-snackbar v-model="snackbar">
      作成できました
      <v-btn
        :to="`../../applications/` + response.application_id"
        color="green darken-1"
        text
        @click="sackbar = false"
        >OK</v-btn
      >
    </v-snackbar>
  </v-container>
</template>

<script>
import axios from "axios";
import Icon from "./shered/Icon";
import ImageUploader from "./shered/ImageUploader";
import { mapActions, mapGetters } from "vuex";
import {
  titlePlaceholder,
  remarksPlaceholder,
  remarksHint
} from "../use/inputFormText";
import { remarksTitle, applicationType } from "../use/applicationDetail";
import { dayPrint } from "../use/dataFormat";
export default {
  data: () => ({
    response: {
      application_id: null,
      applicant: { trapid: null },
      created_at: null,
      current_detail: {
        title: null,
        type: null,
        amount: 0,
        remarks: null,
        created_at: null,
        paid_at: null
      }
    },
    snackbar: false,
    date: null,
    menu: false,
    traPID: [],
    valid: true,
    title: "",
    amount: "",
    remarks: "",
    imageBlobs: [],
    amountRules: [
      v => !!v || "必須の項目です",
      v => !!String(v).match("^[1-9][0-9]*$") || "金額が不正です"
    ],
    nullRules: [v => !!v || "必須の項目です"]
  }),
  mounted() {
    this.$refs.firstfocus.focus();
  },
  computed: {
    ...mapGetters({ traPIDs: "trap_ids" }),
    computedDateFormatted() {
      return this.formatDate(this.date);
    },
    form() {
      return {
        traPID: this.traPID
      };
    }
  },

  // todo 返金対象者周りのポスト等
  async created() {
    await this.getUsers();
  },
  methods: {
    ...mapActions({
      getUsers: "getUserList"
    }),
    submit() {
      if (this.$refs.form.validate()) {
        let form = new FormData();
        let paid_at = new Date(this.date);
        let details = {
          type: this.$route.params.type,
          title: this.title,
          remarks: this.remarks,
          paid_at: paid_at.toISOString(),
          amount: Number(this.amount),
          repaid_to_id: this.traPID
        };
        form.append("details", JSON.stringify(details));
        this.imageBlobs.forEach(imageBlob => {
          form.append("images", imageBlob);
        });
        axios
          .post("/api/applications", form, {
            headers: { "content-type": "multipart/form-data" }
          })
          .then(
            response => (
              (this.response = response.data), (this.snackbar = true)
            )
          );
      }
    },
    formatDate(date) {
      if (!date) return null;

      const [year, month, day] = date.split("-");
      return `${year}年${month.replace(/^0/, "")}月${day.replace(/^0/, "")}日`;
    },
    returnToday: function() {
      const date = new Date();
      return dayPrint(date);
    },
    returnType: function(type) {
      return applicationType(type);
    },
    returnRemarksTitle: function(type) {
      return remarksTitle(type);
    },
    returnTitlePlaceholder: function(type) {
      return titlePlaceholder(type);
    },
    returnRemarksPlaceholder: function(type) {
      return remarksPlaceholder(type);
    },
    returnRemarksHint: function(type) {
      return remarksHint(type);
    }
  },
  props: {},
  components: {
    Icon,
    ImageUploader
  }
};
</script>
