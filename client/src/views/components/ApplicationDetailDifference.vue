<template>
  <div class="difference">
    <v-row class="pa-5">
      <v-col cols="2" class="pl-0">
        <v-row>
          <v-col>{{ itemPrint(this.item) }}</v-col>
        </v-row>
      </v-col>
      <v-col cols="10" class="pa-2">
        <v-row :style="{ background: red }">
          <v-col cols="1">-</v-col>
          <v-col cols="11"> {{ dataPrint(this.item, this.pre) }}</v-col>
        </v-row>
        <v-row :style="{ background: green }">
          <v-col cols="1">+ </v-col>
          <v-col cols="11"> {{ dataPrint(this.item, this.now) }}</v-col>
        </v-row>
      </v-col>
    </v-row>
  </div>
</template>
<script>
import { applicationType } from "@/use/applicationDetail";
import { numberFormat, dayPrint } from "@/use/dataFormat";

export default {
  name: "ApplicationDetailDifference",
  props: {
    item: String,
    pre: null,
    now: null
  },
  data: () => {
    return { red: "#ffeef0", green: "#e6ffed" };
  },
  methods: {
    numberFormat(price) {
      return numberFormat(price);
    },
    dayPrint(time) {
      return dayPrint(time);
    },
    itemPrint(item) {
      switch (item) {
        case "type":
          return "申請様式";
        case "title":
          return "タイトル";
        case "remarks":
          return "詳細";
        case "amount":
          return "金額";
        case "paid_at":
          return "支払日";
      }
    },
    dataPrint(item, data) {
      switch (item) {
        case "type":
          return applicationType(data) + "申請";
        case "title":
          return data;
        case "remarks":
          return data;
        case "amount":
          return this.numberFormat(data) + "円";
        case "paid_at":
          return this.dayPrint(data);
        default:
          return "error";
      }
    }
  }
};
</script>
