import axios from "axios";

export const me = {
  state: {
    trap_id: "",
    is_admin: false
  },
  mutations: {
    setMe(state, newState) {
      Object.keys(state).forEach(key => {
        state[key] = newState[key];
      });
    }
  },
  actions: {
    async getMe(commit) {
      try {
        const response = await axios.get("/api/users/me");
        commit("setMe", response.data);
      } catch (err) {
        console.log(err);
        if (err.response.status === 401) {
          try {
            commit("toggleLoginDialog");
          } catch (err) {
            console.log(err);
          }
        }
      }
    }
  }
};
