import { defineStore } from "pinia";
import { ref } from "vue";

export const useToastStore = defineStore("toast", () => {
  const isVisible = ref(false);
  const message = ref("");
  const color = ref("success");

  const show = (msg: string, type: "success" | "error" = "success") => {
    message.value = msg;
    color.value = type;
    isVisible.value = true;
  };

  return {
    isVisible,
    message,
    color,
    show
  };
});
