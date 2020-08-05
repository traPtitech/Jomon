import Vue from "vue";
import VueRouter from "vue-router";
import ApplicationListPage from "../views/ApplicationListPage.vue";
import ApplicationDetailPage from "../views/ApplicationDetailPage.vue";
import AdminPage from "../views/AdminPage";
import NewApplicationPage from "../views/NewApplicationPage.vue";

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
  }
];

const router = new VueRouter({
  mode: "history",
  base: process.env.BASE_URL,
  routes
});

router.beforeEach(async (_to, _from, next) => {
  if (!store.state.me) {
    await store.dispatch.getMe();
  }
  if (!store.state.me) {
    sessionStorage.setItem(`destination`, to.fullPath);
    redirect2AuthEndpoint();
  }
  next();
});

export default router;
