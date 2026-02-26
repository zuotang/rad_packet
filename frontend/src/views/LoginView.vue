<script setup>
import { computed, reactive, ref } from "vue";
import { useRouter } from "vue-router";
import { useAuthStore } from "../stores/auth";

const auth = useAuthStore();
const router = useRouter();
const loading = ref(false);
const errorMsg = ref("");

const form = reactive({
  phone: "",
  email: "",
  device_hash: "web-device-demo",
  country: "US",
  language: "zh-CN",
});

const inviteCode = computed(() => localStorage.getItem("invite_code") || "");

async function submit() {
  errorMsg.value = "";
  if (!form.phone && !form.email) {
    errorMsg.value = "手机号或邮箱至少填写一个";
    return;
  }
  loading.value = true;
  try {
    await auth.login(form);
    router.push("/");
  } catch (err) {
    errorMsg.value = err?.response?.data?.message || "登录失败，请稍后重试";
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <section class="card" style="max-width: 520px; margin: 40px auto">
    <h2>登录 / 注册</h2>
    <p class="muted">首次登录会自动创建账号与钱包。</p>
    <p v-if="inviteCode" class="badge">已捕获邀请码: {{ inviteCode }}</p>

    <div class="form" style="margin-top: 12px">
      <input v-model.trim="form.phone" placeholder="手机号（可选）" />
      <input v-model.trim="form.email" placeholder="邮箱（可选）" />
      <input v-model.trim="form.device_hash" placeholder="设备指纹（建议保留默认）" />
      <div class="row">
        <input v-model.trim="form.country" placeholder="国家，如 US" />
        <input v-model.trim="form.language" placeholder="语言，如 zh-CN" />
      </div>
      <p v-if="errorMsg" class="error">{{ errorMsg }}</p>
      <button :disabled="loading" @click="submit">{{ loading ? "登录中..." : "登录" }}</button>
    </div>
  </section>
</template>
