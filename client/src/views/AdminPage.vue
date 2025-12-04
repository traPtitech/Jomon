<template>
  <div v-if="isAdmin" class="admin pa-4">
    <v-chip v-for="(user, index) in adminList" :key="index">
      {{ user }}
    </v-chip>
    <v-autocomplete
      v-model="removeAdminUsers"
      :items="adminList"
      label="管理権限を削除"
      multiple
    />
    <simple-button
      :label="'設定'"
      :variant="'warning'"
      @click="removeAdmin()"
    />
    <v-autocomplete
      v-model="addAdminUsers"
      :items="notAdminList"
      label="管理権限を追加"
      multiple
    />
    <simple-button :label="'設定'" :variant="'warning'" @click="addAdmin()" />
  </div>
  <div v-else>権限がありません</div>
</template>

<script setup lang="ts">
import { useMeStore } from "@/stores/me";
import { useUserListStore } from "@/stores/userList";
import SimpleButton from "@/views/shared/SimpleButton.vue";
import axios from "axios";
import { storeToRefs } from "pinia";
import { onMounted, ref } from "vue";

const meStore = useMeStore();
const { isAdmin } = storeToRefs(meStore);

const userListStore = useUserListStore();
const { adminList, notAdminList } = storeToRefs(userListStore);
const { fetchUserList } = userListStore;

const addAdminUsers = ref<string[]>([]);
const removeAdminUsers = ref<string[]>([]);

const addAdmin = async () => {
  if (addAdminUsers.value.length > 0) {
    await Promise.all(
      addAdminUsers.value.map(user =>
        axios
          .put("api/users/admins", { trap_id: user, to_admin: true })
          .catch(e => alert(e))
      )
    );
    fetchUserList();
    addAdminUsers.value = [];
  }
};

const removeAdmin = async () => {
  if (removeAdminUsers.value.length > 0) {
    await Promise.all(
      removeAdminUsers.value.map(user =>
        axios
          .put("api/users/admins", { trap_id: user, to_admin: false })
          .catch(e => alert(e))
      )
    );
    fetchUserList();
    removeAdminUsers.value = [];
  }
};

onMounted(() => {
  fetchUserList();
});
</script>
