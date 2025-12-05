import axios from "axios";
import { defineStore } from "pinia";
import { computed, ref } from "vue";

interface User {
  trap_id: string;
  is_admin: boolean;
}

export const useUserListStore = defineStore("userList", () => {
  const userList = ref<User[]>([]);

  const trapIds = computed(() => userList.value.map(user => user.trap_id));

  const adminList = computed(() =>
    userList.value.filter(user => user.is_admin).map(user => user.trap_id)
  );

  const notAdminList = computed(() =>
    userList.value.filter(user => !user.is_admin).map(user => user.trap_id)
  );

  const fetchUserList = async () => {
    const response = await axios.get("/api/users");
    userList.value = response.data;
  };

  return {
    userList,
    trapIds,
    adminList,
    notAdminList,
    fetchUserList
  };
});
