<template>
  <!-- repuid_to_userについては実際のサーバーでうまくいくか確認する。おそらくリスト取得がフォーカス後なのでうまくいかない -->
  <v-container>
    <v-form ref="form" v-model="valid" lazy-validation>
      <v-card class="ml-2 mr-2 mt-2 pa-3" tile>
        <v-row class="ml-4 mr-4" :justify="`space-between`">
          <v-col cols="12" sm="8" class="pt-0 pb-0">
            <h1>
              <v-row>
                <v-col cols="8">
                  <v-select
                    v-model="type_object"
                    :items="types"
                    item-text="jpn"
                    item-value="type"
                    label="Select"
                    persistent-hint
                    return-object
                    single-line
                    dense
                  ></v-select>
                </v-col>
                <v-col cols="4">申請</v-col>
              </v-row>
            </h1>
          </v-col>

          <v-col cols="12" sm="4" class="pt-0 pb-0">
            <div>申請日: {{ returnDate(this.detail.core.created_at) }}</div>
            <v-divider></v-divider>
            <div>
              申請者:
              <Icon :user="this.detail.core.applicant.trap_id" :size="20" />
              {{ this.detail.core.applicant.trap_id }}
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
            v-model="title_change"
            :rules="nullRules"
            label="概要"
            filled
            :placeholder="returnTitlePlaceholder(this.type_object.type)"
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
                    color="primary"
                  ></v-text-field>
                </template>
                <v-date-picker
                  v-model="paid_at_change"
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
                v-model="amount_change"
                :rules="amountRules"
                type="number"
                label="支払金額"
                placeholder="100"
                class="pa-0"
                height="25"
                suffix="円"
              ></v-text-field>
            </v-col>
          </v-row>
        </div>

        <div>
          <v-autocomplete
            ref="traPID"
            v-model="repaid_to_id_change"
            :rules="[
              () =>
                !(repaid_to_id_change === 0) || '返金対象者は一人以上必要です'
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
            v-model="remarks_change"
            :rules="nullRules"
            filled
            :label="returnRemarksTitle(this.type_object.type)"
            :placeholder="returnRemarksPlaceholder(this.type_object.type)"
            :hint="returnRemarksHint(this.type_object.type)"
            auto-grow
          ></v-textarea>
        </div>

        <div>
          <h3 class="ml-0 mr-0">画像</h3>
          <div
            :class="$style.image"
            :key="path"
            v-for="(path, index) in this.detail.core.images"
          >
            <span v-if="images_change[index]">
              <v-btn
                rounded
                color="primary"
                name="delete"
                @click="deleteImage(index)"
              >
                delete
              </v-btn>
              <v-img :src="`/api/images/${path}`" max-width="80%" />
            </span>
            <span v-else>
              <v-btn
                rounded
                color="primary"
                name="cancel"
                @click="cancelDeleteImage(index)"
              >
                cancel
              </v-btn>
            </span>
          </div>

          <h3 class="ml-0 mr-0">画像を追加</h3>
          <image-uploader v-model="imageBlobs" />
        </div>
      </v-card>

      <!-- todo focusしていないところのvalidateが機能していない -->

      <v-btn :disabled="!valid" @click.stop="submit" class="ma-3"
        >修正する
      </v-btn>
      <v-btn class="ma-3" @click="deleteFix">取り消す</v-btn>
    </v-form>
    <!-- ここ作成したらokを押しても押さなくても自動遷移 -->
    <v-snackbar v-model="snackbar">
      変更できました
      <v-btn
        :to="`/applications/` + response.application_id"
        color="green darken-1"
        text
        @click="afterChange()"
        >OK
      </v-btn>
    </v-snackbar>
  </v-container>
</template>

<script>
import axios from "axios";
import Icon from "../shared/Icon";
import ImageUploader from "../shared/ImageUploader";
import { mapActions } from "vuex";
import { mapState, mapMutations } from "vuex";
import {
  titlePlaceholder,
  remarksPlaceholder,
  remarksHint
} from "../../use/inputFormText";
import { remarksTitle } from "../../use/applicationDetail";
import { dayPrint } from "../../use/dataFormat";

export default {
  data: function () {
    return {
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
      type_object: { jpn: "", type: "" },
      types: [
        { jpn: "部費利用", type: "club" },
        { jpn: "大会等旅費補助", type: "contest" },
        { jpn: "イベント交通費補助", type: "event" },
        { jpn: "渉外交通費補助", type: "public" }
      ],
      snackbar: false,
      menu: false,
      valid: true,
      amountRules: [
        v => !!v || "",
        v => !!String(v).match("^[1-9][0-9]*$") || "金額が不正です"
      ],
      nullRules: [v => !!v || ""],
      type_change: "",
      title_change: "",
      remarks_change: "",
      paid_at_change: "",
      amount_change: "",
      repaid_to_id_change: [],
      Images: [],
      imageBlobs: [],
      // todo返金リスト配列
      changeRules: [v => (v !== this.detail.core.repayment_logs && !!v) || ""]
    };
  },
  async created() {
    this.title_change = this.detail.core.current_detail.title;
    this.type_object.type = this.detail.core.current_detail.type;
    this.title_change = this.detail.core.current_detail.title;
    this.remarks_change = this.detail.core.current_detail.remarks;
    this.paid_at_change = this.detail.core.current_detail.paid_at;
    this.amount_change = this.detail.core.current_detail.amount;
    this.images_change = new Array(this.detail.core.images.length);
    this.images_change.fill(true);
    await this.getUsers();
    const trap_ids = this.detail.core.repayment_logs.map(
      log => log.repaid_to_user.trap_id
    );
    this.repaid_to_id_change = trap_ids;
  },
  mounted() {},
  computed: {
    ...mapState({ detail: "application_detail_paper" }),
    computedDateFormatted() {
      return this.formatDate(this.paid_at_change);
    },
    me() {
      return this.$stor.me;
    },
    form() {
      return {
        repaid_to_id_change: this.repaid_to_id_change
      };
    },
    traPIDs() {
      let trap_ids = [];
      for (let i = 0; i < this.$store.state.userList.length - 1; i++) {
        trap_ids[i] = this.$store.state.userList[i].trap_id;
      }
      return trap_ids;
    }
  },

  methods: {
    ...mapMutations(["deleteFix"]),
    ...mapActions({
      getUsers: "getUserList",
      getApplicationDetail: "getApplicationDetail"
    }),
    async submit() {
      if (this.$refs.form.validate()) {
        this.images_change.forEach((flag, index) => {
          if (!flag) {
            axios.delete("/api/images/" + this.detail.core.images[index]);
          }
        });
        let form = new FormData();
        let date = new Date(this.paid_at_change);
        let details = {
          type: this.type_object.type,
          title: this.title_change,
          remarks: this.remarks_change,
          paid_at: date.toISOString(),
          amount: Number(this.amount_change),
          repaid_to_id: this.repaid_to_id_change
        };
        form.append("details", JSON.stringify(details));
        this.imageBlobs.forEach(imageBlob => {
          form.append("images", imageBlob);
        });
        const response = await axios.patch(
          "/api/applications/" + this.detail.core.application_id,
          form,
          {
            headers: { "content-type": "multipart/form-data" }
          }
        );
        this.response = response.data;
        await this.getApplicationDetail(this.$route.params.id);
        this.snackbar = true;
      }
    },
    formatDate(date) {
      if (!date) return null;

      const [year, month, day] = date.split("-");
      return `${year}年${month.replace(/^0/, "")}月${day.replace(/^0/, "")}日`;
    },
    returnDate: function (date) {
      return dayPrint(date);
    },
    returnRemarksTitle: function (type) {
      return remarksTitle(type);
    },
    returnTitlePlaceholder: function (type) {
      return titlePlaceholder(type);
    },
    returnRemarksPlaceholder: function (type) {
      return remarksPlaceholder(type);
    },
    returnRemarksHint: function (type) {
      return remarksHint(type);
    },
    deleteImage(index) {
      this.images_change[index] = false;
      this.$forceUpdate();
    },
    cancelDeleteImage(index) {
      this.images_change[index] = true;
      this.$forceUpdate();
    },
    afterChange() {
      this.snackbar = false;
      this.deleteFix();
    }
  },
  props: {},
  components: {
    Icon,
    ImageUploader
  }
};
</script>

<style lang="scss" module>
.image {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
}
</style>
