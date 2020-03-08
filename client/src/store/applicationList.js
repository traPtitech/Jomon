import axios from "axios";
import Vue from "vue";

export const State = {
  submitted: "submitted",
  rejected: "rejected",
  fix_required: "fix_required",
  accepted: "accepted",
  fully_repaid: "fully_repaid"
};

export const applicationList = {
  state: [
    {
      application_id: "3fa85f64-5717-4562-b3fc-2c963f66afa6",
      created_at: "2020-03-08T08:20:44.788Z",
      applicant: {
        trap_id: "nagatech",
        is_admin: true
      },
      current_detail: {
        update_user: {
          trap_id: "nagatech",
          is_admin: true
        },
        type: "club",
        title: "夏コミの交通費をお願いします。",
        remarks: "〇〇駅から〇〇駅への移動",
        amount: 1200,
        paid_at: "2020-03-08",
        updated_at: "2020-03-08T08:20:44.788Z"
      },
      current_state: State.submitted
    }
  ],
  mutations: {
    set(state, newState) {
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
        commit("set", response.data);
      } catch (err) {
        console.log(err);
      }
    }
  }
};
