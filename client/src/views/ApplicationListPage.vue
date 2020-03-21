<template>
  <div class="application-list">
    <v-container>
      <v-row>
        <!-- 絞り込み検索欄 -->
        <v-col cols="12" md="3">
          <v-card width="1200px" class="mx-auto mt-5" color="primary">
            <v-card-title style="color:white;font-weight:bold;font-size:1.5em">
              絞り込み
              <v-btn icon @click="show = !show">
                <v-icon>{{
                  show ? "mdi-chevron-up" : "mdi-chevron-down"
                }}</v-icon>
              </v-btn>
            </v-card-title>
            <v-card-text v-show="show" style="background:white" class="pt-4">
              <v-row>
                <v-btn>更新</v-btn>
                <v-btn>条件削除</v-btn>
              </v-row>
              <v-row>
                <v-col cols="5">
                  <v-text-field placeholder="2019/01/01"></v-text-field>
                </v-col>
                <v-col cols="2">
                  <v-text-field placeholder="～" disabled></v-text-field>
                </v-col>
                <v-col cols="5">
                  <v-text-field placeholder="2019/01/01"></v-text-field>
                </v-col>
              </v-row>
              <v-text-field placeholder="2020" label="年度"></v-text-field>
              <v-select
                :items="['', 'club', 'contest', 'event', 'public']"
                label="申請書タイプ"
              ></v-select>
              <v-select
                :items="[
                  'submitted',
                  'rejected',
                  'fix_required',
                  'accepted',
                  'fully_repaid'
                ]"
                label="現在の状態"
              ></v-select>
              <v-autocomplete
                :items="userList"
                label="申請者"
                item-text="trap_id"
                item-value="trap_id"
                multiple
              ></v-autocomplete>
            </v-card-text>
          </v-card>
        </v-col>
        <!-- 申請書一覧 -->
        <v-col cols="12" md="7">
          <v-card width="1200px" class="mx-auto mt-5" color="primary">
            <v-card-title style="color:white;font-weight:bold;font-size:1.5em"
              >申請書一覧</v-card-title
            >
            <v-card-text class="pl-0 pr-0 pb-0">
              <v-list>
                <v-list-item-group color="primary">
                  <Application :list="header" class="pb-0 pt-0"></Application>
                  <v-list-item
                    v-for="list in applicationList"
                    v-bind:key="list.id"
                    :to="'/applications/' + list.id"
                    class="pl-0 pr-0"
                  >
                    <Application :list="list"> </Application>
                  </v-list-item>
                </v-list-item-group>
              </v-list>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>
<script>
import { mapState } from "vuex";
import Application from "./components/Application";
export default {
  name: "ApplicationList",
  computed: {
    ...mapState(["applicationList", "userList"])
  },
  data() {
    return {
      header: {
        title: "タイトル",
        name: "申請者名",
        money: "金額",
        state: ""
      },
      show: null
    };
  },
  created() {
    this.show = this.defaultShow();
  },
  methods: {
    /**
     * 絞り込み画面表示の初期値を画面のサイズによって変える
     */
    defaultShow() {
      switch (this.$vuetify.breakpoint.name) {
        case "xs":
        case "sm":
        case "md":
          return false;
        default:
          return true;
      }
    }
  },
  components: {
    Application
  }
};
</script>
