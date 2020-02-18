//@ts-check
import Vue from "vue";
import Vuex from "vuex";

Vue.use(Vuex);

/**
 * 申請書の状態
 * @enum {string}
 */
const State = {
  /**申請中 */
  submitted: "submitted",
  /**申請却下 */
  rejected: "rejected",
  /**要修正 */
  fix_required: "fix_required",
  /**申請を受諾 */
  accepted: "accepted",
  /**返金済み */
  fully_repaid: "fully_repaid"
};

/**
 * @typedef {object} application 申請書の概要
 * @property {number} id 申請書のuuid
 * @property {string} title 申請書のタイトル
 * @property {string} name 申請者の名前
 * @property {number} money 申請金額
 * @property {string} state 申請書の状態
 */

export default new Vuex.Store({
  state: {
    /**
     * @type application[]
     */
    applicationList: [
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
    ],
    me: { trap_id: "nagatech", display_name: "ながてち", is_admin: true },

    application_detail_paper: {
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
      images: ["3fa85f64-5717-4562-b3fc-2c963f66afa6"]
    },
    logs: [
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
  },
  mutations: {},
  actions: {},
  modules: {}
});
