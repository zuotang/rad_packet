# Backend README

## 技术栈

- Go `1.22`
- Echo `v4`
- GORM + MySQL `8.x`
- Viper（配置读取）
- JWT（鉴权）

## 目录说明

- `cmd/server/main.go`：启动入口
- `configs/config.yaml`：本地默认配置
- `internal/config`：配置加载
- `internal/database`：数据库初始化、AutoMigrate、种子数据
- `internal/models`：GORM 模型
- `internal/service`：核心业务（奖励、任务、邀请、钱包、提现等）
- `internal/http/router`：路由注册
- `internal/http/handlers`：接口处理层
- `internal/http/middleware`：JWT 中间件

## 配置项

支持 `config.yaml` 和环境变量（环境变量优先）：

- `APP_SERVER_PORT`，默认 `8080`
- `APP_MYSQL_DSN`，例如：
  `red_packet:red_packet@tcp(127.0.0.1:3306)/red_packet?charset=utf8mb4&parseTime=True&loc=Local`
- `APP_JWT_SECRET`
- `APP_JWT_TTL_HOURS`，默认 `168`

## 本地启动

```bash
cd backend
go mod tidy
go run ./cmd/server
```

健康检查：

```bash
GET http://localhost:8080/healthz
```

## Docker Compose（仅 MySQL）

项目根目录已提供 `docker-compose.yml`，只用于启动 MySQL：

```bash
docker compose up -d
```

## 生产部署（推荐）

先在本地/CI 编译 Linux 二进制，再上传服务器运行。

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

## API 清单

- `POST /api/auth/login`
- `POST /api/auth/otp`
- `GET /api/config/bootstrap`
- `POST /api/referral/bind`（需 JWT）
- `GET /api/referral/status`（需 JWT）
- `GET /api/reward/summary`（需 JWT）
- `GET /api/reward/records?page=1&size=20&status=`（需 JWT）
- `POST /api/reward/unlock`（需 JWT，风控通过后解冻 pending 奖励）
- `GET /api/task/list?country=`（需 JWT，返回用户任务完成状态）
- `POST /api/task/claim`（需 JWT）
- `GET /api/wallet`（需 JWT）
- `POST /api/withdraw/apply`（需 JWT）
- `GET /api/withdraw/records?page=1&size=20&status=`（需 JWT）
- `GET /api/admin/withdraw/list?page=1&size=20&status=`（需 `X-Admin-Key`）
- `POST /api/admin/withdraw/review`（需 `X-Admin-Key`，状态流转：pending->approved/rejected->paid）
- `GET /api/admin/dashboard`（需 `X-Admin-Key`）
- `GET /api/admin/task/list`（需 `X-Admin-Key`）
- `POST /api/admin/task/save`（需 `X-Admin-Key`）
- `DELETE /api/admin/task/:id`（需 `X-Admin-Key`）
- `GET /api/admin/config/list`（需 `X-Admin-Key`）
- `POST /api/admin/config/upsert`（需 `X-Admin-Key`）
- `GET /api/admin/risk/flags` `POST /api/admin/risk/flag/add`（需 `X-Admin-Key`）
- `GET /api/admin/blacklist/list` `POST /api/admin/blacklist/add`（需 `X-Admin-Key`）

## 核心业务约束

- 发奖事务：`reward + wallet_ledger + wallet` 同事务提交
- 发奖幂等：账本唯一键 `(user_id, ref_type, ref_id)`
- 任务领奖幂等：`user_task_events` 唯一键 `(user_id, event_key)`
- 邀请绑定：防自绑，子用户只允许绑定一次
- 有效邀请：被邀请人首个有效任务触发 `referral_edges.is_valid=true` 并给上级发放 `pending` 邀请奖励
- 解冻奖励：`/api/reward/unlock` 将 `pending` 转 `unlocked`，并同步钱包 `frozen -> balance`

## 常见问题

- 启动 panic：`missing LogValuesFunc callback`  
  已兼容 Echo `v4.12`，如果你是旧代码，请更新 `RequestLoggerWithConfig` 并提供 `LogValuesFunc`。
