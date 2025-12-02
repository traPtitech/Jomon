<template>
  <div :class="$style.container">
    <v-form ref="form" v-model="valid" lazy-validation>
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
            <div>申請日: {{ returnDate(detail.core.created_at) }}</div>
            <div>
              申請者:
              <Icon :user="detail.core.applicant.trap_id" :size="20" />
              {{ detail.core.applicant.trap_id }}
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
          <template #activator="{ on }">
            <v-text-field
              v-model="computedDateFormatted"
              :rules="nullRules"
              label="支払日"
              filled
              readonly
              placeholder="2020年5月2日"
              v-on="on"
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
            () => !(repaid_to_id_change === 0) || '返金対象者は一人以上必要です'
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
            <div v-for="(path, index) in detail.core.images" :key="path">
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

<script>
import { remarksTitle } from "@/use/applicationDetail";
import { dayPrint } from "@/use/dataFormat";
import {
  remarksHint,
  remarksPlaceholder,
  titlePlaceholder
} from "@/use/inputFormText";
import Icon from "@/views/shared/Icon";
import ImageUploader from "@/views/shared/ImageUploader";
import SimpleButton from "@/views/shared/SimpleButton";
import axios from "axios";
import { mapActions, mapMutations, mapState } from "vuex";

export default {
  components: {
    Icon,
    ImageUploader,
    SimpleButton
  },
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
  }
};
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
