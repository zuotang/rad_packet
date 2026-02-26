# Red Packet Project

## 目录结构

- `backend`: Go 1.22 + Echo + GORM + MySQL API
- `frontend`: Vue3 + Pinia + Vite + JS H5
- `admin_frontend`: 独立后台前端（Vue3 + Vite）
- `docker-compose.yml`: 仅用于启动 MySQL

## 文档导航

- 后端详细说明：`backend/README.md`
- 前端详细说明：`frontend/README.md`
- 后台前端说明：`admin_frontend/README.md`

## 已实现能力

- 鉴权：`POST /api/auth/login`、`POST /api/auth/otp`、JWT 中间件
- 邀请：`POST /api/referral/bind`（防自绑、并发幂等、1/2级关系）+ `GET /api/referral/status`
- 奖励：`RewardService.GrantReward`（事务 + 钱包账本幂等）+ `GET /api/reward/summary`
- 任务：`POST /api/task/claim`（事件幂等 + 发奖 + 返回 summary）
- 钱包：`GET /api/wallet`（余额 + 流水分页）
- 提现：`POST /api/withdraw/apply`（pending 状态机入口 + 冻结余额）
- 奖励解冻：`POST /api/reward/unlock`（风控通过后把 pending 奖励释放到可用余额）
- 管理审核：`POST /api/admin/withdraw/review`（`X-Admin-Key`，支持 approved/rejected/paid）
- 配置：`GET /api/config/bootstrap`（任务+奖励档位配置）

## 启动后端

1. 准备 MySQL，创建库：`red_packet`
2. 修改 `backend/configs/config.yaml` 中 `mysql.dsn` 与 `jwt.secret`
3. 启动：

```bash
cd backend
go mod tidy
go run ./cmd/server
```

## 启动前端

```bash
cd frontend
npm install
npm run dev
```

默认会请求 `http://localhost:8080/api`。可通过 `VITE_API_BASE` 覆盖。
开发环境默认通过 Vite 代理请求 `/api`，代理目标默认 `http://localhost:8381`。

## Docker Compose 启动后端

在项目根目录执行：

```bash
docker compose up -d
```

说明：

- Compose 只启动 `mysql`（3306）
- Go 服务建议本地/CI 先编译为二进制，再上传服务器运行

示例（在本地编译 Linux 二进制）：

```bash
cd backend
go mod tidy
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server
```

服务器运行示例：

```bash
APP_MYSQL_DSN="red_packet:red_packet@tcp(127.0.0.1:3306)/red_packet?charset=utf8mb4&parseTime=True&loc=Local" \
APP_JWT_SECRET="change-me" \
APP_SERVER_PORT="8080" \
./server
```

## 邀请落地链路

- 访问 `/r/:code` 会把邀请码存入 `localStorage(invite_code)`
- 登录成功后自动尝试调用 `/referral/bind`
- 已绑定或异常时保留主流程，不阻断登录
