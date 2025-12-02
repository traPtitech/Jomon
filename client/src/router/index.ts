import { useMeStore } from "@/stores/me";
import { redirectAuthEndpoint } from "@/use/api";
import { createRouter, createWebHistory, RouteRecordRaw } from "vue-router";

const routes: Array<RouteRecordRaw> = [
  {
    path: "/",
    name: "All Applications",
    component: () => import("@/views/ApplicationListPage.vue")
  },
  {
    path: "/admin",
    name: "Admin Page",
    component: () => import("@/views/AdminPage.vue")
  },
  {
    path: "/applications/new/:type",
    name: "New Application Page",
    component: () => import("@/views/NewApplicationPage.vue")
  },
  {
    path: "/applications/:id",
    name: "Application",
    component: () => import("@/views/ApplicationDetailPage.vue")
  }
];

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
});

router.beforeEach(async (to, _from, next) => {
  const meStore = useMeStore();
  if (meStore.trapId === "") {
    try {
      await meStore.fetchMe();
    } catch {
      sessionStorage.setItem(`destination`, to.fullPath);
      await redirectAuthEndpoint();
    }
  }
  next();
});

export default router;
