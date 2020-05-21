<template>
  <v-app-bar app color="primary" dark>
    <div class="d-flex align-center">
      <v-img
        alt="Vuetify Logo"
        class="shrink-mr-2"
        contain
        src="https://cdn.vuetifyjs.com/images/logos/vuetify-logo-dark.png"
        transition="scale-transition"
        width="40"
      />

      <router-link to="/">
        <v-toolbar-title class="white--text">Jomon</v-toolbar-title>
      </router-link>
    </div>
    <v-spacer></v-spacer>

    <v-btn v-if="me.is_admin" to="/admin" text>
      管理ページ
    </v-btn>
    <v-menu open-on-hover offset-y>
      <template v-slot:activator="{ on }">
        <v-btn color="primary" dark v-on="on">
          新規作成
        </v-btn>
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
    <Icon :user="this.$store.state.me.trap_id" :size="35" />
  </v-app-bar>
</template>
<script>
import { mapState } from "vuex";
import Icon from "../components/Icon";
export default {
  components: {
    Icon
  },
  computed: {
    ...mapState(["me"])
  },
  data: () => ({
    items: [
      { title: "club", page: "部費利用申請" },
      { title: "contest", page: "大会等旅費補助申請" },
      { title: "event", page: "イベント交通費補助申請" },
      { title: "public", page: "渉外交通費補助" }
    ]
  })
};
</script>
