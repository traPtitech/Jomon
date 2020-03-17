import axios from "axios";
import Vue from "vue";

export const applicationList = {
  state: [
    {
      application_id: "",
      created_at: "",
      applicant: {
        trap_id: "",
        is_admin: false
      },
      current_detail: {
        update_user: {
          trap_id: "",
          is_admin: false
        },
        type: "",
        title: "",
        remarks: "",
        amount: 0,
        paid_at: "",
        updated_at: ""
      },
      current_state: ""
    }
  ],
  mutations: {
    setApplicationList(state, newState) {
      state.splice(0);
      newState.forEach((element, index) => {
        Vue.set(state, index, element);
      });
    }
  },
  actions: {
    async getApplicationList({ commit }) {
      try {
        const response = await axios.get("/api/applications");
        commit("setApplicationList", response.data);
      } catch (err) {
        console.log(err);
      }
    }
  }
};
