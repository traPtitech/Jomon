import axios from "axios";

export const applicationList = {
  state: [],
  mutations: {
    setApplicationList(state, newState) {
      state.splice(0);
      state.push(...newState);
    }
  },
  actions: {
    async getApplicationList({ commit }, params) {
      Object.keys(params).forEach(key => {
        if (
          params[key] === null ||
          params[key] === undefined ||
          params[key] === ""
        ) {
          delete params[key];
        }
      });
      const response = await axios.get("/api/applications", {
        params
      });
      commit("setApplicationList", response.data);
    }
  }
};
