<template>
  <div class="new-applicatoin">
    <v-form ref="form" v-model="valid" lazy-validation>
      <v-card class="ml-2 mr-2 mt-2 pa-3" tile>
        <v-row class="ml-4 mr-4" :justify="`space-between`">
          <h1>{{ returnType($route.params.type) }}申請書</h1>
          <div>
            <div>申請書ID: 自動入力されます</div>
            <v-divider></v-divider>
          </div>
        </v-row>

        <template>
          <v-divider></v-divider>
        </template>
        <v-row class="ml-0 mr-0">
          <h1>タイトル:</h1>
          <v-text-field
            v-model="title"
            :rules="nullRules"
            label="入力してください"
            ref="firstfocus"
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
                      <Icon :user="this.$store.state.me.trap_id" :size="20" />{{
                        this.$store.state.me.trap_id
                      }}
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
                            v-model="amount"
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
                      {{ returnDate(new Date()) }}
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
                          ></v-text-field>
                        </template>
                        <v-date-picker
                          v-model="date"
                          no-title
                          color="primary"
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
          <h3>{{ returnRemarkTitle($route.params.type) }}:</h3>
          <v-textarea
            v-model="remarks"
            :rules="nullRules"
            label="入力してください"
            auto-grow
          ></v-textarea>
        </v-row>
        <h3 class="ml-0 mr-0">払い戻し対象者</h3>
        <v-autocomplete
          ref="traPID"
          v-model="traPID"
          :rules="[() => !!traPID || '返金対象者は一人以上必要です']"
          :items="traPIDs"
          label="返金対象者のtraPidを入力..."
          required
          multiple
        >
        </v-autocomplete>

        <h3 class="ml-0 mr-0">申請書画像リスト</h3>
        <image-uploader v-model="imageBlobs" />
      </v-card>

      <!-- todo focusしていないところのvalidateが機能していない -->
      <v-btn :disabled="!valid" @click.stop="submit" class="ma-3"
        >作成する</v-btn
      >
    </v-form>
    <!-- ここ作成したらokを押しても押さなくても自動遷移 -->
    <v-snacker v-model="snacker">
      作成できました。
      <v-btn
        :to="`../../applications/` * +response.application_id"
        color="green darken-1"
        text
        @click="sacker = false"
        >OK</v-btn
      >
    </v-snacker>
  </div>
</template>

<script>
import axios from "axios";
import Icon from "./shered/Icon";
import ImageUploader from "./shered/ImageUploader";
import { mapActions, mapGetters } from "vuex";
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
      v => !!v || "",
      v => !!String(v).match("^[1-9][0-9]*$") || "金額が不正です"
    ],
    nullRules: [v => !!v || ""]
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
    }
  },
  props: {},
  components: {
    Icon,
    ImageUploader
  }
};
</script>
