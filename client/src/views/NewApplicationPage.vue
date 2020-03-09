<template>
  <div class="new-applicatoin">
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
        <v-text-field label="入力してください"></v-text-field>
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
                <v-col height="100%" cols="8" md="6">
                  <v-card class="pa-2" outlined tile>
                    <Icon :user="`series2`" :size="20" />series2
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
                    <v-row class="ma-1" align="center">
                      <v-col class="pt-0 pb-0" cols="10">
                        <v-text-field
                          type="number"
                          placeholder="金額入力"
                          hide-details
                          class="pa-0"
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
                          readonly
                          placeholder="支払日選択"
                          v-on="on"
                          height="10"
                          hide-details
                        ></v-text-field>
                      </template>
                      <v-date-picker
                        v-model="date"
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
        <h3>{{ returnRemarkTitle($route.params.type) }}:</h3>
        <v-textarea placeholder="入力してください" auto-grow></v-textarea>
      </v-row>
      <h3 class="ml-0 mr-0">払い戻し対象者</h3>
      リストです
      <v-autocomplete
        ref="traPID"
        v-model="traPID"
        :rules="[() => !!traPID || '返金対象者は一人以上必要です']"
        :items="traPIDs"
        placeholder="返金対象者のtraPidを入力..."
        required
        multiple
      >
        <!-- <Icon slot="prepend" :user="traPID" :size="25" /> -->
      </v-autocomplete>

      <!-- <v-btn>返金対象者追加</v-btn> -->

      <h3 class="ml-0 mr-0">申請書画像リスト</h3>

      画像リスト(画像アップロード)
    </v-card>
    <v-btn class="ma-3">作成する</v-btn>
  </div>
</template>

<script>
import Icon from "./components/Icon";
export default {
  data: () => ({
    date: null,
    menu: false,
    traPID: null
  }),
  computed: {
    computedDateFormatted() {
      return this.formatDate(this.date);
    },
    me() {
      return this.$stor.me;
    },
    form() {
      return {
        traPID: this.traPID
      };
    },
    traPIDs() {
      return this.$store.getters.userList;
    }
  },

  methods: {
    formatDate(date) {
      if (!date) return null;

      const [year, month, day] = date.split("-");
      return `${year}年${month}月${day}日`;
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
    Icon
  }
};
</script>
