<script setup>
const emit = defineEmits(["refresh", "share"]);
</script>

<template>
  <div class="bottom-bar">
    <button class="btn secondary" @click="emit('refresh')">刷新数据</button>
    <button class="btn" @click="emit('share')">一键分享冲刺解锁</button>
  </div>
</template>
