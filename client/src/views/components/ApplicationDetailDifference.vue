<template>
  <div :class="$style.container">
    <h3>{{ itemPrint(this.item) }}</h3>
    <div>
      <div :class="$style.before">
        <div :class="$style.mark">−</div>
        {{ dataPrint(this.item, this.pre) }}
      </div>
      <div :class="$style.after">
        <div :class="$style.mark">+</div>
        {{ dataPrint(this.item, this.now) }}
      </div>
    </div>
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

<style lang="scss" module>
.container {
  padding: 8px;
}
.before {
  background: #ffeef0;
  display: flex;
  padding: 8px;
}
.after {
  background: #e6ffed;
  display: flex;
  padding: 8px;
}
.mark {
  width: 8px;
  margin-right: 8px;
  text-align: center;
}
</style>
