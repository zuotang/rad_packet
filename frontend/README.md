# Frontend README

## 技术栈

- Vue `3`
- Pinia
- Vue Router
- Axios
- Vite `5`

## 目录说明

- `src/main.js`：应用入口
- `src/router`：路由和登录态守卫
- `src/api/client.js`：axios 封装（token 注入、401 自动登出）
- `src/stores`：Pinia 状态管理
- `src/views`：页面实现

## 页面清单

- `/login`：登录/注册
- `/r/:code`：邀请落地页入口（写入 `invite_code`）
- `/`：首页（奖励汇总 + 邀请进度）
- `/`：首页支持“解冻待解锁奖励”（调用 `POST /api/reward/unlock`）
- `/tasks`：任务中心（领取任务）
- `/invite`：邀请页（链接 + 分享按钮）
- `/wallet`：钱包页（余额 + 提现）
- `/records`：流水记录

说明：后台管理端已独立为单独项目 `../admin_frontend`。

## 与后端联调

前端请求统一走 `/api`，开发环境由 Vite 代理到后端。

- 默认代理目标：`http://localhost:8381`
- 代理配置文件：`vite.config.js`

可通过环境变量覆盖代理目标：

```bash
VITE_PROXY_TARGET=http://localhost:8381
```

如果你不想走代理，也可以直接指定 API 基地址：

```bash
VITE_API_BASE=http://localhost:8381/api
```

## 启动方式

```bash
cd frontend
npm install
npm run dev
```

默认开发地址：

- `http://localhost:5173`

## 构建

```bash
npm run build
```

## 邀请链路说明

- 用户访问 `/r/:code` 时，前端把 `code` 存入 `localStorage(invite_code)`
- 登录成功后自动尝试调用 `POST /api/referral/bind`
- 绑定失败不会阻断登录主流程

## 常见问题

- 报错：`getActivePinia() was called but there was no active Pinia`  
  原因是路由守卫/axios 拦截器中调用 store 时机早于 `app.use(pinia)`。  
  当前代码已通过共享 pinia 实例修复：`src/stores/pinia.js`。

- 后台接口返回 `ADMIN_UNAUTHORIZED`  
  需要先在 `/admin` 页面填写并保存 `X-Admin-Key`（与后端 `APP_ADMIN_KEY` 一致）。
