import axios from "axios";
export const applicationDetail = {
  state: {
    application_id: "",
    created_at: "",
    applicant: {
      trap_id: "",
      is_admin: false
    },
    current_detail: {
      update_user: {
        trap_id: "",
        is_admin: ""
      },
      type: "",
      title: "",
      remarks: "",
      amount: 0,
      paid_at: "",
      updated_at: "",
      repaid_to_id: []
    },
    current_state: "",
    images: [""],
    comments: [
      {
        comment_id: 0,
        user: {
          trap_id: "",
          is_admin: false
        },
        comment: "",
        created_at: "",
        updated_at: ""
      }
    ],
    state_logs: [
      {
        update_user: {
          trap_id: "",
          is_admin: false
        },
        to_state: "",
        reason: "",
        created_at: ""
      }
    ],
    application_detail_logs: [
      {
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
      }
    ],
    repayment_logs: [
      {
        repaid_by_user: {
          trap_id: "",
          is_admin: false
        },
        repaid_to_user: {
          trap_id: "",
          is_admin: false
        },
        repaid_at: ""
      }
    ]
  },
  getters: {
    logs: state => {
      let logs = [];
      state.comments.forEach(log => {
        logs.push({
          log_type: "comment",
          content: log,
          sort_date: new Date(log.created_at)
        });
      });
      state.state_logs.forEach(log => {
        logs.push({
          log_type: "state",
          content: log,
          sort_date: new Date(log.created_at)
        });
      });
      let pre_record = "";
      state.application_detail_logs.forEach(log => {
        if (pre_record !== "") {
          logs.push({
            log_type: "application",
            content: {
              log: log,
              pre_log: pre_record
            },
            sort_date: new Date(log.updated_at)
          });
        }
        pre_record = log;
      });
      state.repayment_logs.forEach(log => {
        logs.push({
          log_type: "repayment",
          content: log,
          sort_date: new Date(log.repaid_at)
        });
      });
      logs.sort((a, b) => {
        if (a.sort_date > b.sort_date) {
          return 1;
        } else {
          return -1;
        }
      });
      return logs;
    }
  },
  mutations: {
    setApplicationDetail(state, newState) {
      Object.keys(state).forEach(key => {
        state[key] = newState[key];
      });
    }
  },
  actions: {
    async getApplicationDetail({ commit }, applicationId) {
      try {
        const response = await axios.get("/api/applications/" + applicationId);
        commit("setApplicationDetail", response.data);
      } catch (err) {
        console.log(err);
      }
    }
  }
};
