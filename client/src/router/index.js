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
    path: "/applications/new",
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

export default router;
