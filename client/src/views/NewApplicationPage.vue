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
            <v-divider />
            <div>
              申請者:
              <Icon :user="trapId" :size="20" />
              {{ trapId }}
            </div>
            <div>
              <v-divider />
            </div>
          </v-col>
        </v-row>

        <v-divider class="mt-1 mb-2" />

        <div>
          <v-text-field
            ref="firstfocus"
            v-model="title"
            :rules="nullRules"
            label="概要"
            filled
          />
        </div>

        <div>
          <v-row>
            <v-col cols="10" sm="5" class="pb-0 pt-0">
              <v-menu
                v-model="menu"
                :close-on-content-click="false"
                transition="scale-transition"
                location="bottom"
                max-width="290px"
                min-width="290px"
              >
                <template #activator="{ props }">
                  <v-text-field
                    v-model="computedDateFormatted"
                    :rules="nullRules"
                    label="支払日"
                    filled
                    readonly
                    v-bind="props"
                  />
                </template>
                <v-date-picker
                  v-model="date"
                  no-title
                  color="primary"
                  @update:model-value="menu = false"
                />
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
                class="pa-0"
                suffix="円"
              />
            </v-col>
          </v-row>
        </div>

        <div>
          <v-autocomplete
            ref="traPID"
            v-model="traPID"
            :rules="[
              () => !(traPID.length === 0) || '返金対象者は一人以上必要です'
            ]"
            label="返金対象者"
            filled
            :items="traPIDs"
            hint="traQ IDの一部入力で候補が表示されます"
            required
            multiple
          />
        </div>

        <div>
          <v-textarea
            v-model="remarks"
            :rules="nullRules"
            filled
            :label="returnRemarksTitle($route.params.type)"
            :hint="returnRemarksHint($route.params.type)"
            auto-grow
          />
        </div>

        <div>
          <image-uploader v-model="imageBlobs" />
        </div>
      </v-card>

      <!-- todo focusしていないところのvalidateが機能していない -->
      <v-btn :disabled="!valid" class="ma-3" @click.stop="submit">
        作成する
      </v-btn>
    </v-form>
    <!-- ここ作成したらokを押しても押さなくても自動遷移 -->
    <v-snackbar v-model="snackbar">
      作成できました
      <v-btn
        :to="`/applications/` + response.application_id"
        color="green darken-1"
        text
        @click="snackbar = false"
      >
        OK
      </v-btn>
    </v-snackbar>
  </v-container>
</template>

<script setup lang="ts">
import { useMeStore } from "@/stores/me";
import { useUserListStore } from "@/stores/userList";
import { applicationType, remarksTitle } from "@/use/applicationDetail";
import { dayPrint } from "@/use/dataFormat";
import { remarksHint } from "@/use/inputFormText";
import Icon from "@/views/shared/Icon.vue";
import ImageUploader from "@/views/shared/ImageUploader.vue";
import axios from "axios";
import { storeToRefs } from "pinia";
import { computed, onMounted, reactive, ref } from "vue";
import { useRoute } from "vue-router";

const route = useRoute();
const meStore = useMeStore();
const userListStore = useUserListStore();

const { trapId } = storeToRefs(meStore);
const { userList } = storeToRefs(userListStore);
const { fetchUserList } = userListStore;

const form = ref<any>(null); // eslint-disable-line @typescript-eslint/no-explicit-any
const firstfocus = ref<any>(null); // eslint-disable-line @typescript-eslint/no-explicit-any

const response = reactive({
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
});

const snackbar = ref(false);
const date = ref(null);
const menu = ref(false);
const traPID = ref<string[]>([]);
const valid = ref(true);
const title = ref("");
const amount = ref("");
const remarks = ref("");
const imageBlobs = ref<File[]>([]);

const amountRules = [
  (v: unknown) => !!v || "必須の項目です",
  (v: unknown) => !!String(v).match("^[1-9][0-9]*$") || "金額が不正です"
];
const nullRules = [(v: unknown) => !!v || "必須の項目です"];

const traPIDs = computed(() => {
  return userList.value.map(user => user.trap_id);
});

const computedDateFormatted = computed(() => {
  return formatDate(date.value);
});

const formatDate = (date: string | number | Date | null) => {
  if (!date) return null;
  const d = new Date(date);
  return `${d.getFullYear()}年${d.getMonth() + 1}月${d.getDate()}日`;
};

const returnToday = () => {
  const date = new Date();
  return dayPrint(date);
};

const returnType = (type: string) => {
  return applicationType(type);
};

const returnRemarksTitle = (type: string) => {
  return remarksTitle(type);
};

const returnRemarksHint = (type: string) => {
  return remarksHint(type);
};

const submit = () => {
  if (form.value.validate()) {
    const formData = new FormData();
    const paid_at = new Date(date.value || Date.now());
    const details = {
      type: route.params.type,
      title: title.value,
      remarks: remarks.value,
      paid_at: paid_at.toISOString(),
      amount: Number(amount.value),
      repaid_to_id: traPID.value
    };
    formData.append("details", JSON.stringify(details));
    imageBlobs.value.forEach(imageBlob => {
      formData.append("images", imageBlob);
    });
    axios
      .post("/api/applications", formData, {
        headers: { "content-type": "multipart/form-data" }
      })
      .then(res => {
        Object.assign(response, res.data);
        snackbar.value = true;
      });
  }
};

onMounted(async () => {
  await fetchUserList();
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  firstfocus.value.focus();
});
</script>
