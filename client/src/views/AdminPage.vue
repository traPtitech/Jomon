<template>
  <div v-if="this.$store.state.me.is_admin" class="admin">
    <v-chip v-for="(user, index) in adminList" :key="index">{{ user }}</v-chip>
    <v-autocomplete
      v-model="removeAdminUsers"
      :items="adminList"
      label="管理権限を削除"
      multiple
    ></v-autocomplete>
    <simple-button
      :label="'設定'"
      :variant="'warning'"
      @click="removeAdmin()"
    ></simple-button>
    <v-autocomplete
      v-model="addAdminUsers"
      :items="notAdminList"
      label="管理権限を追加"
      multiple
    ></v-autocomplete>
    <simple-button
      :label="'設定'"
      :variant="'warning'"
      @click="addAdmin()"
    ></simple-button>
  </div>
  <div v-else>権限がありません</div>
</template>

<script>
import { mapGetters, mapActions } from "vuex";
import axios from "axios";
import SimpleButton from "@/views/shared/SimpleButton";

export default {
  components: {
    SimpleButton
  },
  computed: {
    ...mapGetters(["adminList", "notAdminList"])
  },
  methods: {
    ...mapActions(["getUserList"]),
    async addAdmin() {
      if (this.addAdminUsers.length > 0) {
        await this.addAdminUsers.forEach(user => {
          axios
            .put("api/users/admins", { trap_id: user, to_admin: true })
            .catch(e => alert(e));
        });
        this.getUserList();
      }
    },
    async removeAdmin() {
      if (this.removeAdminUsers.length > 0) {
        await this.removeAdminUsers.forEach(user => {
          axios
            .put("api/users/admins", { trap_id: user, to_admin: false })
            .catch(e => alert(e));
        });
      }
      this.getUserList();
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
