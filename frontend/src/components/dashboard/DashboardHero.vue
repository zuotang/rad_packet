<script setup>
import { computed } from "vue";

const props = defineProps({
  pending: { type: Number, default: 0 },
  unlocked: { type: Number, default: 0 },
  cashable: { type: Number, default: 0 },
  inviteNeed: { type: Number, default: 5 },
  inviteCount: { type: Number, default: 0 },
  unlocking: { type: Boolean, default: false },
  unlockMessage: { type: String, default: "" },
});

const emit = defineEmits(["unlock"]);

const pct = computed(() => {
  if (props.inviteNeed <= 0) return 0;
  return Math.min(1, props.inviteCount / props.inviteNeed);
});
</script>

<template>
  <section class="hero-card card">
    <div class="row space">
      <div>
        <p class="muted">待解锁奖励</p>
        <div class="hero-big">{{ pending.toFixed(2) }}</div>
        <p class="muted">邀请或完成任务后可解锁到可提现余额</p>
      </div>
      <div class="pill">
        <span>可提现</span>
        <b>{{ cashable.toFixed(2) }}</b>
      </div>
    </div>

    <div class="progress-wrap">
      <div class="row space">
        <strong>解锁进度</strong>
        <span class="muted">{{ inviteCount }}/{{ inviteNeed }}</span>
      </div>
      <div class="progress-bar">
        <span class="progress-fill" :style="{ width: `${(pct * 100).toFixed(0)}%` }"></span>
      </div>
      <div class="row space">
        <span class="muted">已解锁 {{ unlocked.toFixed(2) }}</span>
        <span class="muted">还差 {{ Math.max(0, inviteNeed - inviteCount) }} 人</span>
      </div>
    </div>

    <div class="row" style="margin-top: 10px">
      <button class="btn" :disabled="unlocking" @click="emit('unlock')">
        {{ unlocking ? "解锁中..." : "解冻待解锁奖励" }}
      </button>
      <p class="muted" v-if="unlockMessage">{{ unlockMessage }}</p>
    </div>
  </section>
</template>
