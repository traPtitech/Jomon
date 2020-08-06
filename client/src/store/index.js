import Vue from "vue";
import Vuex from "vuex";
import { me } from "./me";
import { applicationList } from "./applicationList";
import { applicationDetail } from "./applicationDetail";
import { userList } from "./userList";

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
