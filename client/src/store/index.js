import Vue from "vue";
import Vuex from "vuex";
import { me } from "@/store/me";
import { applicationList } from "@/store/applicationList";
import { applicationDetail } from "@/store/applicationDetail";
import { userList } from "@/store/userList";

Vue.use(Vuex);
export default new Vuex.Store({
  state: {},
  mutations: {},
  actions: {},
  modules: {
    me: me,
    applicationList: applicationList,
    application_detail_paper: applicationDetail,
    userList: userList
  }
});
