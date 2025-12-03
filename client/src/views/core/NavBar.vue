<template>
  <v-app-bar color="#011A27" theme="dark">
    <router-link to="/">
      <v-img
        alt="Jomon Logo"
        class="shrink-mr-2"
        :src="logo"
        transition="scale-transition"
        width="144"
        style="display: inline-block; cursor: pointer"
      />
    </router-link>
    <v-spacer />

    <v-btn v-if="isAdmin" to="/admin" variant="text"> 管理ページ </v-btn>
    <v-menu location="bottom">
      <template #activator="{ props }">
        <v-btn v-bind="props" variant="text"> 新規作成 </v-btn>
      </template>

      <v-list>
        <v-list-item
          v-for="(item, index) in items"
          :key="index"
          :to="'/applications/new/' + item.title"
        >
          <v-list-item-title>{{ item.page }}</v-list-item-title>
        </v-list-item>
      </v-list>
    </v-menu>
    <Icon :user="trapId" :size="35" />
  </v-app-bar>
</template>
<script lang="ts">
import logo from "@/assets/logo.svg";
import { useMeStore } from "@/stores/me";
import Icon from "@/views/shared/Icon.vue";
import { mapState } from "pinia";

export default {
  components: {
    Icon
  },
  data() {
    return {
      logo,
      items: [
        { title: "club", page: "部費利用申請" },
        { title: "contest", page: "大会等旅費補助申請" },
        { title: "event", page: "イベント交通費補助申請" },
        { title: "public", page: "渉外交通費補助" }
      ]
    };
  },
  computed: {
    ...mapState(useMeStore, ["trapId", "isAdmin"])
  }
};
</script>
