import { useMeStore } from "@/stores/me";
import { redirectAuthEndpoint } from "@/use/api";
import AdminPage from "@/views/AdminPage.vue";
import ApplicationDetailPage from "@/views/ApplicationDetailPage.vue";
import ApplicationListPage from "@/views/ApplicationListPage.vue";
import NewApplicationPage from "@/views/NewApplicationPage.vue";
import { createRouter, createWebHistory } from "vue-router";

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
