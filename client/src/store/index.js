import Vue from "vue";
import Vuex from "vuex";
import { me } from "./me";
import { applicationList } from "./list";
import { detail } from "./detail";
import { logs } from "./logs";

Vue.use(Vuex);
export default new Vuex.Store({
  state: {},
  mutations: {},
  actions: {},
  modules: {
    me: me,
    applicationList: applicationList,
    application_detail_paper: detail,
    logs: logs
  }
});
