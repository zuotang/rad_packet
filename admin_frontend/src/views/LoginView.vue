<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";
import { useAdminStore } from "../stores/admin";

const router = useRouter();
const admin = useAdminStore();
const key = ref(admin.adminKey);
const error = ref("");

function submit() {
  if (!key.value.trim()) {
    error.value = "请输入管理员密钥";
    return;
  }
  admin.setKey(key.value.trim());
  router.push("/");
}
</script>

<template>
  <section class="card login">
    <h2>后台登录</h2>
    <p class="muted">请输入后端配置的 `APP_ADMIN_KEY`</p>
    <div class="form">
      <input v-model="key" placeholder="X-Admin-Key" />
      <p class="error" v-if="error">{{ error }}</p>
      <button class="btn" @click="submit">进入后台</button>
    </div>
  </section>
</template>
