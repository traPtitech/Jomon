export const token = {
  state: {
    loginDialog: false
  },
  mutations: {
    toggleLoginDialog(state) {
      state.loginDialog = !state.loginDialog;
    }
  },
  actions: {}
};
