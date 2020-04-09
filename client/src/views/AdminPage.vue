<template>
  <div class="admin">
    <v-chip v-for="(user, index) in adminList" :key="index">{{ user }}</v-chip>
    <v-autocomplete
      v-model="removeAdminUsers"
      :items="adminList"
      label="管理権限を削除"
      multiple
    ></v-autocomplete>
    <v-btn @click="removeAdmin()">設定</v-btn>
    <v-autocomplete
      v-model="addAdminUsers"
      :items="notAdminList"
      label="管理権限を追加"
      multiple
    ></v-autocomplete>
    <v-btn @click="addAdmin()">設定</v-btn>
  </div>
</template>

<script>
import { mapGetters, mapActions } from "vuex";
import axios from "axios";
export default {
  computed: {
    ...mapGetters(["adminList", "notAdminList"])
  },
  methods: {
    ...mapActions(["getUserList"]),
    addAdmin() {
      if (this.addAdminUsers.length > 0) {
        this.addAdminUsers.forEach(user => {
          axios
            .put("api/users/admins", { trap_id: user, to_admin: true })
            .catch(e => console.log(e));
        });
      }
    },
    removeAdmin() {
      if (this.removeAdminUsers.length > 0) {
        this.removeAdminUsers.forEach(user => {
          axios
            .put("api/users/admins", { trap_id: user, to_admin: false })
            .catch(e => console.log(e));
        });
      }
    }
  },
  data() {
    return {
      addAdminUsers: [],
      removeAdminUsers: []
    };
  },
  created() {
    this.getUserList();
  }
};
</script>
