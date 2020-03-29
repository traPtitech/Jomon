<template>
  <div class="admin">
    <P v-for="(user, index) in adminList" :key="index">{{ user }}</P>
    <v-autocomplete
      :items="adminList"
      label="管理権限を削除"
      multiple
    ></v-autocomplete>
    <v-autocomplete
      :items="notAdminList"
      label="管理権限を削追加"
      multiple
    ></v-autocomplete>
    <v-btn>設定</v-btn>
  </div>
</template>

<script>
import { mapGetters, mapActions } from "vuex";
export default {
  computed: {
    ...mapGetters(["adminList", "notAdminList"])
  },
  methods: {
    ...mapActions(["getUserList"])
  },
  created() {
    this.getUserList();
  }
};
</script>
