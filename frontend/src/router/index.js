import { createRouter, createWebHistory } from "vue-router";
import LoginView from "../views/LoginView.vue";
import HomeView from "../views/HomeView.vue";
import TasksView from "../views/TasksView.vue";
import InviteView from "../views/InviteView.vue";
import WalletView from "../views/WalletView.vue";
import RecordsView from "../views/RecordsView.vue";
import { useAuthStore } from "../stores/auth";
import pinia from "../stores/pinia";

const routes = [
  { path: "/login", component: LoginView },
  { path: "/r/:code", component: LoginView, props: true },
  { path: "/", component: HomeView, meta: { auth: true } },
  { path: "/tasks", component: TasksView, meta: { auth: true } },
  { path: "/invite", component: InviteView, meta: { auth: true } },
  { path: "/wallet", component: WalletView, meta: { auth: true } },
  { path: "/records", component: RecordsView, meta: { auth: true } },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach((to) => {
  if (to.params.code) {
    localStorage.setItem("invite_code", String(to.params.code));
  }
  const auth = useAuthStore(pinia);
  if (to.meta.auth && !auth.token) {
    return "/login";
  }
  return true;
});

export default router;
