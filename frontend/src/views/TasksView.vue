<script setup>
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import api from "../api/client";
import { useRewardStore } from "../stores/reward";
import { useWalletStore } from "../stores/wallet";
import PageHeader from "../components/PageHeader.vue";
import TaskCenterCard from "../components/dashboard/TaskCenterCard.vue";
import BottomActionBar from "../components/dashboard/BottomActionBar.vue";

const router = useRouter();
const reward = useRewardStore();
const wallet = useWalletStore();

const loading = ref(true);
const errorMsg = ref("");
const hintMsg = ref("");
const claimingId = ref(0);
const tasks = ref([]);

async function loadTasks() {
  loading.value = true;
  errorMsg.value = "";
  try {
    const res = await api.get("/task/list");
    tasks.value = res.data.items || [];
  } catch (err) {
    errorMsg.value = err?.response?.data?.message || "任务加载失败";
  } finally {
    loading.value = false;
  }
}

onMounted(loadTasks);

async function claim(task) {
  claimingId.value = task.id;
  hintMsg.value = "";
  errorMsg.value = "";
  try {
    const res = await api.post("/task/claim", {
      task_id: task.id,
      event_key: `${task.id}-${Date.now()}`,
      meta_json: JSON.stringify({ from: "task_page" }),
    });
    await Promise.all([reward.fetchSummary(), wallet.fetchWallet(), loadTasks()]);
    hintMsg.value = `任务「${task.name}」领取成功，获得 ${res.data.spin_count || 0} 次抽奖`;
  } catch (err) {
    errorMsg.value = err?.response?.data?.message || "领取失败";
  } finally {
    claimingId.value = 0;
  }
}

function share() {
  window.open("https://wa.me/?text=" + encodeURIComponent("来一起做任务拿奖励"), "_blank", "noopener");
}
</script>

<template>
  <PageHeader title="任务中心" desc="完成任务，提升解锁进度并获得奖励。" />
  <TaskCenterCard :tasks="tasks" :loading="loading" :claiming-id="claimingId" @claim="claim" @open-ledger="router.push('/records')" />
  <section class="card" v-if="hintMsg"><p class="muted">{{ hintMsg }}</p></section>
  <section class="card" v-if="errorMsg"><p class="error">{{ errorMsg }}</p></section>
  <BottomActionBar @refresh="loadTasks" @share="share" />
</template>
