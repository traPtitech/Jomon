<template>
  <!-- repuid_to_userについては実際のサーバーでうまくいくか確認する。おそらくリスト取得がフォーカス後なのでうまくいかない -->
  <div class="fix-applicatoin">
    <v-form ref="form" v-model="valid" lazy-validation>
      <v-card class="ml-2 mr-2 mt-2 pa-3" tile>
        <v-row class="ml-4 mr-4" :justify="`space-between`">
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
          <h1>申請書</h1>

          <div>
            <div>申請書ID: {{ this.detail.core.application_id }}</div>
            <v-divider></v-divider>
          </div>
        </v-row>

        <template>
          <v-divider></v-divider>
        </template>
        <v-row class="ml-0 mr-0">
          <h1>タイトル:</h1>
          <v-text-field
            v-model="title_change"
            :rules="nullRules"
            label="入力してください"
          ></v-text-field>
        </v-row>

        <div>
          <v-container class="pa-0">
            <v-row>
              <!-- 以下は左列 -->
              <v-col cols="12" md="6">
                <v-row no-gutters>
                  <v-col cols="4" md="6">
                    <v-card height="100%" class="pa-2" outlined tile>
                      申請者trapid
                    </v-card>
                  </v-col>
                  <v-col cols="8" md="6">
                    <v-card height="100%" class="pa-2" outlined tile>
                      <Icon
                        :user="this.detail.core.applicant.trap_id"
                        :size="20"
                      />{{ this.detail.core.applicant.trap_id }}
                    </v-card>
                  </v-col>
                </v-row>
                <v-row no-gutters>
                  <v-col cols="4" md="6">
                    <v-card height="100%" class="pa-2" outlined tile>
                      申請金額
                    </v-card>
                  </v-col>
                  <v-col cols="8" md="6">
                    <v-card height="100%" class="pa-0" outlined tile>
                      <v-row class="pr-2 pl-2" align="center">
                        <v-col class="pb-1 pt-2" cols="10">
                          <v-text-field
                            v-model="amount_change"
                            :rules="amountRules"
                            type="number"
                            label="金額入力"
                            hide-details
                            class="pa-0"
                            height="25"
                          ></v-text-field
                        ></v-col>
                        <v-col class="pt-0 pb-0" cols="2">円</v-col>
                      </v-row>
                    </v-card>
                  </v-col>
                </v-row>
              </v-col>
              <!-- 以上左列以下右列 -->
              <v-col cols="12" md="6">
                <v-row no-gutters>
                  <v-col cols="4" md="6">
                    <v-card height="100%" class="pa-2" outlined tile>
                      申請書作成日
                    </v-card>
                  </v-col>
                  <v-col height="100%" cols="8" md="6">
                    <v-card class="pa-2" outlined tile>
                      {{ returnDate(this.detail.core.created_at) }}
                    </v-card>
                  </v-col>
                </v-row>
                <v-row no-gutters>
                  <v-col cols="4" md="6">
                    <v-card height="100%" class="pa-2" outlined tile>
                      支払った日
                    </v-card>
                  </v-col>
                  <v-col cols="8" md="6">
                    <v-card height="100%" class="pa-2" outlined tile>
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
                            readonly
                            label="支払日選択"
                            v-on="on"
                            height="10"
                            hide-details
                            color="primary"
                          ></v-text-field>
                        </template>
                        <v-date-picker
                          v-model="paid_at_change"
                          no-title
                          @input="menu = false"
                        ></v-date-picker>
                      </v-menu>
                    </v-card>
                  </v-col>
                </v-row>
              </v-col>
              <!-- 以上右列 -->
            </v-row>
          </v-container>
        </div>
        <v-row class="ml-0 mr-0">
          <h3>{{ returnRemarkTitle(this.type_object.type) }}:</h3>
          <v-textarea
            v-model="remarks_change"
            :rules="nullRules"
            label="入力してください"
            auto-grow
          ></v-textarea>
        </v-row>
        <h3 class="ml-0 mr-0">払い戻し対象者</h3>
        <v-autocomplete
          ref="traPID"
          v-model="repaid_to_id_change"
          :rules="[
            () => !!repaid_to_id_change || '返金対象者は一人以上必要です'
          ]"
          :items="traPIDs"
          label="返金対象者のtraPidを入力..."
          required
          multiple
        >
        </v-autocomplete>
        <h3 class="ml-0 mr-0">申請書画像リスト</h3>
        <div :key="path" v-for="(path, index) in this.detail.core.images">
          <span v-if="images_change[index]">
            <v-btn
              rounded
              color="primary"
              name="delete"
              @click="deleteImage(index)"
            >
              delete
            </v-btn>
            <v-img :src="'/api/images/' + path" max-width="80%" />
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
      </v-card>

      <!-- todo focusしていないところのvalidateが機能していない -->

      <v-btn :disabled="!valid" @click.stop="submit" class="ma-3"
        >修正する</v-btn
      ><v-btn class="ma-3" @click="deleteFix">取り消す</v-btn>
      <v-dialog persistent v-model="open_dialog">
        <v-card class="pa-3">
          <v-card-title class="headline">以下の内容で修正しました</v-card-title>
          <v-row :justify="`space-between`">
            <v-col cols="4" md="2">申請書id</v-col>
            <v-col cols="8" md="10">{{ response.application_id }}</v-col>
          </v-row>
          <v-row :justify="`space-between`">
            <v-col cols="4" md="2">申請書タイプ</v-col>
            <v-col cols="8" md="10">{{ response.current_detail.type }}</v-col>
          </v-row>
          <v-row :justify="`space-between`">
            <v-col cols="4" md="2">タイトル</v-col>
            <v-col cols="8" md="10">{{ response.current_detail.title }}</v-col>
          </v-row>
          <v-row :justify="`space-between`">
            <v-col cols="4" md="2">申請者</v-col>
            <v-col cols="8" md="10">{{ response.applicant.trap_id }}</v-col>
          </v-row>
          <v-row :justify="`space-between`">
            <v-col cols="4" md="2">申請金額</v-col>
            <v-col cols="8" md="10">{{ response.current_detail.amount }}</v-col>
          </v-row>
          <v-row :justify="`space-between`">
            <v-col cols="4" md="2">申請書作成日</v-col>
            <v-col cols="8" md="10">{{ response.created_at }}</v-col>
          </v-row>
          <v-row :justify="`space-between`">
            <v-col cols="4" md="2">支払った日</v-col>
            <v-col cols="8" md="10">{{
              response.current_detail.paid_at
            }}</v-col>
          </v-row>
          <v-row :justify="`space-between`">
            <v-col cols="4" md="2">詳細</v-col>
            <v-col cols="8" md="10">{{ response.current_detail.title }}</v-col>
          </v-row>
          <v-row :justify="`space-between`">
            <v-col cols="4" md="2">返金対象者</v-col>
            <v-col cols="8" md="10">保留</v-col>
          </v-row>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn
              :to="`../../applications/` + response.application_id"
              color="primary"
              text
              @click="[(open_dialog = false), deleteFix()]"
              >OK</v-btn
            >
          </v-card-actions>
        </v-card>
      </v-dialog>
    </v-form>
  </div>
</template>

<script>
import axios from "axios";
import Icon from "../shered/Icon";
import ImageUploader from "../shered/ImageUploader";
import { mapActions } from "vuex";
import { mapState, mapMutations } from "vuex";
export default {
  data: function() {
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
      open_dialog: false,
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
      let trap_ids = new Array();
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
        this.open_dialog = true;
      }
    },
    formatDate(date) {
      if (!date) return null;

      const [year, month, day] = date.split("-");
      return `${year}年${month.replace(/^0/, "")}月${day.replace(/^0/, "")}日`;
    },
    returnDate: function(date) {
      const normalizedDate = new Date(date);
      return (
        normalizedDate.getFullYear() +
        "年" +
        (normalizedDate.getMonth() + 1) +
        "月" +
        normalizedDate.getDate() +
        "日"
      );
    },
    returnType: function(type) {
      switch (type) {
        case "club":
          return "部費利用";
        case "contest":
          return "大会等旅費補助";
        case "event":
          return "イベント交通費補助";
        case "public":
          return "渉外交通費補助";
        default:
          return "タイプが間違っています";
      }
    },
    returnRemarkTitle: function(type) {
      switch (type) {
        case "club":
          return "購入したものの概要";
        case "contest":
          return "旅程";
        case "event":
          return "乗車区間";
        case "public":
          return "乗車区間";
        default:
          return "タイプが間違っています";
      }
    },
    deleteImage(index) {
      this.images_change[index] = false;
      this.$forceUpdate();
    },
    cancelDeleteImage(index) {
      this.images_change[index] = true;
      this.$forceUpdate();
    }
  },
  props: {},
  components: {
    Icon,
    ImageUploader
  }
};
</script>
