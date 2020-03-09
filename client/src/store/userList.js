//import axios from "axios";
export const userList = {
  state: [
    {
      trap_id: "user1",
      is_admin: true
    },
    {
      trap_id: "user2",
      is_admin: false
    },
    {
      trap_id: "user3",
      is_admin: false
    },
    {
      trap_id: "user4",
      is_admin: false
    }
  ],
  // actions: {
  //   getUsers: function({ commit }) {
  //     return axios.get("https://q.trap.jp/api/1.0/users").then(response => {
  //       commit("setUsers", response.data);
  //     });
  //   }
  // },
  // mutations: {
  //   setUsers: function(userList, users) {
  //     userList.state = users;
  //   }
  // },
  getters: {
    userList: function(state) {
      let users = [];
      state.forEach(function(value) {
        users.push(value.trap_id);
      });
      return users;
    }
  }
};
