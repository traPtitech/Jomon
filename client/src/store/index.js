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
    ]
  },
  mutations: {},
  actions: {},
  modules: {}
});
