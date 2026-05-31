# sub2api-distributor 验收与联调文档

## 1. 目标

这份文档用于统一说明：

- 如何初始化本地演示数据
- 如何启动前后端
- 如何运行自动化测试
- 如何执行真实 API 验收
- 当前验收覆盖到了哪些功能

## 2. 项目位置

- 新项目：`/Users/lhl/Desktop/code/sub2api-distributor`
- 主系统：`/Users/lhl/Desktop/code/sub2api`
- 主数据库：`sub2api`

## 3. 当前联调地址

- 前端推荐地址：`http://127.0.0.1:5177/`
- 后端 API：`http://127.0.0.1:8091/api`
- 健康检查：`http://127.0.0.1:8091/health`

说明：

- 后端 CORS 同时兼容 `5173 / 5176 / 5177`
- 当前文档、手工联调和前端默认使用都以 `5177` 为推荐端口

## 4. 初始化演示数据

```bash
cd /Users/lhl/Desktop/code/sub2api-distributor
./scripts/seed_demo_data.sh
```

脚本会：

1. 执行 `backend/migrations/001_distributor_tables.sql`
2. 执行 `backend/migrations/002_seed_demo_operator.sql`
3. 写入一个演示分销商和一个演示被邀请用户
4. 清理 demo 邀请链上的历史 `accrue / transfer` 流水和历史 review 测试分销商资料
5. 重置并生成演示提现单与提现事件

## 5. 演示账号

### 分销商账号

- 邮箱：`dist_demo@example.com`
- 密码：`Viewer123!`
- 角色：`distributor`
- 特征：已存在返利流水、邀请码和提现演示数据

### 被邀请人账号

- 邮箱：`invitee_demo@example.com`
- 密码：`Viewer123!`
- 角色：主系统普通用户
- 说明：默认不登录分销系统，仅用于演示邀请关系和返利来源

### 运营账号

- 邮箱：`operator_demo@local.dev`
- 密码：`Distributor123!`
- 角色：`operator`

## 6. 启动步骤

### 6.1 启动后端

```bash
cd /Users/lhl/Desktop/code/sub2api-distributor
SERVER_PORT=8091 \
DATABASE_DSN='postgres://sub2api@localhost:5432/sub2api?sslmode=disable' \
JWT_SECRET='sub2api-distributor-dev-secret' \
go run -mod=mod ./backend/cmd/server
```

### 6.2 启动前端

```bash
cd /Users/lhl/Desktop/code/sub2api-distributor/frontend
pnpm install
pnpm dev --host 127.0.0.1 --port 5177
```

如果要让分销商页面里的邀请链接直接跳转到主系统注册页，请额外配置：

```bash
VITE_MAIN_APP_BASE_URL=http://127.0.0.1:5173
```

### 6.3 Docker 启动

```bash
cd /Users/lhl/Desktop/code/sub2api-distributor/deploy
cp .env.example .env
docker compose up -d --build
```

说明：

- Docker 方式会自动构建前端并由后端统一托管
- 启动后可以直接访问根路径，不需要单独再跑 `pnpm dev`

### 6.4 systemd 托管启动

```bash
sudo bash /Users/lhl/Desktop/code/sub2api-distributor/deploy/install.sh
sudo systemctl enable --now sub2api-distributor
```

## 7. 自动化验证命令

### 7.1 后端测试

```bash
cd /Users/lhl/Desktop/code/sub2api-distributor
go test ./backend/internal/app ./backend/internal/auth ./backend/internal/config ./backend/internal/distributor ./backend/internal/httpapi/...
```

当前覆盖重点：

- 本地 CORS 白名单
- JWT 签发与解析
- 配置默认值和环境覆盖
- 提现状态机和金额逻辑

### 7.2 前端测试

```bash
cd /Users/lhl/Desktop/code/sub2api-distributor/frontend
pnpm test
```

当前覆盖重点：

- 金额/时间/状态格式化
- session 持久化读写
- 菜单配置
- 基础 HTTP 请求包装

### 7.3 前端构建

```bash
cd /Users/lhl/Desktop/code/sub2api-distributor/frontend
pnpm build
```

构建产物默认输出到：

- `backend/internal/web/dist`

### 7.4 API 验收

```bash
cd /Users/lhl/Desktop/code/sub2api-distributor
./scripts/api_acceptance.sh
```

验收结果会输出到：

- `test-results/api-acceptance-YYYYMMDD-HHMMSS.md`
- `test-results/api-acceptance-latest.md`

说明补充：

- 验收报告按真实请求顺序记录，因此同一份报告中的不同 section 可能代表不同时点的数据状态
- 例如前面的 `dashboard` 是运营操作前的快照，后面的 `ops/withdrawals` 可能已经是运营标记已打款之后的快照
- 如果你要回到统一初始状态，请先重新执行 `./scripts/seed_demo_data.sh`

## 8. API 验收覆盖范围

脚本当前会真实调用以下接口：

### 认证

- `POST /auth/login`
- `GET /me`
- `POST /auth/logout`

### 分销商端

- `GET /portal/dashboard`
- `GET /portal/invite-meta`
- `GET /portal/invitees`
- `GET /portal/rebates`
- `GET /portal/withdrawals`
- `GET /portal/settlement-profile`
- `PUT /portal/settlement-profile`
- `POST /portal/withdrawals`
- `POST /portal/withdrawals/:id/cancel`

### 运营端

- `GET /ops/distributors`
- `GET /ops/users/lookup`
- `GET /ops/distributors/:userId`
- `PUT /ops/distributors/:userId/profile`
- `GET /ops/withdrawals`
- `GET /ops/withdrawals/:id`
- `POST /ops/withdrawals/:id/mark-paid`
- `POST /ops/withdrawals/:id/cancel`

## 9. 页面人工验收建议

推荐按下面顺序做一遍人工检查：

1. 用分销商账号登录，确认能看到左侧 5 个菜单
2. 在概览页确认邀请码卡片、金额卡片、最近提现/返利摘要显示正常
3. 进入邀请用户、返利明细、提现申请，确认表格可滚动
4. 发起提现并取消一笔提现
5. 退出后用运营账号登录，确认只看到 2 个运营菜单
6. 在分销商管理里搜索并查看用户
7. 在提现管理中查看详情、标记已打款、取消申请

## 10. 常见问题

### 10.1 登录提示 `distributor not enabled`

说明当前用户没有启用中的 `distributor_profiles` 记录。可以：

- 通过 `seed_demo_data.sh` 重新生成演示分销商
- 或者使用运营端 `分销商管理` 页面为某个用户开通

说明补充：

- 现在运营端重新启用分销商时，会自动确保该用户在主库里有 `user_affiliates` 记录
- 所以开通后应当可以立刻看到邀请码，不需要再手工补数据库

### 10.2 浏览器能打开页面但登录接口报跨域

请确认：

- 前端是否运行在 `127.0.0.1` 或 `localhost`
- 端口是否在 `5173 / 5176 / 5177` 之内
- 后端是否启动在 `8091`

### 10.3 演示账号和页面数据不一致

优先重新执行：

```bash
cd /Users/lhl/Desktop/code/sub2api-distributor
./scripts/seed_demo_data.sh
```
