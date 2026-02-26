<script setup>
import { onMounted, ref } from "vue";
import PageHeader from "../components/PageHeader.vue";
import InviteAssistCard from "../components/dashboard/InviteAssistCard.vue";
import BottomActionBar from "../components/dashboard/BottomActionBar.vue";
import { useReferralStore } from "../stores/referral";
import { useConfigStore } from "../stores/config";

const referral = useReferralStore();
const config = useConfigStore();
const hint = ref("");
const error = ref("");

onMounted(async () => {
  try {
    await Promise.all([referral.fetchStatus(), config.fetchBootstrap()]);
  } catch (err) {
    error.value = err?.response?.data?.message || "加载邀请数据失败";
  }
});

const inviteNeed = () => {
  if (config.rewardTiers?.length && config.rewardTiers[0]?.target) return Number(config.rewardTiers[0].target);
  return 5;
};

async function copyLink(link) {
  try {
    await navigator.clipboard.writeText(link);
    hint.value = "邀请链接已复制";
  } catch {
    error.value = "复制失败，请手动复制";
  }
}

function share() {
  const code = referral.myCode || "";
  const link = `https://h5.xxx.com/r/${code}`;
  const text = `我在 LuckyRain 做任务拿奖励，邀请码 ${code}：${link}`;
  window.open(`https://wa.me/?text=${encodeURIComponent(text)}`, "_blank", "noopener");
}
</script>

<template>
  <PageHeader title="邀请助力" desc="有效邀请会触发上级奖励，奖励先冻结再解锁。" />
  <InviteAssistCard :code="referral.myCode" :invite-count="referral.inviteCount" :invite-need="inviteNeed()" @copy="copyLink" />

  <section class="card">
    <h3>邀请规则</h3>
    <div class="list">
      <div class="list-item">同设备重复注册不计入有效邀请</div>
      <div class="list-item">异常 IP、风险账户不计入有效邀请</div>
      <div class="list-item">被邀请人完成首个有效任务后，才计入邀请</div>
      <div class="list-item">邀请奖励先进入 pending，风控通过后可解冻</div>
    </div>
  </section>

  <section class="card" v-if="hint"><p class="muted">{{ hint }}</p></section>
  <section class="card" v-if="error"><p class="error">{{ error }}</p></section>
  <BottomActionBar @refresh="referral.fetchStatus" @share="share" />
</template>
