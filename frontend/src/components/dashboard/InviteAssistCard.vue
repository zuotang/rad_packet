<script setup>
import { computed } from "vue";

const props = defineProps({
  code: { type: String, default: "" },
  inviteCount: { type: Number, default: 0 },
  inviteNeed: { type: Number, default: 5 },
});

const emit = defineEmits(["copy"]);

const link = computed(() => `https://h5.xxx.com/r/${props.code || ""}`);

function shareWA() {
  const text = `我在 LuckyRain 做任务拿奖励，邀请码 ${props.code}：${link.value}`;
  window.open(`https://wa.me/?text=${encodeURIComponent(text)}`, "_blank", "noopener");
}

function shareFB() {
  window.open(
    `https://www.facebook.com/sharer/sharer.php?u=${encodeURIComponent(link.value)}`,
    "_blank",
    "noopener"
  );
}
</script>

<template>
  <section class="card">
    <div class="row space">
      <div>
        <h3>邀请助力</h3>
        <p class="muted">好友注册并完成首个有效任务后计入有效邀请</p>
      </div>
      <span class="pill">{{ inviteCount }}/{{ inviteNeed }}</span>
    </div>

    <div class="invite-box">
      <div class="row space">
        <span class="muted">邀请码</span>
        <strong class="badge">{{ code || "-" }}</strong>
      </div>
      <input :value="link" readonly />
      <div class="row">
        <button class="btn secondary" @click="emit('copy', link)">复制链接</button>
        <button class="btn" @click="shareWA">WhatsApp</button>
        <button class="btn" @click="shareFB">Facebook</button>
      </div>
    </div>
  </section>
</template>
