export const applicationDetail = {
  state: {
    application_id: "3fa85f64-5717-4562-b3fc-2c963f66afa6",
    created_at: "Thu, 1 May 2008 02:00:00 +0900",
    applicant: {
      trap_id: "nagatech",
      display_name: "ながてち",
      is_admin: true
    },
    current_state: "submitted",
    type: "club",
    title: "夏コミの交通費をお願いします。",
    remarks: "〇〇駅から〇〇駅への移動",
    paid_at: "Thu, 1 May 2008 02:00:00 +0900",
    ammount: 1200,
    repaid_to_id: ["nagatech", "series2"],
    images: ["3fa85f64-5717-4562-b3fc-2c963f66afa6"],
    comments: [
      {
        comment_id: 0,
        user: {
          trap_id: "nagatech",
          display_name: "ながてち",
          is_admin: true
        },
        comment: "コメント内容",
        created_at: "Thu, 1 May 2008 02:00:00 +0900",
        updated_at: "Thu, 1 May 2008 02:00:00 +0900"
      }
    ],
    state_logs: [
      {
        update_user: {
          trap_id: "nagatech",
          display_name: "ながてち",
          is_admin: true
        },
        to_state: "submitted",
        reason: "これは雑すぎますね。",
        created_at: "Thu, 1 May 2008 02:00:00 +0900"
      }
    ],
    application_detail_logs: [
      {
        update_user: {
          trap_id: "nagatech",
          display_name: "ながてち",
          is_admin: true
        },
        type: "club",
        title: "夏コミの交通費をお願いします。",
        remarks: "〇〇駅から〇〇駅への移動",
        ammount: 1200,
        paid_at: "2019-12-13",
        updated_at: "Thu, 1 May 2008 02:00:00 +0900"
      }
    ],
    repayment_logs: [
      {
        repaid_by_user: {
          trap_id: "nagatech",
          display_name: "ながてち",
          is_admin: true
        },
        repaid_to_user: {
          trap_id: "nagatech",
          display_name: "ながてち",
          is_admin: true
        },
        repaid_at: "Thu, 1 May 2008 02:00:00 +0900"
      }
    ]
  },
  getters: {
    logs: state => {
      let logs = [];
      state.comments.forEach(log => {
        logs.push({ log_type: "comment", content: log });
      });
      state.state_logs.forEach(log => {
        logs.push({ log_type: "state", content: log });
      });
      state.application_detail_logs.forEach(log => {
        logs.push({ log_type: "application", content: log });
      });
      state.repayment_logs.forEach(log => {
        logs.push({ log_type: "repayment", content: log });
      });
      //TODO:時系列順にlogを並び変える
      return logs;
    }
  }
};
