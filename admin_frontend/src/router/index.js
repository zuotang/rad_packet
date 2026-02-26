import { createRouter, createWebHistory } from "vue-router";
import { useAdminStore } from "../stores/admin";
import LoginView from "../views/LoginView.vue";
import DashboardView from "../views/DashboardView.vue";

const routes = [
  { path: "/login", component: LoginView },
  { path: "/", component: DashboardView, meta: { auth: true } },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach((to) => {
  const admin = useAdminStore();
  if (to.meta.auth && !admin.adminKey) {
    return "/login";
  }
  return true;
});

export default router;
