<script setup>
import { computed } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useAuthStore } from "./stores/auth";

const route = useRoute();
const router = useRouter();
const auth = useAuthStore();
const showNav = computed(() => route.path !== "/login" && route.path !== "/" && !route.path.startsWith("/r/"));

function logout() {
  auth.logout();
  router.push("/login");
}
</script>

<template>
  <div class="layout">
    <header v-if="showNav" class="topbar card">
      <div class="brand">Red Packet</div>
      <nav class="nav">
        <router-link to="/" active-class="active">首页</router-link>
        <router-link to="/tasks" active-class="active">任务</router-link>
        <router-link to="/invite" active-class="active">邀请</router-link>
        <router-link to="/wallet" active-class="active">钱包</router-link>
        <router-link to="/records" active-class="active">记录</router-link>
      </nav>
      <button class="btn secondary" @click="logout">退出</button>
    </header>
    <main class="main">
      <router-view />
    </main>
  </div>
</template>
