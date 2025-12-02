import axios from "axios";

export const userList = {
  state: [
    {
      trap_id: "",
      is_admin: false
    }
  ],
  mutations: {
    setUserList(state, newState) {
      state.splice(0);
      state.push(...newState);
    }
  },
  actions: {
    async getUserList({ commit }) {
      const response = await axios.get("/api/users");
      commit("setUserList", response.data);
    }
  },
  getters: {
    trap_ids: state => {
      return state.map(data => data.trap_id);
    },
    adminList: state => {
      let admin = [];
      state.forEach(user => {
        if (user.is_admin) admin.push(user.trap_id);
      });
      return admin;
    },
    notAdminList: state => {
      let notAdmin = [];
      state.forEach(user => {
        if (!user.is_admin) notAdmin.push(user.trap_id);
      });
      return notAdmin;
    }
  }
};
