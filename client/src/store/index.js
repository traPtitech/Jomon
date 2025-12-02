import { applicationDetail } from "@/store/applicationDetail";
import { applicationList } from "@/store/applicationList";
import { me } from "@/store/me";
import { userList } from "@/store/userList";
import { createStore } from "vuex";

export default createStore({
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
