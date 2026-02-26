# Admin Frontend

独立后台管理前端项目（与业务前端分离）。

## 技术栈

- Vue 3
- Vite 5
- Pinia
- Axios

## 启动

```bash
cd admin_frontend
npm install
npm run dev
```

默认地址：`http://localhost:5174`

## 联调配置

- 默认代理：`/api -> http://localhost:8381`
- 管理接口基地址：`/api/admin`
- 可通过 `.env.development` 调整 `VITE_PROXY_TARGET`

## 登录说明

- 打开 `/login`
- 输入后端配置的 `APP_ADMIN_KEY`
- 保存后进入后台页面

## 功能模块

- 运营看板
- 提现审核
- 任务配置（增删改）
- 系统配置（key/value）
- 风控标记与黑名单
