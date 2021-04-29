<template>
  <div :class="$style.container">
    <h3>{{ changedItem }}</h3>
    <div>
      <div :class="$style.before">
        <div :class="$style.mark">−</div>
        <div v-html="renderedBefore" />
      </div>
      <div :class="$style.after">
        <div :class="$style.mark">+</div>
        <div v-html="renderedAfter" />
      </div>
    </div>
  </div>
</template>
<script>
import { applicationType } from "@/use/applicationDetail";
import { numberFormat, dayPrint } from "@/use/dataFormat";
import { render } from "@/use/markdown";

export default {
  name: "ApplicationDetailDifference",
  props: {
    item: String,
    pre: null,
    now: null
  },
  data: () => {
    return { changedItem: "", renderedBefore: "", renderedAfter: "" };
  },
  async mounted() {
    this.changedItem = (item => {
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
    })(this.item);
    const preRender = (item, data) => {
      switch (item) {
        case "type":
          return applicationType(data) + "申請";
        case "title":
          return data;
        case "remarks":
          return data;
        case "amount":
          return numberFormat(data) + "円";
        case "paid_at":
          return dayPrint(data);
        default:
          return "error";
      }
    };
    this.renderedBefore = await render(preRender(this.item, this.pre));
    this.renderedAfter = await render(preRender(this.item, this.now));
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
