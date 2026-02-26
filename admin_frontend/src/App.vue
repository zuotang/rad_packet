<script setup>
import { useRouter } from "vue-router";
import { useAdminStore } from "./stores/admin";

const router = useRouter();
const admin = useAdminStore();

function logout() {
  admin.logout();
  router.push("/login");
}
</script>

<template>
  <div class="layout">
    <header class="top card">
      <strong>Red Packet Admin</strong>
      <div class="row">
        <span class="muted" v-if="admin.adminKey">已登录后台</span>
        <button class="btn secondary" v-if="admin.adminKey" @click="logout">退出</button>
      </div>
    </header>
    <main class="main">
      <router-view />
    </main>
  </div>
</template>
