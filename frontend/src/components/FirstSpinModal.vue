<script setup>
const props = defineProps({
  show: { type: Boolean, default: false },
  title: { type: String, default: "" },
  subtitle: { type: String, default: "" },
  tip: { type: String, default: "点击任意位置继续" },
});

const emit = defineEmits(["close"]);

function close() {
  emit("close");
}
</script>

<template>
  <div v-if="show" class="fs-mask show" aria-hidden="false" @click="close">
    <div class="fs-banner" role="dialog" aria-modal="true" aria-label="首次抽奖提示">
      <div class="fs-bg"></div>
      <div class="fs-glow"></div>
      <div class="fs-edge top"></div>
      <div class="fs-edge bottom"></div>
      <div class="fs-sweep"></div>
      <div class="fs-beam top"></div>
      <div class="fs-beam bottom"></div>

      <div class="fs-text">
        <div>
          <span class="t1">{{ title }}</span>
          <span class="t2">{{ subtitle }}</span>
        </div>
      </div>
    </div>
    <div class="fs-tip">{{ tip }}</div>
  </div>
</template>

<style scoped>
.fs-mask {
  position: fixed;
  inset: 0;
  z-index: 99999;
  background: radial-gradient(circle at 50% 35%, rgba(0, 0, 0, 0.45), rgba(0, 0, 0, 0.82) 70%, rgba(0, 0, 0, 0.92) 100%);
  display: none;
  align-items: center;
  justify-content: center;
  padding: 18px;
}

.fs-mask.show {
  display: flex;
}

.fs-banner {
  width: min(92vw, 520px);
  height: 132px;
  border-radius: 14px;
  position: relative;
  overflow: hidden;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.45);
  transform: translateZ(0);
  animation: fs-pop 0.22s ease-out both;
}

.fs-bg {
  position: absolute;
  inset: 0;
  background: linear-gradient(180deg, #ff6a47 0%, #ff3b2f 55%, #ff2e2e 100%);
}

.fs-bg::after {
  content: "";
  position: absolute;
  inset: -2px;
  opacity: 0.18;
  background: linear-gradient(transparent 0 94%, rgba(255, 255, 255, 0.55) 94% 100%),
    linear-gradient(90deg, transparent 0 94%, rgba(255, 255, 255, 0.55) 94% 100%);
  background-size: 18px 18px;
  mix-blend-mode: overlay;
}

.fs-glow {
  position: absolute;
  inset: -50px;
  background: radial-gradient(circle at 50% 50%, rgba(255, 210, 74, 0.35), rgba(255, 210, 74, 0) 55%),
    radial-gradient(circle at 20% 50%, rgba(255, 255, 255, 0.18), rgba(255, 255, 255, 0) 55%),
    radial-gradient(circle at 80% 50%, rgba(255, 255, 255, 0.18), rgba(255, 255, 255, 0) 55%);
  filter: blur(6px);
  pointer-events: none;
}

.fs-edge {
  position: absolute;
  left: 0;
  right: 0;
  height: 10px;
  background: linear-gradient(90deg, transparent, rgba(255, 245, 210, 0.95), transparent);
  filter: blur(0.2px);
  opacity: 0.9;
}

.fs-edge.top {
  top: -2px;
}

.fs-edge.bottom {
  bottom: -2px;
}

.fs-sweep {
  position: absolute;
  top: -40%;
  left: -40%;
  width: 55%;
  height: 180%;
  background: linear-gradient(
    90deg,
    transparent,
    rgba(255, 255, 255, 0.2),
    rgba(255, 255, 255, 0.65),
    rgba(255, 255, 255, 0.18),
    transparent
  );
  transform: skewX(-18deg);
  animation: fs-sweep 1.55s ease-in-out infinite;
  mix-blend-mode: screen;
  filter: blur(0.2px);
  pointer-events: none;
}

.fs-text {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 14px 16px;
  font-weight: 1000;
  line-height: 1.25;
  letter-spacing: 0.5px;
}

.fs-text .t1 {
  font-size: 34px;
  color: #fff;
  text-shadow: 0 4px 10px rgba(0, 0, 0, 0.25);
}

.fs-text .t2 {
  display: block;
  margin-top: 10px;
  font-size: 34px;
  color: #ffd24a;
  text-shadow: 0 4px 10px rgba(0, 0, 0, 0.2);
}

.fs-tip {
  position: absolute;
  bottom: -40px;
  left: 50%;
  transform: translateX(-50%);
  color: rgba(255, 255, 255, 0.82);
  font-weight: 900;
  font-size: 13px;
  letter-spacing: 0.2px;
  text-shadow: 0 2px 8px rgba(0, 0, 0, 0.25);
  user-select: none;
  white-space: nowrap;
}

.fs-beam {
  position: absolute;
  left: -20%;
  width: 140%;
  height: 10px;
  opacity: 0.8;
  filter: blur(0.4px);
  mix-blend-mode: screen;
  pointer-events: none;
  background: linear-gradient(
    90deg,
    transparent,
    rgba(255, 255, 255, 0.15),
    rgba(255, 245, 210, 0.9),
    rgba(255, 255, 255, 0.2),
    transparent
  );
}

.fs-beam.top {
  top: -12px;
  animation: fs-beam-left 1.6s ease-in-out infinite;
}

.fs-beam.bottom {
  bottom: -12px;
  animation: fs-beam-right 1.6s ease-in-out infinite;
}

@keyframes fs-beam-left {
  0% {
    transform: translateX(-20%);
    opacity: 0;
  }
  20% {
    opacity: 0.9;
  }
  100% {
    transform: translateX(20%);
    opacity: 0;
  }
}

@keyframes fs-beam-right {
  0% {
    transform: translateX(20%);
    opacity: 0;
  }
  20% {
    opacity: 0.9;
  }
  100% {
    transform: translateX(-20%);
    opacity: 0;
  }
}

@keyframes fs-sweep {
  0% {
    left: -60%;
    opacity: 0;
  }
  10% {
    opacity: 0.95;
  }
  55% {
    opacity: 0.95;
  }
  100% {
    left: 120%;
    opacity: 0;
  }
}

@keyframes fs-pop {
  from {
    transform: scale(0.96);
    opacity: 0;
  }
  to {
    transform: scale(1);
    opacity: 1;
  }
}
</style>
