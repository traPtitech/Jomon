//@ts-check
export const logs = {
  state: [
    {
      log_type: "comment",
      content: {
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
    },
    {
      log_type: "state",
      content: {
        update_user: {
          trap_id: "nagatech",
          display_name: "ながてち",
          is_admin: true
        },
        to_state: "submitted",
        reason: "これは雑すぎますね。",
        created_at: "Thu, 1 May 2008 02:00:00 +0900"
      }
    },
    {
      log_type: "application",
      content: {
        update_user: {
          trap_id: "nagatech",
          display_name: "ながてち",
          is_admin: true
        },
        type: "club",
        title: "夏コミの交通費をお願いします。",
        remarks: "〇〇駅から〇〇駅への移動",
        ammount: 1200,
        paid_at: "Thu, 1 May 2008 02:00:00 +0900",
        updated_at: "Thu, 1 May 2008 02:00:00 +0900"
      }
    },
    {
      log_type: "repayment",
      content: {
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
    }
  ]
};
