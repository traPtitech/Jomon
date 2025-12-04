<template>
  <div class="application-list">
    <v-container>
      <v-row justify="space-around">
        <!-- 絞り込み検索欄 -->
        <v-col cols="12" md="3">
          <v-card class="mx-auto mt-5" color="primary">
            <v-card-title
              style="color: white; font-weight: bold; font-size: 1.5em"
            >
              絞り込み
              <v-btn variant="text" icon color="white" @click="show = !show">
                <v-icon color="white">
                  {{ show ? "mdi-chevron-up" : "mdi-chevron-down" }}
                </v-icon>
              </v-btn>
            </v-card-title>
            <v-card-text
              v-show="show"
              style="background: white; color: black"
              class="pt-4"
            >
              <v-form>
                <div>
                  <v-btn
                    color="primary"
                    class="ma-1"
                    @click="getApplicationList(params)"
                  >
                    <v-icon>mdi-reload</v-icon>
                  </v-btn>
                  <v-btn color="primary" class="ma-1" @click="resetParams()">
                    <v-icon>mdi-close</v-icon>
                  </v-btn>
                </div>
                <div>
                  <v-btn
                    outlined
                    color="primary"
                    class="ma-1"
                    @click="sortCreatedAt()"
                  >
                    日付順
                  </v-btn>
                  <v-btn
                    outlined
                    color="primary"
                    class="ma-1"
                    @click="sortTitle()"
                  >
                    タイトル順
                  </v-btn>
                </div>
                <div :class="$style.date_range">
                  <v-text-field
                    v-model="params.submitted_since"
                    placeholder="2019-01-01"
                    :rules="dayRule"
                  />
                  <span :class="$style.tilde">〜</span>
                  <v-text-field
                    v-model="params.submitted_until"
                    placeholder="2019-01-01"
                    :rules="dayRule"
                  />
                </div>
                <v-text-field
                  v-model="params.financial_year"
                  :rules="yearRule"
                  placeholder="2020"
                  label="年度"
                />
                <v-select
                  v-model="params.type"
                  :items="type_items"
                  item-title="jpn"
                  item-value="type"
                  label="申請タイプ"
                />
                <v-select
                  v-model="params.current_state"
                  :items="state_items"
                  item-title="jpn"
                  item-value="state"
                  label="現在の状態"
                />
                <v-select
                  v-model="params.applicant"
                  :items="userList"
                  label="申請者"
                  item-title="trap_id"
                  item-value="trap_id"
                />
              </v-form>
            </v-card-text>
          </v-card>
        </v-col>
        <!-- 申請一覧 -->
        <v-col cols="12" md="7">
          <v-card class="mx-auto mt-5" color="primary">
            <v-card-title
              style="color: white; font-weight: bold; font-size: 1.5em"
            >
              申請一覧
            </v-card-title>
            <v-card-text class="pl-0 pr-0 pb-0">
              <v-list>
                <template v-if="applicationList.length > 0">
                  <Application :list="header" class="pb-0 pt-0" />
                  <v-list-item
                    v-for="(list, index) in applicationList"
                    :key="index"
                    :to="'/applications/' + list.application_id"
                    class="pl-0 pr-0"
                  >
                    <Application :list="list" />
                  </v-list-item>
                </template>
                <div v-else>該当する申請はありません。</div>
              </v-list>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script setup lang="ts">
import { useApplicationListStore } from "@/stores/applicationList";
import { useUserListStore } from "@/stores/userList";
import { storeToRefs } from "pinia";
import { onMounted, reactive, ref } from "vue";
import { useDisplay } from "vuetify";
import Application from "./components/Application.vue";

const applicationListStore = useApplicationListStore();
const userListStore = useUserListStore();

const { applicationList } = storeToRefs(applicationListStore);
const { userList } = storeToRefs(userListStore);
const { fetchApplicationList } = applicationListStore;
const { fetchUserList } = userListStore;

const display = useDisplay();

const sort = {
  created_at: "created_at",
  inv_created_at: "-created_at",
  title: "title",
  inv_title: "-title"
};

const type_items = [
  { type: "club", jpn: "部費利用申請" },
  { type: "contest", jpn: "大会等旅費補助申請" },
  { type: "event", jpn: "イベント交通費補助申請" },
  { type: "public", jpn: "渉外交通費補助" }
];

const state_items = [
  { state: "submitted", jpn: "提出済み" },
  { state: "rejected", jpn: "却下" },
  { state: "fix_required", jpn: "要修正" },
  { state: "accepted", jpn: "払い戻し待ち" },
  { state: "fully_repaid", jpn: "払い戻し完了" }
];

const header = {
  current_detail: {
    title: "タイトル",
    amount: "金額"
  },
  applicant: {
    trap_id: "申請者ID"
  },
  current_state: ""
};

const show = ref(false);
const params = reactive({
  sort: sort.created_at,
  current_state: "",
  financial_year: "",
  applicant: "",
  type: "",
  submitted_since: "",
  submitted_until: ""
});

const dayRule = [
  (value: string) =>
    !value || /^[0-9]{4}-[0-9]{1,2}-[0-9]{1,2}$/.test(value) || "Invalid Day."
];
const yearRule = [
  (value: string) => !value || /^[0-9]{4}$/.test(value) || "Invalid Year."
];

const defaultShow = () => {
  switch (display.name.value) {
    case "xs":
    case "sm":
      return false;
    default:
      return true;
  }
};

const resetParams = () => {
  Object.assign(params, {
    sort: sort.created_at,
    current_state: "",
    financial_year: "",
    applicant: "",
    type: "",
    submitted_since: "",
    submitted_until: ""
  });
};

const sortCreatedAt = () => {
  if (params.sort === sort.created_at) params.sort = sort.inv_created_at;
  else params.sort = sort.created_at;
  fetchApplicationList(params);
};

const sortTitle = () => {
  if (params.sort === sort.title) params.sort = sort.inv_title;
  else params.sort = sort.title;
  fetchApplicationList(params);
};

const getApplicationList = (p: Record<string, unknown>) => {
  fetchApplicationList(p);
};

onMounted(async () => {
  const p1 = fetchApplicationList({});
  const p2 = fetchUserList();
  await Promise.all([p1, p2]);
  show.value = defaultShow();
});
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
