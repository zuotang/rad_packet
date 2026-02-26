<script setup>
import { onMounted, ref } from "vue";
import api from "../api/client";
import { useWalletStore } from "../stores/wallet";
import PageHeader from "../components/PageHeader.vue";
import BottomActionBar from "../components/dashboard/BottomActionBar.vue";

const wallet = useWalletStore();
const tab = ref("wallet");
const loading = ref(true);
const error = ref("");
const rewardItems = ref([]);
const withdrawItems = ref([]);

async function loadAll() {
  loading.value = true;
  error.value = "";
  try {
    const [walletRes, rewardRes, withdrawRes] = await Promise.all([
      wallet.fetchWallet(),
      api.get("/reward/records?page=1&size=20"),
      api.get("/withdraw/records?page=1&size=20"),
    ]);
    rewardItems.value = rewardRes.data.items || [];
    withdrawItems.value = withdrawRes.data.items || [];
    return walletRes;
  } catch (err) {
    error.value = err?.response?.data?.message || "记录加载失败";
  } finally {
    loading.value = false;
  }
}

onMounted(loadAll);

function share() {
  window.open("https://wa.me/?text=" + encodeURIComponent("来一起做任务拿奖励"), "_blank", "noopener");
}

function fmtTime(v) {
  if (!v) return "-";
  return new Date(v).toLocaleString();
}
</script>

<template>
  <PageHeader title="流水与记录" desc="任务奖励、提现申请、账本流水统一查看。" />

  <section class="card">
    <div class="tabs">
      <button :class="['tab-btn', tab === 'wallet' ? 'active' : '']" @click="tab = 'wallet'">钱包流水</button>
      <button :class="['tab-btn', tab === 'reward' ? 'active' : '']" @click="tab = 'reward'">奖励记录</button>
      <button :class="['tab-btn', tab === 'withdraw' ? 'active' : '']" @click="tab = 'withdraw'">提现记录</button>
    </div>

    <div v-if="loading" class="muted">加载中...</div>

    <div class="list" v-else-if="tab === 'wallet'">
      <div class="list-item" v-for="item in wallet.ledgers" :key="item.id">
        <div class="row space">
          <strong>{{ item.type }}</strong>
          <span :style="{ color: item.amount >= 0 ? '#0f766e' : '#b91c1c' }">{{ item.amount >= 0 ? '+' : '' }}{{ Number(item.amount).toFixed(2) }}</span>
        </div>
        <p class="muted">{{ item.ref_type }} / {{ item.ref_id }}</p>
        <p class="muted">{{ fmtTime(item.created_at) }}</p>
      </div>
      <p class="muted" v-if="!wallet.ledgers.length">暂无钱包流水</p>
    </div>

    <div class="list" v-else-if="tab === 'reward'">
      <div class="list-item" v-for="item in rewardItems" :key="item.id || item.ID">
        <div class="row space">
          <strong>{{ item.reward_type || item.RewardType || 'reward' }}</strong>
          <span class="badge">{{ item.status || item.Status }}</span>
        </div>
        <p class="muted">金额：{{ Number(item.amount || item.Amount || 0).toFixed(2) }}</p>
        <p class="muted">来源：{{ item.source_type || item.SourceType }} / {{ item.source_id || item.SourceID }}</p>
      </div>
      <p class="muted" v-if="!rewardItems.length">暂无奖励记录</p>
    </div>

    <div class="list" v-else>
      <div class="list-item" v-for="item in withdrawItems" :key="item.id || item.ID">
        <div class="row space">
          <strong>#{{ item.id || item.ID }}</strong>
          <span class="badge">{{ item.status || item.Status }}</span>
        </div>
        <p class="muted">金额：{{ Number(item.amount || item.Amount || 0).toFixed(2) }}</p>
        <p class="muted">备注：{{ item.note || item.Note || '-' }}</p>
      </div>
      <p class="muted" v-if="!withdrawItems.length">暂无提现记录</p>
    </div>
  </section>

  <section class="card" v-if="error"><p class="error">{{ error }}</p></section>
  <BottomActionBar @refresh="loadAll" @share="share" />
</template>
