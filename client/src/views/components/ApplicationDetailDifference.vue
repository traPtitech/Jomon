<template>
  <div :class="$style.container">
    <h3>{{ changedItem }}</h3>
    <div>
      <div :class="$style.before">
        <div :class="$style.mark">−</div>
        <!-- eslint-disable-next-line vue/no-v-html -->
        <div v-html="renderedBefore" />
      </div>
      <div :class="$style.after">
        <div :class="$style.mark">+</div>
        <!-- eslint-disable-next-line vue/no-v-html -->
        <div v-html="renderedAfter" />
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { applicationType } from "@/use/applicationDetail";
import { dayPrint, numberFormat } from "@/use/dataFormat";
import { render } from "@/use/markdown";
import { onMounted, ref } from "vue";

const props = withDefaults(
  defineProps<{
    item?: string;
    pre?: string | number | null;
    now?: string | number | null;
  }>(),
  {
    item: "",
    pre: null,
    now: null
  }
);

const changedItem = ref("");
const renderedBefore = ref("");
const renderedAfter = ref("");

onMounted(async () => {
  switch (props.item) {
    case "type":
      changedItem.value = "申請様式";
      break;
    case "title":
      changedItem.value = "タイトル";
      break;
    case "remarks":
      changedItem.value = "詳細";
      break;
    case "amount":
      changedItem.value = "金額";
      break;
    case "paid_at":
      changedItem.value = "支払日";
      break;
    default:
      changedItem.value = "";
  }

  const preRender = (item: string, data: string | number | null) => {
    if (data === null) return "error";
    switch (item) {
      case "type":
        return applicationType(String(data)) + "申請";
      case "title":
        return String(data);
      case "remarks":
        return String(data);
      case "amount":
        return numberFormat(Number(data)) + "円";
      case "paid_at":
        return dayPrint(String(data));
      default:
        return "error";
    }
  };

  renderedBefore.value = await render(preRender(props.item, props.pre));
  renderedAfter.value = await render(preRender(props.item, props.now));
});
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
