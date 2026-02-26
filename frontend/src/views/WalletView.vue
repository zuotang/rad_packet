<script setup>
import { onMounted, ref } from "vue";
import api from "../api/client";
import { useWalletStore } from "../stores/wallet";
import { useConfigStore } from "../stores/config";
import PageHeader from "../components/PageHeader.vue";
import WalletWithdrawCard from "../components/dashboard/WalletWithdrawCard.vue";
import BottomActionBar from "../components/dashboard/BottomActionBar.vue";

const wallet = useWalletStore();
const config = useConfigStore();
const withdrawing = ref(false);
const hint = ref("");
const error = ref("");
const withdrawItems = ref([]);

async function loadPage() {
  try {
    await config.fetchBootstrap();
    await wallet.fetchWallet();
    const wr = await api.get("/withdraw/records?page=1&size=10");
    withdrawItems.value = wr.data.items || [];
  } catch (err) {
    error.value = err?.response?.data?.message || "钱包数据加载失败";
  }
}

onMounted(loadPage);

async function applyWithdraw(amount) {
  hint.value = "";
  error.value = "";
  withdrawing.value = true;
  try {
    await wallet.applyWithdraw(amount);
    await loadPage();
    hint.value = "提现申请已提交";
  } catch (err) {
    error.value = err?.response?.data?.message || "提现申请失败";
  } finally {
    withdrawing.value = false;
  }
}

function share() {
  window.open("https://wa.me/?text=" + encodeURIComponent("来一起做任务拿奖励"), "_blank", "noopener");
}
</script>

<template>
  <PageHeader title="钱包与提现" desc="可用余额达到门槛后可提交提现申请。" />
  <WalletWithdrawCard
    :balance="wallet.balance"
    :frozen="wallet.frozen"
    :min-withdraw="config.withdrawMin"
    :submitting="withdrawing"
    @withdraw="applyWithdraw"
    @rules="$router.push('/records')"
  />

  <section class="card">
    <h3>最近提现申请</h3>
    <div class="list" v-if="withdrawItems.length">
      <div class="list-item" v-for="item in withdrawItems" :key="item.id || item.ID">
        <div class="row space">
          <strong>#{{ item.id || item.ID }}</strong>
          <span class="badge">{{ item.status || item.Status }}</span>
        </div>
        <p class="muted">金额：{{ Number(item.amount || item.Amount || 0).toFixed(2) }}</p>
        <p class="muted">{{ (item.note || item.Note || "-") }}</p>
      </div>
    </div>
    <p class="muted" v-else>暂无提现记录</p>
  </section>

  <section class="card" v-if="hint"><p class="muted">{{ hint }}</p></section>
  <section class="card" v-if="error"><p class="error">{{ error }}</p></section>
  <BottomActionBar @refresh="loadPage" @share="share" />
</template>
