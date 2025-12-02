<template>
  <div :class="$style.container">
    <v-form ref="formRef" v-model="valid" lazy-validation>
      <div>
        <div :class="$style.header">
          <h1 :class="$style.title">
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
              :class="$style.selector"
            />
            申請
          </h1>

          <div>
            <div>申請日: {{ returnDate(detailCore.created_at) }}</div>
            <div>
              申請者:
              <Icon :user="detailCore.applicant.trap_id" :size="20" />
              {{ detailCore.applicant.trap_id }}
            </div>
          </div>
        </div>

        <v-text-field
          v-model="title_change"
          :rules="nullRules"
          label="概要"
          filled
          :placeholder="returnTitlePlaceholder(type_object.type)"
        />

        <v-menu
          v-model="menu"
          :close-on-content-click="false"
          transition="scale-transition"
          offset-y
        >
          <template #activator="{ props }">
            <v-text-field
              v-model="computedDateFormatted"
              :rules="nullRules"
              label="支払日"
              filled
              readonly
              placeholder="2020年5月2日"
              v-bind="props"
            />
          </template>
          <v-date-picker
            v-model="paid_at_change"
            no-title
            @input="menu = false"
          />
        </v-menu>

        <v-text-field
          v-model="amount_change"
          :rules="amountRules"
          type="number"
          label="支払金額"
          placeholder="100"
          suffix="円"
        />

        <v-autocomplete
          ref="traPID"
          v-model="repaid_to_id_change"
          :rules="[
            () =>
              repaid_to_id_change.length > 0 || '返金対象者は一人以上必要です'
          ]"
          label="返金対象者"
          filled
          :items="traPIDs"
          placeholder="traQIDs"
          hint="traQ IDの一部入力で候補が表示されます"
          required
          multiple
        />

        <v-textarea
          v-model="remarks_change"
          :rules="nullRules"
          filled
          :label="returnRemarksTitle(type_object.type)"
          :placeholder="returnRemarksPlaceholder(type_object.type)"
          :hint="returnRemarksHint(type_object.type)"
          auto-grow
        />

        <div>
          <h3>画像</h3>
          <div :class="$style.image_container">
            <div v-for="(path, index) in detailCore.images" :key="path">
              <div v-if="images_change[index]">
                <img :src="`/api/images/${path}`" />
                <v-btn
                  rounded
                  color="primary"
                  name="delete"
                  @click="deleteImage(index)"
                >
                  delete
                </v-btn>
              </div>
              <div v-else>
                <v-btn
                  rounded
                  color="primary"
                  name="cancel"
                  @click="cancelDeleteImage(index)"
                >
                  cancel
                </v-btn>
              </div>
            </div>
          </div>
          <h3>画像を追加</h3>
          <image-uploader v-model="imageBlobs" />
        </div>
      </div>

      <div class="$style.button_container">
        <simple-button :label="`修正する`" @click.stop="submit" />
        <simple-button :label="`取り消す`" @click="deleteFix" />
      </div>
    </v-form>
    <!-- ここ作成したらokを押しても押さなくても自動遷移 -->
    <v-snackbar v-model="snackbar">
      変更できました
      <v-btn
        :to="`/applications/` + response.application_id"
        color="green darken-1"
        text
        @click="afterChange()"
      >
        OK
      </v-btn>
    </v-snackbar>
  </div>
</template>

<script setup lang="ts">
import { useApplicationDetailStore } from "@/stores/applicationDetail";
import { useUserListStore } from "@/stores/userList";
import { remarksTitle } from "@/use/applicationDetail";
import { dayPrint } from "@/use/dataFormat";
import {
  remarksHint,
  remarksPlaceholder,
  titlePlaceholder
} from "@/use/inputFormText";
import Icon from "@/views/shared/Icon.vue";
import ImageUploader from "@/views/shared/ImageUploader.vue";
import SimpleButton from "@/views/shared/SimpleButton.vue";
import axios from "axios";
import { storeToRefs } from "pinia";
import { computed, onMounted, reactive, ref } from "vue";
import { useRoute } from "vue-router";

const route = useRoute();
const applicationDetailStore = useApplicationDetailStore();
const userListStore = useUserListStore();

const { core: detailCore } = storeToRefs(applicationDetailStore);
const { userList } = storeToRefs(userListStore);
const { fetchApplicationDetail, deleteFix } = applicationDetailStore;
const { fetchUserList } = userListStore;

const formRef = ref<any>(null); // eslint-disable-line @typescript-eslint/no-explicit-any

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

const type_object = reactive({ jpn: "", type: "" });
const types = [
  { jpn: "部費利用", type: "club" },
  { jpn: "大会等旅費補助", type: "contest" },
  { jpn: "イベント交通費補助", type: "event" },
  { jpn: "渉外交通費補助", type: "public" }
];

const snackbar = ref(false);
const menu = ref(false);
const valid = ref(true);

const amountRules = [
  (v: unknown) => !!v || "",
  (v: unknown) => !!String(v).match("^[1-9][0-9]*$") || "金額が不正です"
];
const nullRules = [(v: unknown) => !!v || ""];

const title_change = ref("");
const remarks_change = ref("");
const paid_at_change = ref("");
const amount_change = ref("");
const repaid_to_id_change = ref<string[]>([]);
const imageBlobs = ref<File[]>([]);
const images_change = ref<boolean[]>([]);

// const changeRules = [
//   (v: unknown) => (v !== detailCore.value.repayment_logs && !!v) || ""
// ];

const formatDate = (date: string) => {
  if (!date) return null;
  const [year, month, day] = date.split("-");
  return `${year}年${month.replace(/^0/, "")}月${day.replace(/^0/, "")}日`;
};

const computedDateFormatted = computed(() => formatDate(paid_at_change.value));

const traPIDs = computed(() => {
  return userList.value.map(user => user.trap_id);
});

const returnDate = (date: string) => dayPrint(date);
const returnRemarksTitle = (type: string) => remarksTitle(type);
const returnTitlePlaceholder = (type: string) => titlePlaceholder(type);
const returnRemarksPlaceholder = (type: string) => remarksPlaceholder(type);
const returnRemarksHint = (type: string) => remarksHint(type);

const deleteImage = (index: number) => {
  images_change.value[index] = false;
};

const cancelDeleteImage = (index: number) => {
  images_change.value[index] = true;
};

const afterChange = () => {
  snackbar.value = false;
  deleteFix();
};

const submit = async () => {
  if (formRef.value.validate()) {
    images_change.value.forEach((flag, index) => {
      if (!flag) {
        axios.delete("/api/images/" + detailCore.value.images[index]);
      }
    });
    const form = new FormData();
    const date = new Date(paid_at_change.value);
    const details = {
      type: type_object.type,
      title: title_change.value,
      remarks: remarks_change.value,
      paid_at: date.toISOString(),
      amount: Number(amount_change.value),
      repaid_to_id: repaid_to_id_change.value
    };
    form.append("details", JSON.stringify(details));
    imageBlobs.value.forEach(imageBlob => {
      form.append("images", imageBlob);
    });
    const res = await axios.patch(
      "/api/applications/" + detailCore.value.application_id,
      form,
      {
        headers: { "content-type": "multipart/form-data" }
      }
    );
    Object.assign(response, res.data);
    await fetchApplicationDetail(route.params.id as string);
    snackbar.value = true;
  }
};

onMounted(async () => {
  title_change.value = detailCore.value.current_detail.title;
  type_object.type = detailCore.value.current_detail.type;
  // title_change.value = detailCore.value.current_detail.title; // Duplicate in original
  remarks_change.value = detailCore.value.current_detail.remarks;
  paid_at_change.value = detailCore.value.current_detail.paid_at;
  amount_change.value = String(detailCore.value.current_detail.amount);
  images_change.value = new Array(detailCore.value.images.length).fill(true);

  await fetchUserList();
  const trap_ids = detailCore.value.repayment_logs.map(
    log => log.repaid_to_user.trap_id
  );
  repaid_to_id_change.value = trap_ids;
});
</script>

<style lang="scss" module>
.container {
  height: fit-content;
  margin: 12px;
  padding: 12px;
  box-shadow: 0 3px 1px -2px rgb(0 0 0 / 20%), 0 2px 2px 0 rgb(0 0 0 / 14%),
    0 1px 5px 0 rgb(0 0 0 / 12%);
}
.header {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
}
.title {
  display: flex;
  flex: 1;
}
.selector {
  flex: 0;
}
.section {
  margin: 16px 0;
  border-bottom: 1px solid $color-grey;
}
.section_title {
  color: $color-text-primary-disabled;
}
.section_item {
  margin-left: 8px;
}
.image_container {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(min(100%, 360px), 1fr));
  gap: 16px;
  padding: 8px;
}
</style>
