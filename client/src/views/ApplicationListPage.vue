<template>
  <div class="application-list">
    <v-container>
      <v-row justify="space-around">
        <!-- 絞り込み検索欄 -->
        <v-col cols="12" md="3">
          <v-card width="1200px" class="mx-auto mt-5" color="primary">
            <v-card-title style="color: white; font-weight: bold; font-size: 1.5em">
              絞り込み
              <v-btn icon @click="show = !show">
                <v-icon color="white">{{ show ? "mdi-chevron-up" : "mdi-chevron-down" }} </v-icon>
              </v-btn>
            </v-card-title>
            <v-card-text v-show="show" style="background: white" class="pt-4">
              <v-form>
                <div>
                  <v-btn color="primary" @click="getApplicationList(params)" class="ma-1">
                    <v-icon>mdi-reload</v-icon>
                  </v-btn>
                  <v-btn color="primary" @click="resetParams()" class="ma-1">
                    <v-icon>mdi-close</v-icon>
                  </v-btn>
                </div>
                <div>
                  <v-btn outlined color="primary" @click="sortCreatedAt()" class="ma-1"
                    >日付順
                  </v-btn>
                  <v-btn outlined color="primary" @click="sortTitle()" class="ma-1"
                    >タイトル順
                  </v-btn>
                </div>
                <div :class="$style.date_range">
                  <v-text-field
                    v-model="params.submitted_since"
                    placeholder="2019-01-01"
                    :rules="dayRule"
                  ></v-text-field>
                  <span :class="$style.tilde">〜</span>
                  <v-text-field
                    v-model="params.submitted_until"
                    placeholder="2019-01-01"
                    :rules="dayRule"
                  ></v-text-field>
                </div>
                <v-text-field
                  v-model="params.financial_year"
                  :rules="yearRule"
                  placeholder="2020"
                  label="年度"
                ></v-text-field>
                <v-select
                  v-model="params.type"
                  :items="type_items"
                  item-text="jpn"
                  item-value="type"
                  label="申請タイプ"
                ></v-select>
                <v-select
                  v-model="params.current_state"
                  :items="state_items"
                  item-text="jpn"
                  item-value="state"
                  label="現在の状態"
                ></v-select>
                <v-select
                  v-model="params.applicant"
                  :items="userList"
                  label="申請者"
                  item-text="trap_id"
                  item-value="trap_id"
                ></v-select>
              </v-form>
            </v-card-text>
          </v-card>
        </v-col>
        <!-- 申請一覧 -->
        <v-col cols="12" md="7">
          <v-card width="1200px" class="mx-auto mt-5" color="primary">
            <v-card-title style="color: white; font-weight: bold; font-size: 1.5em"
              >申請一覧
            </v-card-title>
            <v-card-text class="pl-0 pr-0 pb-0">
              <v-list>
                <v-list-item-group v-if="applicationList.length > 0" color="primary">
                  <Application :list="header" class="pb-0 pt-0"></Application>
                  <v-list-item
                    v-for="(list, index) in applicationList"
                    v-bind:key="index"
                    :to="'/applications/' + list.application_id"
                    class="pl-0 pr-0"
                  >
                    <Application :list="list"></Application>
                  </v-list-item>
                </v-list-item-group>
                <div v-else>該当する申請はありません。</div>
              </v-list>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script>
import { mapState, mapActions } from "vuex";
import Application from "./components/Application";

let sort = {
  created_at: "created_at",
  inv_created_at: "-created_at",
  title: "title",
  inv_title: "-title"
};
export default {
  name: "ApplicationList",
  computed: {
    ...mapState(["applicationList", "userList"])
  },
  methods: {
    ...mapActions(["getApplicationList", "getUserList"]),
    /**
     * 絞り込み画面表示の初期値を画面のサイズによって変える
     */
    defaultShow() {
      switch (this.$vuetify.breakpoint.name) {
        case "xs":
        case "sm":
          return false;
        default:
          return true;
      }
    },
    /**
     * 絞り込みリセット
     */
    resetParams() {
      this.params = {
        sort: sort.created_at,
        current_state: "",
        financial_year: "",
        applicant: "",
        type: "",
        submitted_since: "",
        submitted_until: ""
      };
    },
    /**
     * 作成日でソート
     */
    sortCreatedAt() {
      if (this.params.sort === sort.created_at) this.params.sort = sort.inv_created_at;
      else this.params.sort = sort.created_at;
      this.getApplicationList(this.params);
    },
    /**
     * タイトルでソート
     */
    sortTitle() {
      if (this.params.sort === sort.title) this.params.sort = sort.inv_title;
      else this.params.sort = sort.title;
      this.getApplicationList(this.params);
    }
  },
  data() {
    return {
      type_items: [
        { type: "club", jpn: "部費利用申請" },
        { type: "contest", jpn: "大会等旅費補助申請" },
        { type: "event", jpn: "イベント交通費補助申請" },
        { type: "public", jpn: "渉外交通費補助" }
      ],
      type: { type: "", jpn: "" },
      state_items: [
        { state: "submitted", jpn: "提出済み" },
        { state: "rejected", jpn: "却下" },
        { state: "fix_required", jpn: "要修正" },
        { state: "accepted", jpn: "払い戻し待ち" },
        { state: "fully_repaid", jpn: "払い戻し完了" }
      ],
      state: { state: "", jpn: "" },
      header: {
        current_detail: {
          title: "タイトル",
          amount: "金額"
        },
        applicant: {
          trap_id: "申請者ID"
        },
        current_state: ""
      },
      show: null,
      params: {
        sort: sort.created_at,
        current_state: "",
        financial_year: "",
        applicant: "",
        type: "",
        submitted_since: "",
        submitted_until: ""
      },
      dayRule: [
        value => !value || /^[0-9]{4}-[0-9]{1,2}-[0-9]{1,2}$/.test(value) || "Invalid Day."
      ],
      yearRule: [value => !value || /^[0-9]{4}$/.test(value) || "Invalid Year."]
    };
  },
  async created() {
    const p1 = this.getApplicationList({});
    const p2 = this.getUserList();
    await Promise.all([p1, p2]);
    this.show = this.defaultShow();
  },
  components: {
    Application
  }
};
</script>

<style lang="scss" module>
.date_range {
  display: grid;
  grid-template-columns: 3fr 1fr 3fr;
  align-items: baseline;
}
.tilde {
  text-align: center;
}
</style>
