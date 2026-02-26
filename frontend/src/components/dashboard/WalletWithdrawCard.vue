<script setup>
import { ref, watch } from "vue";

const props = defineProps({
  balance: { type: Number, default: 0 },
  frozen: { type: Number, default: 0 },
  minWithdraw: { type: Number, default: 5 },
  submitting: { type: Boolean, default: false },
});

const emit = defineEmits(["withdraw", "rules"]);
const amount = ref(5);

watch(
  () => props.minWithdraw,
  (v) => {
    amount.value = v;
  },
  { immediate: true }
);

function submit() {
  emit("withdraw", Number(amount.value));
}
</script>

<template>
  <section class="card">
    <div class="row space">
      <div>
        <h3>钱包与提现</h3>
        <p class="muted">达到门槛后可提交提现申请</p>
      </div>
      <span class="pill">余额 {{ Number(balance).toFixed(2) }}</span>
    </div>

    <div class="row space">
      <span class="muted">可用</span>
      <strong>{{ Number(balance).toFixed(2) }}</strong>
    </div>
    <div class="row space">
      <span class="muted">冻结</span>
      <strong>{{ Number(frozen).toFixed(2) }}</strong>
    </div>
    <div class="row space">
      <span class="muted">最低提现</span>
      <strong>{{ Number(minWithdraw).toFixed(2) }}</strong>
    </div>

    <div class="row">
      <input v-model="amount" type="number" min="0.01" step="0.01" />
      <button :disabled="submitting" @click="submit">{{ submitting ? "提交中..." : "申请提现" }}</button>
      <button class="btn secondary" @click="emit('rules')">规则</button>
    </div>
  </section>
</template>
