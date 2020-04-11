export const token = {
  state: {
    // loginDialog: false
    loginDialog: true
  },
  mutations: {
    toggleLoginDialog(state) {
      state.loginDialog = !state.loginDialog;
    }
  },
  actions: {}
};
