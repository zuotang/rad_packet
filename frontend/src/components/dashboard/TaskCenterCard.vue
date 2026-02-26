<script setup>
defineProps({
  tasks: { type: Array, default: () => [] },
  claimingId: { type: Number, default: 0 },
  loading: { type: Boolean, default: false },
});

const emit = defineEmits(["claim", "open-ledger"]);
</script>

<template>
  <section class="card">
    <div class="row space">
      <div>
        <h3>任务中心</h3>
        <p class="muted">完成任务可获得抽奖次数</p>
      </div>
      <button class="btn secondary" @click="emit('open-ledger')">流水</button>
    </div>

    <div v-if="loading" class="muted">任务加载中...</div>
    <div v-else-if="!tasks.length" class="muted">暂无可用任务</div>
    <div v-else class="list">
      <div class="list-item row" v-for="task in tasks" :key="task.id">
        <div>
          <strong>{{ task.name }}</strong>
          <p class="muted">抽奖次数 +{{ Math.round(Number(task.reward_amount || 0)) }}</p>
        </div>
        <button :disabled="claimingId === task.id" @click="emit('claim', task)">
          {{ claimingId === task.id ? "处理中..." : "领取" }}
        </button>
      </div>
    </div>
  </section>
</template>
