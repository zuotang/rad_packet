
<script setup>
import { computed, onMounted, ref } from "vue";
import api from "../api/client";
import { useAuthStore } from "../stores/auth";
import { useConfigStore } from "../stores/config";
import { useReferralStore } from "../stores/referral";
import { useRewardStore } from "../stores/reward";
import { useWalletStore } from "../stores/wallet";
import FirstSpinModal from "../components/FirstSpinModal.vue";

const auth = useAuthStore();
const config = useConfigStore();
const referral = useReferralStore();
const reward = useRewardStore();
const wallet = useWalletStore();

const loading = ref(true);
const spinning = ref(false);
const currentDeg = ref(0);
const toastVisible = ref(true);
const toastText = ref("");
const activeTab = ref("show");
const errorMessage = ref("");
const hintMessage = ref("");
const showInviteModal = ref(false);
const showExitModal = ref(false);
const showFirstSpin = ref(false);

const lotteryStatus = ref({
  spin_count: 0,
  target: 60,
  balance: 0,
  pending: 0,
  needed: 0,
  unlockable: 0,
});

const withdrawRecords = ref([]);
const mockWithdrawRecords = [
  { id: "mock-1", amount: 60, status: "已提现", created_at: "2026-02-25" },
  { id: "mock-2", amount: 60, status: "已提现", created_at: "2026-02-24" },
  { id: "mock-3", amount: 60, status: "已提现", created_at: "2026-02-23" },
];
const displayWithdrawRecords = computed(() =>
  withdrawRecords.value.length ? withdrawRecords.value : mockWithdrawRecords
);

const prizeSegments = [
  { label: "提现红包", type: "mid" },
  { label: "0.01元~3元", type: "small" },
  { label: "高级红包", type: "mid" },
  { label: "谢谢参与", type: "thanks" },
  { label: "提现红包", type: "mid" },
  { label: "0.01元~3元", type: "small" },
  { label: "高级红包", type: "mid" },
  { label: "谢谢参与", type: "thanks" },
  { label: "提现红包", type: "mid" },
  { label: "0.01元~3元", type: "small" },
  { label: "高级红包", type: "mid" },
  { label: "谢谢参与", type: "thanks" },
];
const segmentAngle = 360 / prizeSegments.length;

const guideSteps = [
  "先把“待获得金额 3 元”补齐（你这里显示还差 3 元）。",
  "优先完成高权重任务：邀请助力、下单、签到等（以活动规则为准）。",
  "剩余次数不足时，先去“扫码助力”拉满，再回来抽奖更快。",
  "注意倒计时：到期后现金失效，建议尽快完成兑换并提现。",
];

const targetAmount = computed(() => Number(lotteryStatus.value.target || config.withdrawMin || 60));
const balance = computed(() => Number(lotteryStatus.value.balance || 0));
const pending = computed(() => Number(lotteryStatus.value.pending || 0));
const displayAmount = computed(() => balance.value + pending.value);
const neededAmount = computed(() => Math.max(0, Number(lotteryStatus.value.needed || targetAmount.value - displayAmount.value)));
const leftTimes = computed(() => Number(lotteryStatus.value.spin_count || 0));
const displayName = computed(() => {
  if (auth.user?.phone) {
    return maskPhone(auth.user.phone);
  }
  if (auth.user?.email) {
    return auth.user.email;
  }
  return "用户";
});

const exitProgress = computed(() => {
  if (targetAmount.value <= 0) return 0;
  return Math.max(0, Math.min(100, (displayAmount.value / targetAmount.value) * 100));
});

function formatMoney(value) {
  const amount = Number(value || 0);
  if (!Number.isFinite(amount)) return "0";
  const fixed = amount.toFixed(2);
  return fixed.replace(/\.00$/, "");
}

function maskPhone(phone) {
  const raw = String(phone || "");
  if (raw.length < 7) return raw;
  return `${raw.slice(0, 3)}****${raw.slice(-4)}`;
}

function formatDate(value) {
  if (!value) return "";
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return "";
  return new Intl.DateTimeFormat("zh-CN", { year: "numeric", month: "2-digit", day: "2-digit" }).format(date);
}

function showFirstSpinModal() {
  showFirstSpin.value = true;
  localStorage.setItem("first_spin_shown", "1");
}

function hideFirstSpinModal() {
  showFirstSpin.value = false;
}

async function loadWithdrawRecords() {
  const res = await api.get("/withdraw/records?page=1&size=5");
  withdrawRecords.value = (res.data.items || []).map((item) => ({
    id: item.id ?? item.ID,
    amount: Number(item.amount ?? item.Amount ?? 0),
    status: item.status ?? item.Status,
    created_at: item.created_at ?? item.CreatedAt,
  }));
}

async function loadStatus() {
  const res = await api.get("/lottery/status");
  lotteryStatus.value = res.data || lotteryStatus.value;
}

async function loadAll() {
  loading.value = true;
  errorMessage.value = "";
  try {
    await Promise.all([
      config.fetchBootstrap(),
      referral.fetchStatus(),
      reward.fetchSummary(),
      wallet.fetchWallet(),
      loadStatus(),
      loadWithdrawRecords(),
    ]);
  } catch (err) {
    errorMessage.value = err?.response?.data?.message || "首页数据加载失败";
  } finally {
    loading.value = false;
  }
}

async function spin() {
  if (spinning.value) return;
  if (!localStorage.getItem("first_spin_shown")) {
    showFirstSpinModal();
    return;
  }
  if (leftTimes.value <= 0) {
    toastText.value = "暂无抽奖次数，完成任务或邀请好友获取次数。";
    toastVisible.value = true;
    return;
  }

  spinning.value = true;
  toastVisible.value = false;

  let segmentIndex = 0;
  let amount = 0;
  try {
    const res = await api.post("/lottery/spin");
    segmentIndex = res.data.segment_index ?? 0;
    amount = Number(res.data.amount || 0);
    lotteryStatus.value.spin_count = res.data.spin_count ?? lotteryStatus.value.spin_count;
  } catch (err) {
    errorMessage.value = err?.response?.data?.message || "抽奖失败";
    toastText.value = "抽奖未成功，请稍后再试";
    toastVisible.value = true;
    spinning.value = false;
    return;
  }

  const extraTurns = randInt(4, 6);
  const targetDeg = 360 * extraTurns + (360 - (segmentIndex * segmentAngle + segmentAngle / 2));
  currentDeg.value += targetDeg;

  setTimeout(async () => {
    try {
      await Promise.all([reward.fetchSummary(), wallet.fetchWallet(), loadStatus(), loadWithdrawRecords()]);
      if (amount > 0) {
        hintMessage.value = `抽奖成功，获得 ${formatMoney(amount)} 元`;
        toastText.value = `本次抽中：${formatMoney(amount)} 元`;
      } else {
        hintMessage.value = "本次未中奖，继续努力";
        toastText.value = "谢谢参与，再试一次";
      }
    } catch (err) {
      errorMessage.value = err?.response?.data?.message || "刷新抽奖数据失败";
    } finally {
      spinning.value = false;
      toastVisible.value = true;
    }
  }, 3300);
}

function randInt(min, max) {
  return Math.floor(Math.random() * (max - min + 1)) + min;
}

async function applyWithdraw() {
  errorMessage.value = "";
  hintMessage.value = "";
  if (balance.value < targetAmount.value) {
    toastText.value = `还差 ${formatMoney(neededAmount.value)} 元才能提现`;
    toastVisible.value = true;
    return;
  }
  try {
    await wallet.applyWithdraw(targetAmount.value);
    await Promise.all([wallet.fetchWallet(), loadWithdrawRecords(), loadStatus()]);
    hintMessage.value = "提现申请已提交";
  } catch (err) {
    errorMessage.value = err?.response?.data?.message || "提现申请失败";
  }
}

function shareInvite(channel) {
  const code = referral.myCode || "";
  const link = `${window.location.origin}/r/${code}`;
  const text = `我在 LuckyRain 做任务拿红包，邀请码 ${code}，点这里一起参与：${link}`;
  if (navigator.share) {
    navigator.share({ title: "邀请", text, url: link }).catch(() => {});
    return;
  }
  if (channel === "copy") {
    navigator.clipboard.writeText(`${text}\n${link}`).catch(() => {});
    return;
  }
  window.open(`https://wa.me/?text=${encodeURIComponent(text)}`, "_blank", "noopener");
}

function openInviteModal() {
  showInviteModal.value = true;
}

function closeInviteModal() {
  showInviteModal.value = false;
}

function openExitModal() {
  showExitModal.value = true;
}

function closeExitModal() {
  showExitModal.value = false;
}

onMounted(loadAll);
</script>

<template>
  <div class="app">
    <div class="topbar">
      <div class="toprow">
        <div class="iconbtn" aria-label="back" @click="openExitModal">‹</div>
        <div class="title">现金大转盘</div>
        <div class="right">
          <div class="iconbtn" title="sound">🔊</div>
          <div class="iconbtn" title="rules" @click="activeTab = 'guide'">规则</div>
          <div class="iconbtn" title="detail" @click="activeTab = 'show'">明细</div>
        </div>
      </div>
    </div>

    <div class="content">
      <div class="banner">
        <div class="avatar"></div>
        <div class="txt">用户 {{ displayName }} 已累计获得金额：</div>
        <div class="amount">{{ formatMoney(displayAmount) }}元</div>
      </div>

      <div class="card">
        <div class="cardhead">
          <div class="name">
            <div class="avatar"></div>
            <div>我的账户</div>
          </div>
          <div class="pill">
            现金账户
            <span class="check">✓</span>
            <b>{{ formatMoney(displayAmount) }}</b>
          </div>
        </div>

        <div class="headline">离可提现仅差<span class="hot">{{ formatMoney(neededAmount) }}</span>元</div>

        <div class="bigmoney">{{ formatMoney(displayAmount) }}<span>元</span></div>

        <div class="stats">
          <div class="stat">
            <div class="k">已获得金额</div>
            <div class="v">{{ formatMoney(displayAmount) }}<small>元</small></div>
          </div>
          <div class="stat">
            <div class="k">待获得金额</div>
            <div class="v red">{{ formatMoney(neededAmount) }}<small>元</small></div>
          </div>
        </div>
      </div>

      <div class="wish">祝你今天顺利提现 <b>{{ formatMoney(targetAmount) }}元</b>!</div>

      <div class="wheelWrap">
        <div class="wheelFrame">
          <div class="mini">
            <div class="chip">再抽 {{ leftTimes }} 次</div>
          </div>

          <div class="pointer"></div>

          <div id="wheel" class="wheel" :style="{ '--rot': currentDeg + 'deg' }">
            <div v-for="(seg, idx) in prizeSegments" :key="idx" class="prize" :style="{ '--a': idx * segmentAngle + 'deg' }">
              <span :class="{ gray: seg.type === 'thanks' }">{{ seg.label }}</span>
            </div>
          </div>

          <div id="spinBtn" class="centerBtn" role="button" aria-label="抽奖" @click="spin">
            <div class="tag">提现冲刺阶段</div>
            <div class="main">抽奖</div>
            <div class="sub">还剩<i id="leftTimes">{{ leftTimes }}</i>次</div>
          </div>

          <div v-if="toastVisible" class="toast" id="toast">
            {{ toastText || "你是活动优质用户，超容易提现" }}
            <small>活动优质用户指 30 日内未参与现金大转盘活动的用户</small>
          </div>
        </div>
      </div>

      <div class="bottomArea">
        <div class="ctaWrap">
          <button class="cta ctaPulse" id="easyWithdrawBtn" @click="applyWithdraw">超容易提现{{ formatMoney(targetAmount) }}元</button>
          <div class="handHint" aria-hidden="true">👆</div>
        </div>
        <div class="countdown">23:29:34:8 后现金失效</div>

        <div class="scan">
          <div class="btn" title="扫码助力" @click="openInviteModal">
            <div class="qr"></div>
          </div>
          <div class="label" @click="openInviteModal">扫码助力</div>
        </div>
      </div>

      <section v-if="loading" class="status-card">数据加载中…</section>
      <section v-if="errorMessage" class="status-card error">{{ errorMessage }}</section>

      <div class="panel">
        <div class="tabs">
          <button class="tab" :class="{ active: activeTab === 'help' }" data-tab="help" @click="activeTab = 'help'">助力记录</button>
          <button class="tab" :class="{ active: activeTab === 'show' }" data-tab="show" @click="activeTab = 'show'">提现晒单</button>
          <button class="tab" :class="{ active: activeTab === 'guide' }" data-tab="guide" @click="activeTab = 'guide'">提现攻略</button>
        </div>

        <div class="tabBody">
          <div class="tabPage" :class="{ active: activeTab === 'help' }" id="tab-help">
            <div v-if="!referral.directInvites?.length" class="empty">暂无助力记录，邀请好友一起赚红包。</div>
            <div v-for="(item, idx) in referral.directInvites" :key="item.id || idx" class="receiptItem">
              <div class="uavatar" :class="idx % 2 === 0 ? 'a2' : 'a3'"></div>
              <div class="info">
                <div class="row1">
                  <div class="uname">好友 {{ maskPhone(item.phone || '******') }}</div>
                  <div class="date">{{ formatDate(item.created_at || item.CreatedAt) }}</div>
                </div>
                <div class="desc">扫码助力成功</div>
              </div>
              <div class="badge">
                <div class="money">+1</div>
                <div class="status">助力</div>
              </div>
            </div>
          </div>

          <div class="tabPage" :class="{ active: activeTab === 'show' }" id="tab-show">
            <div v-if="!displayWithdrawRecords.length" class="empty">暂无提现记录，继续抽奖累积金额。</div>
            <div v-for="item in displayWithdrawRecords" :key="item.id" class="receiptItem">
              <div class="uavatar a1"></div>
              <div class="info">
                <div class="row1">
                  <div class="uname">用户 {{ displayName }}</div>
                  <div class="date">{{ formatDate(item.created_at) }}</div>
                </div>
                <div class="desc">获得 {{ formatMoney(item.amount) }} 元现金，现金打款已到账！</div>
              </div>
              <div class="badge">
                <div class="money">{{ formatMoney(item.amount) }}<span>元</span></div>
                <div class="status">{{ item.status || "已提现" }}</div>
              </div>
            </div>
          </div>

          <div class="tabPage" :class="{ active: activeTab === 'guide' }" id="tab-guide">
            <div class="guideCard">
              <div class="gtitle">提现攻略</div>
              <ol class="glist">
                <li v-for="(line, idx) in guideSteps" :key="idx">{{ line }}</li>
              </ol>
              <div class="gtip">小提示：如你属于“活动优质用户”，通常更容易达成提现门槛。</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <FirstSpinModal
    :show="showFirstSpin"
    title="太棒啦！你是首次参与提现"
    subtitle="送你首次提现福利"
    tip="点击任意位置继续"
    @close="hideFirstSpinModal"
  />

  <div v-if="showInviteModal" id="inviteModal" class="modalMask show" aria-hidden="false">
    <div class="inviteCard" role="dialog" aria-modal="true" aria-label="邀请好友">
      <button class="modalClose" id="inviteClose" aria-label="关闭" @click="closeInviteModal">×</button>

      <div class="inviteMain">
        <div class="inviteAvatar"></div>
        <div class="inviteName" id="inviteName">{{ displayName }}</div>

        <div class="inviteText">
          You are <b>very close</b> to withdrawing <span class="green">¥{{ formatMoney(targetAmount) }}</span>, only missing
        </div>

        <div class="inviteBig">
          <span id="inviteDiff">{{ formatMoney(neededAmount) }}</span><em>¥</em>
        </div>

        <button class="inviteBtn primary" id="btnShareWA" @click="shareInvite('whatsapp')">Share via WhatsApp</button>
        <button class="inviteBtn outline" id="btnShareTG" @click="shareInvite('telegram')">Share via Telegram</button>
      </div>

      <div class="inviteBottom">
        <div class="inviteTip">Invite these friends to withdraw faster</div>

        <div class="shareMore">
          <button class="miniShare" data-share="messenger" @click="shareInvite('messenger')">Messenger</button>
          <button class="miniShare" data-share="x" @click="shareInvite('x')">X</button>
          <button class="miniShare" data-share="discord" @click="shareInvite('discord')">Discord</button>
          <button class="miniShare" data-share="email" @click="shareInvite('email')">Email</button>
          <button class="miniShare" data-share="copy" @click="shareInvite('copy')">Copy Link</button>
        </div>
      </div>
    </div>
  </div>

  <div v-if="showExitModal" id="exitModal" class="modalMask show" aria-hidden="false">
    <div class="modalCard" role="dialog" aria-modal="true" aria-label="退出确认">
      <button class="modalClose" id="exitClose" aria-label="关闭" @click="closeExitModal">×</button>

      <div class="modalTitle">确定要退出吗？</div>
      <div class="modalSub">仅差<span id="exitNeed">{{ formatMoney(neededAmount) }}</span>元即可提现</div>

      <div class="modalMoneyBox">
        <div class="moneyLine">¥ <b id="exitMoney">{{ formatMoney(displayAmount) }}</b></div>

        <div class="progressWrap">
          <div class="progressBar">
            <div class="progressFill" id="exitProgress" :style="{ width: exitProgress.toFixed(0) + '%' }"></div>
          </div>
        </div>

        <div class="modalHint">仅差<span id="exitNeed2">{{ formatMoney(neededAmount) }}</span>元即可提现{{ formatMoney(targetAmount) }}元</div>
        <div class="soonTag">即将提现</div>
      </div>

      <button class="modalPrimary" id="exitStay" @click="closeExitModal">继续抽奖，加速提现</button>
      <button class="modalGhost" id="exitLeave" @click="closeExitModal">狠心离开</button>
    </div>
  </div>
</template>

<style scoped>
:global(body) {
  margin: 0;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, "PingFang SC",
    "Hiragino Sans GB", "Noto Sans CJK SC", "Microsoft YaHei", sans-serif;
  background: linear-gradient(180deg, #ff3b2f, #ffbd8e 45%, #fff2ed 100%);
  min-height: 100vh;
  color: #fff;
}

:global(*) {
  box-sizing: border-box;
}

.app {
  max-width: 430px;
  margin: 0 auto;
  min-height: 100vh;
  position: relative;
  padding-bottom: 24px;
}

.topbar {
  position: sticky;
  top: 0;
  z-index: 20;
  padding: 6px 10px 4px;
  background: linear-gradient(180deg, rgba(255, 70, 42, 0.98), rgba(255, 70, 42, 0.82));
  backdrop-filter: blur(8px);
}

.toprow {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.iconbtn {
  width: 32px;
  height: 32px;
  border-radius: 18px;
  display: grid;
  place-items: center;
  color: #fff;
  background: rgba(255, 255, 255, 0.14);
  border: 1px solid rgba(255, 255, 255, 0.22);
}

.title {
  flex: 1;
  text-align: center;
  font-weight: 800;
  font-size: 20px;
  letter-spacing: 0.5px;
}

.right {
  display: flex;
  gap: 8px;
}

.banner {
  margin-bottom: 4px;
  background: rgba(255, 255, 255, 0.18);
  border: 1px solid rgba(255, 255, 255, 0.28);
  border-radius: 22px;
  padding: 6px 8px;
  display: flex;
  align-items: center;
  gap: 10px;
}

.avatar {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: radial-gradient(circle at 35% 30%, #dff3ff 0 35%, #7ab7ff 36% 70%, #2b57ff 71% 100%);
  border: 2px solid rgba(255, 255, 255, 0.55);
  flex: 0 0 auto;
}

.banner .txt {
  font-size: 14px;
  color: #fff;
  opacity: 0.95;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.banner .amount {
  margin-left: auto;
  font-weight: 900;
  color: #ffe66d;
}

.content {
  padding: 8px 12px 0;
}

.card {
  background: #fff6b8;
  color: #1f1f1f;
  border-radius: 18px;
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.14);
  padding: 8px 10px 8px;
  position: relative;
  overflow: hidden;
}

.card::after {
  content: "";
  position: absolute;
  right: -18px;
  top: -8px;
  width: 120px;
  height: 120px;
  background: radial-gradient(circle at 30% 30%, rgba(255, 255, 255, 0.7), rgba(255, 210, 74, 0.1) 55%, rgba(255, 210, 74, 0) 75%);
  transform: rotate(12deg);
  pointer-events: none;
}

.cardhead {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  margin-bottom: 4px;
}

.name {
  display: flex;
  align-items: center;
  gap: 10px;
  font-weight: 800;
}

.pill {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  background: #ffeaa3;
  color: #b35c00;
  border: 1px solid rgba(179, 92, 0, 0.18);
  padding: 6px 10px;
  border-radius: 999px;
  font-size: 13px;
  font-weight: 800;
  white-space: nowrap;
}

.check {
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: #20c05c;
  display: inline-grid;
  place-items: center;
  color: #fff;
  font-size: 12px;
  line-height: 1;
}

.headline {
  font-size: 14px;
  font-weight: 900;
  margin: 3px 0 4px;
}

.headline .hot {
  color: #ff2f2f;
}

.bigmoney {
  font-size: 48px;
  font-weight: 1000;
  color: #ff2b2b;
  text-align: center;
  margin: 2px 0 4px;
  letter-spacing: 1px;
  line-height: 1.1;
}

.bigmoney span {
  font-size: 20px;
  font-weight: 900;
  margin-left: 4px;
}

.stats {
  background: #fff;
  border-radius: 14px;
  display: grid;
  grid-template-columns: 1fr 1fr;
  overflow: hidden;
  border: 1px solid rgba(0, 0, 0, 0.06);
}

.stat {
  padding: 6px 6px;
  text-align: center;
}

.stat:not(:last-child) {
  border-right: 1px solid rgba(0, 0, 0, 0.06);
}

.stat .k {
  font-size: 12px;
  color: #666;
  margin-bottom: 2px;
  font-weight: 700;
}

.stat .v {
  font-size: 22px;
  font-weight: 1000;
  color: #111;
}

.stat .v small {
  font-size: 14px;
  font-weight: 900;
  color: #333;
}

.stat .v.red {
  color: #ff2f2f;
}

.wish {
  margin: 4px 0 4px;
  background: rgba(255, 255, 255, 0.22);
  border: 1px solid rgba(255, 255, 255, 0.28);
  border-radius: 999px;
  text-align: center;
  padding: 4px 10px;
  font-weight: 900;
  color: #fff;
}

.wish b {
  color: #ffe66d;
}

.wheelWrap {
  margin-top: 4px;
  position: relative;
  display: flex;
  justify-content: center;
}

.wheelFrame {
  width: 300px;
  max-width: 92vw;
  aspect-ratio: 1 / 1;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.18);
  padding: 12px;
  box-shadow: 0 18px 30px rgba(0, 0, 0, 0.14);
  position: relative;
}

.wheel {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  position: relative;
  overflow: hidden;
  background: conic-gradient(
    #d0f5de 0 30deg,
    #ffefe1 30deg 60deg,
    #ffeeaf 60deg 90deg,
    #dbcdcd 90deg 120deg,
    #d0f5de 120deg 150deg,
    #ffefe1 150deg 180deg,
    #ffeeaf 180deg 210deg,
    #dbcdcd 210deg 240deg,
    #d0f5de 240deg 270deg,
    #ffefe1 270deg 300deg,
    #ffeeaf 300deg 330deg,
    #dbcdcd 330deg 360deg
  );
  border: 10px solid rgb(255 241 209);
  transform: rotate(var(--rot, 0deg));
  transition: transform 3.2s cubic-bezier(0.12, 0.9, 0.1, 1);
}

.prize {
  position: absolute;
  left: 50%;
  top: 50%;
  height: 0;
  transform: rotate(calc(var(--a) + 15deg)) translateY(-122px) translateX(-50%);
  transform-origin: 0 0;
  pointer-events: none;
}

.prize span {
  display: flex;
  align-items: center;
  gap: 6px;
  transform-origin: left center;
  font-size: 13px;
  font-weight: 900;
  white-space: nowrap;
  color: #7a3b00;
  text-shadow: 0 1px 0 rgba(255, 255, 255, 0.55);
}

.prize span.gray {
  color: #fff;
  text-shadow: none;
}

.centerBtn {
  position: absolute;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
  width: 118px;
  height: 118px;
  border-radius: 50%;
  background: radial-gradient(circle at 30% 25%, #ff9b6a, #ff3b2f 60%, #e00000 100%);
  box-shadow: 0 16px 30px rgba(0, 0, 0, 0.25);
  border: 8px solid #ffd24a;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  z-index: 1;
}

.centerBtn .tag {
  font-size: 11px;
  font-weight: 900;
  background: rgb(255 67 67);
  padding: 5px 8px;
  border-radius: 999px;
  border: 1px solid rgba(255, 255, 255, 0.25);
  margin-bottom: 6px;
  position: absolute;
  top: -10px;
  z-index: 2;
}

.centerBtn .main {
  font-size: 30px;
  font-weight: 1000;
  letter-spacing: 2px;
  margin-top: -2px;
}

.centerBtn .sub {
  font-size: 11px;
  opacity: 0.95;
  margin-top: 4px;
  font-weight: 800;
}

.pointer {
  position: absolute;
  top: -2px;
  left: 50%;
  top: calc(50% - 74px);
  z-index: 1;
  transform: translateX(-50%);
  width: 0;
  height: 0;
  border-left: 14px solid transparent;
  border-right: 14px solid transparent;
  border-bottom: 22px solid #ffdf7d;
  filter: drop-shadow(0 6px 8px rgba(0, 0, 0, 0.16));
}

.toast {
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
  bottom: 52px;
  width: min(92vw, 360px);
  background: rgba(0, 0, 0, 0.45);
  padding: 10px 12px;
  border-radius: 14px;
  border: 1px solid rgba(255, 255, 255, 0.18);
  backdrop-filter: blur(6px);
  font-weight: 900;
  text-align: center;
  z-index: 3;
}

.toast small {
  display: block;
  font-weight: 700;
  opacity: 0.85;
  margin-top: 6px;
  line-height: 1.25;
}

.bottomArea {
  margin-top: 10px;
  padding: 0 14px;
  position: relative;
}

.cta {
  width: 100%;
  border: none;
  border-radius: 999px;
  padding: 12px 12px;
  background: linear-gradient(180deg, #ff2f2f, #ff5c2e);
  color: #fff;
  font-size: 16px;
  font-weight: 1000;
  box-shadow: 0 14px 22px rgba(0, 0, 0, 0.16);
  cursor: pointer;
}

.countdown {
  text-align: center;
  margin-top: 6px;
  color: rgba(255, 255, 255, 0.85);
  font-weight: 800;
  font-size: 14px;
}

.scan {
  position: absolute;
  right: 14px;
  bottom: 68px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.scan .btn {
  width: 58px;
  height: 58px;
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.18);
  border: 1px solid rgba(255, 255, 255, 0.28);
  display: grid;
  place-items: center;
  box-shadow: 0 12px 20px rgba(0, 0, 0, 0.12);
}

.qr {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  background:
    linear-gradient(90deg, #fff 0 40%, transparent 40% 60%, #fff 60% 100%),
    linear-gradient(#fff 0 40%, transparent 40% 60%, #fff 60% 100%);
  background-size: 12px 12px;
  background-position: 0 0;
  opacity: 0.95;
}

.scan .label {
  font-size: 12px;
  font-weight: 900;
  background: rgba(255, 255, 255, 0.18);
  border: 1px solid rgba(255, 255, 255, 0.24);
  padding: 6px 10px;
  border-radius: 999px;
}

.mini {
  position: absolute;
  right: 10px;
  top: 12px;
  transform: translateY(0);
  z-index: 8;
  display: flex;
  align-items: center;
  gap: 8px;
}

.mini .chip {
  background: #fff;
  color: #ff3b2f;
  font-weight: 1000;
  font-size: 12px;
  padding: 6px 10px;
  border-radius: 999px;
  border: 1px solid rgba(0, 0, 0, 0.06);
  box-shadow: 0 10px 18px rgba(0, 0, 0, 0.1);
}

.panel {
  margin-top: 8px;
  background: #fff;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(0, 0, 0, 0.06);
}

.tabs {
  display: flex;
  background: #fff;
}

.tab {
  flex: 1;
  padding: 10px 8px;
  font-size: 16px;
  font-weight: 1000;
  border: none;
  background: transparent;
  color: #b06a2a;
  position: relative;
  cursor: pointer;
}

.tab.active {
  color: #ff3b2f;
}

.tab.active::after {
  content: "";
  position: absolute;
  left: 22%;
  right: 22%;
  bottom: 0;
  height: 3px;
  border-radius: 3px;
  background: #ff3b2f;
}

.tabBody {
  padding: 10px 12px 12px;
}

.tabPage {
  display: none;
}

.tabPage.active {
  display: block;
}

.receiptItem {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 8px;
  border-radius: 14px;
  background: #fffaf0;
  border: 1px solid rgba(0, 0, 0, 0.06);
}

.receiptItem + .receiptItem {
  margin-top: 10px;
}

.uavatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  flex: 0 0 auto;
  border: 2px solid rgba(255, 255, 255, 0.8);
  box-shadow: 0 8px 14px rgba(0, 0, 0, 0.08);
  background: radial-gradient(circle at 30% 30%, #ffe7d6, #ff9a7e 55%, #ff3b2f 100%);
}

.uavatar.a2 {
  background: radial-gradient(circle at 30% 30%, #dff3ff, #7ab7ff 55%, #2b57ff 100%);
}

.uavatar.a3 {
  background: radial-gradient(circle at 30% 30%, #e8ffe3, #5be49b 55%, #0bbd5c 100%);
}

.info {
  flex: 1;
  min-width: 0;
}

.row1 {
  display: flex;
  align-items: center;
  gap: 10px;
}

.uname {
  font-weight: 1000;
  color: #7a3b00;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 52%;
}

.date {
  margin-left: auto;
  font-weight: 900;
  font-size: 12px;
  color: #b06a2a;
  background: #ffeaa3;
  border: 1px solid rgba(179, 92, 0, 0.18);
  padding: 4px 8px;
  border-radius: 999px;
  white-space: nowrap;
}

.desc {
  margin-top: 6px;
  font-size: 13px;
  font-weight: 800;
  color: #6b3a14;
  opacity: 0.92;
}

.badge {
  width: 72px;
  border-radius: 14px;
  border: 2px solid rgba(18, 185, 90, 0.22);
  background: linear-gradient(180deg, #eafff2, #ffffff);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 8px 6px;
  flex: 0 0 auto;
}

.money {
  font-weight: 1000;
  color: #12b95a;
  font-size: 22px;
  line-height: 1;
}

.money span {
  font-size: 12px;
  font-weight: 900;
  margin-left: 2px;
}

.status {
  margin-top: 6px;
  font-size: 12px;
  font-weight: 1000;
  color: #12b95a;
  background: rgba(18, 185, 90, 0.12);
  border: 1px solid rgba(18, 185, 90, 0.18);
  padding: 4px 8px;
  border-radius: 999px;
}

.guideCard {
  background: #fffaf0;
  border: 1px solid rgba(0, 0, 0, 0.06);
  border-radius: 14px;
  padding: 12px 12px 10px;
}

.gtitle {
  font-size: 16px;
  font-weight: 1000;
  color: #7a3b00;
  margin-bottom: 8px;
}

.glist {
  margin: 0;
  padding-left: 18px;
  color: #6b3a14;
  font-weight: 800;
  line-height: 1.6;
  font-size: 14px;
}

.gtip {
  margin-top: 10px;
  font-size: 12px;
  font-weight: 900;
  color: #ff3b2f;
  background: rgba(255, 59, 47, 0.08);
  border: 1px solid rgba(255, 59, 47, 0.12);
  padding: 8px 10px;
  border-radius: 12px;
}

.modalMask {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.55);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  padding: 18px;
}

.modalCard {
  width: min(92vw, 360px);
  background: #fff6cf;
  border-radius: 18px;
  box-shadow: 0 18px 40px rgba(0, 0, 0, 0.25);
  position: relative;
  padding: 18px 16px 16px;
}

.modalClose {
  position: absolute;
  right: -10px;
  top: -10px;
  width: 34px;
  height: 34px;
  border-radius: 17px;
  border: none;
  background: rgba(0, 0, 0, 0.829);
  color: #fff;
  font-size: 20px;
  line-height: 34px;
  cursor: pointer;
}

.modalTitle {
  text-align: center;
  font-size: 18px;
  font-weight: 1000;
  color: #b44f00;
  margin-top: 8px;
}

.modalSub {
  text-align: center;
  font-size: 16px;
  font-weight: 1000;
  color: #b44f00;
  margin-top: 6px;
}

.modalSub span {
  color: #ff3b2f;
}

.modalMoneyBox {
  margin: 14px auto 14px;
  background: #fff;
  border-radius: 14px;
  padding: 14px 14px 12px;
  border: 1px solid rgba(0, 0, 0, 0.06);
  position: relative;
}

.moneyLine {
  text-align: center;
  color: #ff3b2f;
  font-weight: 1000;
  font-size: 16px;
}

.moneyLine b {
  font-size: 34px;
}

.soonTag {
  position: absolute;
  right: 12px;
  top: 12px;
  background: #ff7a2e;
  color: #fff;
  font-weight: 1000;
  font-size: 12px;
  padding: 4px 8px;
  border-radius: 999px;
}

.progressWrap {
  margin-top: 10px;
}

.progressBar {
  height: 10px;
  background: #ffe0d0;
  border-radius: 999px;
  overflow: hidden;
}

.progressFill {
  height: 100%;
  background: linear-gradient(90deg, #ff3b2f, #ff8a2a);
  border-radius: 999px;
}

.modalHint {
  text-align: center;
  margin-top: 10px;
  font-weight: 900;
  font-size: 13px;
  color: #b44f00;
}

.modalHint span {
  color: #ff3b2f;
}

.modalPrimary {
  width: 100%;
  border: none;
  cursor: pointer;
  border-radius: 999px;
  padding: 14px 12px;
  background: linear-gradient(180deg, #ff8a2a, #ff3b2f);
  color: #fff;
  font-size: 16px;
  font-weight: 1000;
  box-shadow: 0 12px 20px rgba(0, 0, 0, 0.16);
}

.modalGhost {
  width: 100%;
  margin-top: 10px;
  border: none;
  background: transparent;
  color: #b44f00;
  font-weight: 1000;
  font-size: 14px;
  cursor: pointer;
}

.ctaWrap {
  position: relative;
}

.cta.ctaPulse {
  position: relative;
  animation: cta-breathe 1.05s ease-in-out infinite;
  transform-origin: center;
}

@keyframes cta-breathe {
  0% {
    transform: scale(1);
    filter: brightness(1);
  }
  50% {
    transform: scale(1.05);
    filter: brightness(1.06);
  }
  100% {
    transform: scale(1);
    filter: brightness(1);
  }
}

.handHint {
  position: absolute;
  right: 90px;
  top: 46px;
  font-size: 40px;
  line-height: 1;
  pointer-events: none;
  transform: rotate(-20deg);
  animation: hand-move 0.6s ease-in-out infinite;
  filter: drop-shadow(0 8px 10px rgba(0, 0, 0, 0.2));
}

@keyframes hand-move {
  0% {
    transform: translate(0, 0) rotate(-20deg) scale(1);
  }
  50% {
    transform: translate(-10px, -8px) rotate(-20deg) scale(0.98);
  }
  100% {
    transform: translate(0, 0) rotate(-20deg) scale(1);
  }
}

.inviteCard {
  width: min(92vw, 380px);
  background: #fff6d6;
  border-radius: 18px;
  box-shadow: 0 18px 44px rgba(0, 0, 0, 0.28);
  position: relative;
  padding: 18px 16px 14px;
}

.inviteMain {
  background: #fffdf6;
  border: 1px solid rgba(0, 0, 0, 0.06);
  border-radius: 16px;
  padding: 18px 14px 14px;
  text-align: center;
  position: relative;
}

.inviteAvatar {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  margin: -44px auto 6px;
  border: 4px solid #fff;
  background: radial-gradient(circle at 35% 30%, #dff3ff 0 35%, #7ab7ff 36% 70%, #2b57ff 71% 100%);
  box-shadow: 0 10px 18px rgba(0, 0, 0, 0.15);
}

.inviteName {
  display: inline-block;
  margin: 2px auto 10px;
  padding: 6px 18px;
  border-radius: 999px;
  background: linear-gradient(90deg, #f8d7b3, #f4b67f);
  color: #d85a00;
  font-weight: 1000;
  border: 1px solid rgba(0, 0, 0, 0.06);
}

.inviteText {
  font-size: 18px;
  font-weight: 1000;
  color: #c05b00;
  line-height: 1.35;
}

.inviteText .green {
  color: #0bbd5c;
}

.inviteBig {
  margin: 14px 0 14px;
  font-size: 74px;
  font-weight: 1000;
  color: #ff3b2f;
  letter-spacing: 1px;
}

.inviteBig em {
  font-style: normal;
  font-size: 26px;
  font-weight: 1000;
  margin-left: 6px;
  color: #ff3b2f;
}

.inviteBtn {
  width: 100%;
  border: none;
  border-radius: 999px;
  padding: 14px 12px;
  font-weight: 1000;
  font-size: 16px;
  cursor: pointer;
}

.inviteBtn.primary {
  background: linear-gradient(180deg, #2dd071, #0bbd5c);
  color: #fff;
  box-shadow: 0 12px 22px rgba(0, 0, 0, 0.14);
}

.inviteBtn.outline {
  margin-top: 10px;
  background: transparent;
  border: 2px solid rgba(240, 140, 40, 0.55);
  color: #d86a00;
}

.inviteBottom {
  margin-top: 12px;
  background: rgba(0, 0, 0, 0.55);
  border-radius: 14px;
  padding: 10px 12px 12px;
  color: #f6f2e6;
  position: relative;
}

.inviteBottom::before {
  content: "";
  position: absolute;
  left: 50%;
  top: -10px;
  transform: translateX(-50%);
  width: 0;
  height: 0;
  border-left: 10px solid transparent;
  border-right: 10px solid transparent;
  border-bottom: 10px solid rgba(0, 0, 0, 0.55);
}

.inviteTip {
  text-align: center;
  font-weight: 1000;
  margin-bottom: 10px;
  opacity: 0.95;
}

.shareMore {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  justify-content: center;
}

.miniShare {
  border: none;
  cursor: pointer;
  font-weight: 1000;
  font-size: 12px;
  padding: 8px 10px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.14);
  color: #fff;
  border: 1px solid rgba(255, 255, 255, 0.22);
}

.status-card {
  margin-top: 12px;
  background: rgba(255, 255, 255, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 12px;
  padding: 10px 12px;
  font-weight: 800;
}

.status-card.error {
  color: #ffefe5;
}

.status-card.hint {
  color: #fff8d6;
}

.empty {
  text-align: center;
  padding: 16px 10px;
  font-weight: 800;
  color: #7a3b00;
}

@media (max-width: 360px) {
  .bigmoney {
    font-size: 42px;
  }

  .centerBtn {
    width: 110px;
    height: 110px;
  }

  .wheelFrame {
    padding: 10px;
  }

  .handHint {
    right: 62px;
    top: -12px;
    transform: rotate(-20deg) scale(0.92);
  }
}
</style>



