**整套 Vue3+Pinia+纯JS + Go(Echo)+GORM+MySQL**的开发建议和“提示词”

**高裂变、高转化、但规则透明可落地**的增长平台工程方案。

> 工程落地文档：
> - 后端运行与部署：`backend/README.md`
> - 前端联调与代理：`frontend/README.md`
> - 项目总览与启动：`README_PROJECT.md`

---

## 1) 产品机制建议（能猛、也能做得住）

### 裂变主循环（建议你只做 1 个主循环）

* 新用户注册 → 得“待解锁奖励”(pending)
* 完成任务/邀请 → 解锁一部分
* 达到门槛 → 可提现（或兑换）
* 提现后继续给“下一档奖励目标”

**关键**：把增长动力从“返钱”转移到“任务奖励 + 阶梯解锁”，避免陷入无限返利模型。

### 推荐奖励层级

* 建议 **最多 2 级**（1级主奖励，2级弱激励），更易控成本/控风控。
* 每笔奖励都要有：`来源事件(邀请/任务/广告)` + `可回溯订单号/事件ID`。

---

## 2) 后端 Go + Echo + GORM + MySQL：模块拆分（你照这个建仓）

### 服务模块

1. **auth**：手机号/邮箱/第三方登录、设备绑定、JWT
2. **user**：用户资料、地区/语言、钱包状态
3. **referral**：邀请码、上下级绑定、归因、反作弊
4. **reward**：奖励计算、发放、解锁、过期
5. **wallet**：余额、冻结、提现申请、流水
6. **task**：任务中心（观看广告、下载、签到、分享等）
7. **risk**：风控规则、设备/IP、黑名单、限频
8. **admin**：运营后台（配置奖励档位、任务、提现审核）

### 工程关键点（非常重要）

* **幂等**：奖励发放必须支持幂等（同一个事件只发一次）
* **事务**：发奖励 = 写 reward + 写 wallet_ledger + 更新 wallet，必须事务一致
* **延迟可控**：提现走 `withdraw_request` 状态机（pending/approved/rejected/paid）
* **配置中心**：奖励档位不要写死在代码里，放 DB 或配置文件+热更新

---

## 3) MySQL 表设计（最小可跑版本）

### 用户与邀请关系

* `users(id, country, language, created_at, ...)`
* `referral_codes(user_id, code, created_at)`
* `referral_edges(id, parent_user_id, child_user_id, level, created_at, UNIQUE(child_user_id))`

  * **child 只能绑定一次**（防止归因混乱）
  * level 只存 1/2

### 奖励与钱包

* `wallets(user_id, balance, frozen, updated_at)`
* `wallet_ledgers(id, user_id, amount, type, ref_type, ref_id, created_at, UNIQUE(ref_type, ref_id, user_id))`
* `rewards(id, user_id, status, amount, unlock_amount, expire_at, source_type, source_id, created_at)`

  * status: pending/unlocked/expired

### 任务与事件

* `tasks(id, type, reward_rule_id, enabled, country_scope, ...)`
* `user_task_events(id, user_id, task_id, event_key, meta_json, created_at, UNIQUE(user_id, event_key))`

  * event_key：比如广告回调ID、下载归因ID

### 风控

* `device_fingerprints(id, user_id, device_hash, first_ip, last_ip, created_at, UNIQUE(device_hash))`
* `risk_flags(id, user_id, reason, score, created_at)`
* `blacklists(type, value, note, created_at)`（ip/device_hash/phone/email）

---

## 4) API 设计建议（Echo 路由清单）

* `POST /api/auth/login` `POST /api/auth/otp`
* `POST /api/referral/bind`（填邀请码绑定上级，一次性）
* `GET  /api/referral/status`（我邀请了谁/进度）
* `GET  /api/reward/summary`（待解锁/可用/已过期）
* `POST /api/task/claim`（完成任务领奖，带 event_key 幂等）
* `GET  /api/wallet`（余额+流水分页）
* `POST /api/withdraw/apply`
* `GET  /api/config/bootstrap`（把奖励档位、任务列表、地区文案一次性下发）

---

## 5) 前端 Vue3 + Pinia + 纯JS：页面与状态拆分

### 页面（最小闭环）

1. 登录/注册（含邀请码预填）
2. 首页（奖励进度条 + 立即邀请按钮）
3. 任务中心（列表 + 领取）
4. 邀请页（邀请码 + 海报图 + WhatsApp/Facebook 分享）
5. 钱包页（余额 + 提现）
6. 提现记录/奖励记录

### Pinia store 建议

* `useAuthStore`：token、用户信息
* `useConfigStore`：bootstrap 配置（奖励档位/任务/文案/国家）
* `useReferralStore`：邀请进度、邀请列表
* `useRewardStore`：奖励汇总、解锁进度
* `useWalletStore`：余额、流水、提现状态

### 分享落地链路（裂变关键）

* 入口链接：`https://h5.xxx.com/r/<code>`
* H5 打开后：

  * 存 `code` 到 localStorage
  * 注册成功后调用 `/referral/bind`
  * 已登录用户访问 `/r/<code>` 要提示“你已绑定过上级/或不能绑定自己”

---

## 6) 风控与反刷（不做这块你会被刷爆）

最低要做这些（都不违法、但很“硬”）：

* **同设备注册限制**（device_hash）
* **同 IP 限频**（短期阈值）
* **提现前加强校验**：KYC/手机号实名/钱包一致性（至少一项）
* **奖励延迟释放**：新用户奖励先 frozen，满足条件/过风控后再解冻
* **邀请有效性**：被邀请人需要完成“首个有效行为”（比如完成1个任务）才算有效邀请

---

## 7) 你可以直接用的“提示词”（给 ChatGPT/Claude 写代码用）

下面这些是“工程向提示词”，你复制就能让模型按你的栈输出可用代码：

### A. 后端工程骨架

> 用 Go 1.22 + Echo + GORM + MySQL 写一个 REST API 骨架。要求：JWT 鉴权中间件、统一错误码、请求日志、配置读取(viper)、数据库迁移(用 golang-migrate 或 gorm auto)。输出项目目录结构和关键文件代码。

### B. 邀请绑定（幂等+防自绑）

> 实现 POST /api/referral/bind：入参 code。规则：用户只能绑定一次上级；不能绑定自己；写 referral_edges(level=1)，并可选写 level=2；所有写入必须事务；并发下也不能重复绑定（MySQL 唯一键处理）。给出 handler/service/repo 代码。

### C. 奖励发放（幂等+账本）

> 设计奖励发放函数 GrantReward(userId, amount, refType, refId)。要求：写 rewards 和 wallet_ledgers，并更新 wallets.balance；幂等键为 (user_id, refType, refId)；必须事务；返回发放结果。给出 GORM model 和 service 实现。

### D. 任务领奖（事件幂等）

> 实现 POST /api/task/claim：入参 task_id、event_key、meta_json。要求：同一 event_key 同一用户只能领取一次；根据 task 配置计算奖励；调用 GrantReward；返回更新后的 reward summary。

### E. 前端 Pinia + API 封装

> 用 Vue3 + Pinia + Vite + 纯JS 实现 apiClient(axios)（带 token、401 自动登出），实现 useAuthStore/useReferralStore/useRewardStore/useWalletStore。给出每个 store 的 state/actions 和示例页面调用方式。

---

## 8) 一句话建议（落地最重要的取舍）

* **先做“配置中心 + 幂等发奖 + 钱包账本 + 邀请绑定”**四件事，其它都可以迭代。
* 不要把奖励逻辑散落在 handler 里，全部收敛到 `RewardService`，后面改档位不会炸。

