<script setup>
defineProps({
  visible: { type: Boolean, default: false },
  ledgers: { type: Array, default: () => [] },
});

const emit = defineEmits(["close"]);

function toNum(v) {
  return Number(v || 0);
}
</script>

<template>
  <div v-if="visible" class="modal-mask" @click.self="emit('close')">
    <div class="modal-panel card">
      <div class="row space">
        <h3>奖励流水</h3>
        <button class="btn secondary" @click="emit('close')">关闭</button>
      </div>
      <div class="list" style="max-height: 50vh; overflow: auto">
        <div v-if="!ledgers.length" class="muted">暂无流水</div>
        <div class="list-item" v-for="item in ledgers" :key="item.id">
          <div class="row space">
            <strong>{{ item.type }}</strong>
            <span :style="{ color: toNum(item.amount) >= 0 ? '#0f766e' : '#b91c1c' }">
              {{ toNum(item.amount) >= 0 ? "+" : "" }}{{ Math.abs(toNum(item.amount)).toFixed(2) }}
            </span>
          </div>
          <p class="muted">{{ item.ref_type }} / {{ item.ref_id }}</p>
        </div>
      </div>
    </div>
  </div>
</template>
