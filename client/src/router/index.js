import Vue from "vue";
import VueRouter from "vue-router";
import ApplicationListPage from "../views/ApplicationListPage.vue";
import ApplicationDetailPage from "../views/ApplicationDetailPage.vue";
import AdminPage from "../views/AdminPage";
import NewApplicationPage from "../views/NewApplicationPage.vue";
import store from "./../store/index";

Vue.use(VueRouter);

const routes = [
  {
    path: "/",
    name: "All Applications",
    component: ApplicationListPage
  },
  {
    path: "/admin",
    name: "Admin Page",
    component: AdminPage
  },
  {
    path: "/applications/new/:type",
    name: "New Application Page",
    component: NewApplicationPage
  },
  {
    path: "/applications/:id",
    name: "Application",
    component: ApplicationDetailPage
  },
  {
    path: "/callback",
    name: "callback",
    component: ApplicationListPage,
    beforeEnter: async (to, from, next) => {
      const code = to.query.code;
      const state = to.query.state;
      const sessionState = sessionStorage.getItem(`state`);
      if (state !== sessionState) {
        console.log("state error");
      } else {
        try {
          // Jomon server にapi
          console.log(code);
          const resp = await store.dispatch("getMe");
          await store.commit("setMe", resp.data);
          next("/");
        } catch (e) {
          console.error(e);
        }
      }
    }
  }
];

const router = new VueRouter({
  mode: "history",
  base: process.env.BASE_URL,
  routes
});

export default router;
