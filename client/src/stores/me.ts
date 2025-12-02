import axios from "axios";
import { defineStore } from "pinia";
import { ref } from "vue";

export const useMeStore = defineStore("me", () => {
  const trapId = ref<string>("");
  const isAdmin = ref<boolean>(false);

  const fetchMe = async () => {
    const response = await axios.get("/api/users/me");
    trapId.value = response.data.trap_id;
    isAdmin.value = response.data.is_admin;
  };

  return {
    trapId,
    isAdmin,
    fetchMe
  };
});
