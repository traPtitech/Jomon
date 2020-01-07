<template>
  <div class="application-list">
    <v-card width="1200px" class="mx-auto mt-5" color="primary">
      <v-card-title style="color:white;font-weight:bold;font-size:1.5em"
        >申請書一覧</v-card-title
      >
      <v-card-text class="pl-0 pr-0 pb-0">
        <v-list>
          <v-list-item-group color="primary">
            <Application :list="header" class="pb-0 pt-0"></Application>
            <v-list-item
              v-for="list in applicationList"
              v-bind:key="list.id"
              :to="'/applications/' + list.id"
              class="pl-0 pr-0"
            >
              <Application :list="list"> </Application>
            </v-list-item>
          </v-list-item-group>
        </v-list>
      </v-card-text>
    </v-card>
  </div>
</template>
<script>
import { mapState } from "vuex";
import Application from "./components/Application";
export default {
  name: "ApplicationList",
  computed: {
    ...mapState(["applicationList"])
  },
  data() {
    return {
      header: {
        title: "タイトル",
        name: "申請者名",
        money: "金額",
        state: "状態"
      }
    };
  },
  components: {
    Application
  }
};
</script>
