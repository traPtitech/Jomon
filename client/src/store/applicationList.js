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
      id: 10,
      title: "タイトル1",
      name: "user1",
      money: 100,
      state: State.accepted
    },
    {
      id: 13,
      title: "タイトル2",
      name: "user2",
      money: 500,
      state: State.accepted
    },
    {
      id: 16,
      title: "タイトル3",
      name: "user3",
      money: 2500,
      state: State.accepted
    }
  ]
};
