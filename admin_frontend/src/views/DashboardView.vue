<script setup>
import { onMounted, reactive, ref } from "vue";
import api from "../api/client";

const tab = ref("dashboard");
const loading = ref(false);
const error = ref("");
const hint = ref("");

const dashboard = reactive({
  total_users: 0,
  new_users_today: 0,
  pending_withdraws: 0,
  pending_withdraw_amount: 0,
  rewards_pending_amount: 0,
  rewards_unlocked_amount: 0,
});

const withdrawItems = ref([]);
const taskItems = ref([]);
const configItems = ref([]);
const riskFlags = ref([]);
const blacklists = ref([]);

const reviewForm = reactive({ request_id: 0, status: "approved", note: "" });
const taskForm = reactive({ id: 0, type: "custom", name: "", reward_rule_id: "", reward_amount: 0.2, enabled: true, country_scope: "*" });
const configForm = reactive({ key: "", value: "" });
const riskForm = reactive({ user_id: 0, reason: "", score: 20 });
const blacklistForm = reactive({ type: "ip", value: "", note: "" });

async function loadAll() {
  loading.value = true;
  error.value = "";
  hint.value = "";
  try {
    const [dash, withdraw, tasks, configs, flags, bl] = await Promise.all([
      api.get("/dashboard"),
      api.get("/withdraw/list?page=1&size=20"),
      api.get("/task/list"),
      api.get("/config/list"),
      api.get("/risk/flags?page=1&size=20"),
      api.get("/blacklist/list?page=1&size=20"),
    ]);
    Object.assign(dashboard, dash.data);
    withdrawItems.value = withdraw.data.items || [];
    taskItems.value = tasks.data.items || [];
    configItems.value = configs.data.items || [];
    riskFlags.value = flags.data.items || [];
    blacklists.value = bl.data.items || [];
  } catch (err) {
    error.value = err?.response?.data?.message || "后台数据加载失败，请检查管理员密钥";
  } finally {
    loading.value = false;
  }
}

onMounted(loadAll);

async function reviewWithdraw() {
  try {
    await api.post("/withdraw/review", reviewForm);
    hint.value = "提现审核已提交";
    await loadAll();
  } catch (err) {
    error.value = err?.response?.data?.message || "审核失败";
  }
}

async function saveTask() {
  try {
    await api.post("/task/save", taskForm);
    hint.value = "任务已保存";
    Object.assign(taskForm, { id: 0, type: "custom", name: "", reward_rule_id: "", reward_amount: 0.2, enabled: true, country_scope: "*" });
    await loadAll();
  } catch (err) {
    error.value = err?.response?.data?.message || "任务保存失败";
  }
}

function editTask(t) {
  Object.assign(taskForm, {
    id: t.id || t.ID,
    type: t.type || t.Type,
    name: t.name || t.Name,
    reward_rule_id: t.reward_rule_id || t.RewardRuleID,
    reward_amount: Number(t.reward_amount ?? t.RewardAmount ?? 0),
    enabled: Boolean(t.enabled ?? t.Enabled),
    country_scope: t.country_scope || t.CountryScope || "*",
  });
  tab.value = "tasks";
}

async function deleteTask(id) {
  try {
    await api.delete(`/task/${id}`);
    hint.value = "任务已删除";
    await loadAll();
  } catch (err) {
    error.value = err?.response?.data?.message || "任务删除失败";
  }
}

async function upsertConfig() {
  try {
    await api.post("/config/upsert", configForm);
    hint.value = "配置已保存";
    configForm.key = "";
    configForm.value = "";
    await loadAll();
  } catch (err) {
    error.value = err?.response?.data?.message || "配置保存失败";
  }
}

async function addRiskFlag() {
  try {
    await api.post("/risk/flag/add", riskForm);
    hint.value = "风控标记已添加";
    riskForm.user_id = 0;
    riskForm.reason = "";
    await loadAll();
  } catch (err) {
    error.value = err?.response?.data?.message || "风控标记添加失败";
  }
}

async function addBlacklist() {
  try {
    await api.post("/blacklist/add", blacklistForm);
    hint.value = "黑名单已添加";
    blacklistForm.value = "";
    blacklistForm.note = "";
    await loadAll();
  } catch (err) {
    error.value = err?.response?.data?.message || "黑名单添加失败";
  }
}
</script>

<template>
  <section class="card">
    <h2>后台管理</h2>
    <p class="muted">提现审核、任务配置、奖励参数与风控管理</p>
  </section>

  <section class="card">
    <div class="tabs">
      <button :class="['tab-btn', tab==='dashboard'?'active':'']" @click="tab='dashboard'">概览</button>
      <button :class="['tab-btn', tab==='withdraw'?'active':'']" @click="tab='withdraw'">提现审核</button>
      <button :class="['tab-btn', tab==='tasks'?'active':'']" @click="tab='tasks'">任务配置</button>
      <button :class="['tab-btn', tab==='configs'?'active':'']" @click="tab='configs'">系统配置</button>
      <button :class="['tab-btn', tab==='risk'?'active':'']" @click="tab='risk'">风控</button>
      <button class="btn secondary" :disabled="loading" @click="loadAll">刷新</button>
    </div>

    <div v-if="tab==='dashboard'" class="grid2">
      <div class="list-item"><div class="muted">总用户</div><div class="kpi">{{ dashboard.total_users }}</div></div>
      <div class="list-item"><div class="muted">今日新增</div><div class="kpi">{{ dashboard.new_users_today }}</div></div>
      <div class="list-item"><div class="muted">待审提现笔数</div><div class="kpi">{{ dashboard.pending_withdraws }}</div></div>
      <div class="list-item"><div class="muted">待审提现金额</div><div class="kpi">{{ Number(dashboard.pending_withdraw_amount).toFixed(2) }}</div></div>
      <div class="list-item"><div class="muted">待解锁奖励总额</div><div class="kpi">{{ Number(dashboard.rewards_pending_amount).toFixed(2) }}</div></div>
      <div class="list-item"><div class="muted">已解锁奖励总额</div><div class="kpi">{{ Number(dashboard.rewards_unlocked_amount).toFixed(2) }}</div></div>
    </div>

    <div v-else-if="tab==='withdraw'" class="list">
      <div class="list-item row">
        <input v-model.number="reviewForm.request_id" type="number" placeholder="request_id" />
        <select v-model="reviewForm.status"><option>approved</option><option>rejected</option><option>paid</option></select>
        <input v-model="reviewForm.note" placeholder="note" />
        <button @click="reviewWithdraw">提交审核</button>
      </div>
      <div class="list-item" v-for="w in withdrawItems" :key="w.id || w.ID">
        <div class="row"><strong>#{{ w.id || w.ID }}</strong><span class="badge">{{ w.status || w.Status }}</span></div>
        <p class="muted">user: {{ w.user_id || w.UserID }} | amount: {{ Number(w.amount || w.Amount || 0).toFixed(2) }}</p>
      </div>
    </div>

    <div v-else-if="tab==='tasks'" class="list">
      <div class="list-item form">
        <div class="row"><input v-model="taskForm.name" placeholder="任务名" /><input v-model="taskForm.type" placeholder="type" /></div>
        <div class="row"><input v-model="taskForm.reward_rule_id" placeholder="reward_rule_id" /><input v-model.number="taskForm.reward_amount" type="number" step="0.01" /></div>
        <div class="row"><input v-model="taskForm.country_scope" placeholder="country_scope" /><label class="muted"><input v-model="taskForm.enabled" type="checkbox" /> enabled</label><button @click="saveTask">保存任务</button></div>
      </div>
      <div class="list-item" v-for="t in taskItems" :key="t.id || t.ID">
        <div class="row"><strong>{{ t.name || t.Name }}</strong><span class="badge">{{ (t.enabled ?? t.Enabled) ? 'enabled' : 'disabled' }}</span></div>
        <p class="muted">id: {{ t.id || t.ID }} | reward: {{ Number(t.reward_amount ?? t.RewardAmount ?? 0).toFixed(2) }}</p>
        <div class="row"><button class="btn secondary" @click="editTask(t)">编辑</button><button @click="deleteTask(t.id || t.ID)">删除</button></div>
      </div>
    </div>

    <div v-else-if="tab==='configs'" class="list">
      <div class="list-item row"><input v-model="configForm.key" placeholder="key" /><input v-model="configForm.value" placeholder="value" /><button @click="upsertConfig">保存配置</button></div>
      <div class="list-item" v-for="c in configItems" :key="c.id || c.ID"><div class="row"><strong>{{ c.key || c.Key }}</strong><span class="badge">id {{ c.id || c.ID }}</span></div><p class="muted">{{ c.value || c.Value }}</p></div>
    </div>

    <div v-else class="list">
      <div class="list-item row"><input v-model.number="riskForm.user_id" type="number" placeholder="user_id" /><input v-model.number="riskForm.score" type="number" placeholder="score" /><input v-model="riskForm.reason" placeholder="reason" /><button @click="addRiskFlag">添加标记</button></div>
      <div class="list-item row"><input v-model="blacklistForm.type" placeholder="type" /><input v-model="blacklistForm.value" placeholder="value" /><input v-model="blacklistForm.note" placeholder="note" /><button @click="addBlacklist">添加黑名单</button></div>
      <div class="list-item"><h3>风险标记</h3><div class="list-item" v-for="f in riskFlags" :key="f.id || f.ID"><p class="muted">user {{ f.user_id || f.UserID }} | score {{ f.score || f.Score }} | {{ f.reason || f.Reason }}</p></div></div>
      <div class="list-item"><h3>黑名单</h3><div class="list-item" v-for="b in blacklists" :key="b.id || b.ID"><p class="muted">{{ b.type || b.Type }}: {{ b.value || b.Value }} ({{ b.note || b.Note || '-' }})</p></div></div>
    </div>
  </section>

  <section class="card" v-if="hint"><p class="muted">{{ hint }}</p></section>
  <section class="card" v-if="error"><p class="error">{{ error }}</p></section>
</template>
