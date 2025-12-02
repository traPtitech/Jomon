<template>
  <div v-if="$store.state.me.is_admin" class="admin">
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

<script>
import SimpleButton from "@/views/shared/SimpleButton.vue";
import axios from "axios";
import { mapActions, mapGetters } from "vuex";

export default {
  components: {
    SimpleButton
  },
  data() {
    return {
      addAdminUsers: [],
      removeAdminUsers: []
    };
  },
  computed: {
    ...mapGetters(["adminList", "notAdminList"])
  },
  created() {
    this.getUserList();
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
  }
};
</script>
