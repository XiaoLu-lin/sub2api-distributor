# sub2api-distributor

独立的分销商结算门户系统，面向 `sub2api` 主库提供两套后台：

- 分销商端：查看返利、邀请用户、邀请码、提现申请、收款信息
- 运营端：开通分销商、查看提现申请、线下打款后手动标记已打款

## 项目定位

这个项目不继续在主系统旧 `/affiliate` 页面上扩功能，而是单独提供一套更清晰的分销商结算流程：

1. 分销商邀请用户
2. 被邀请用户在主系统内产生返利
3. 分销商在新系统里查看累计返利和可提现金额
4. 分销商提交提现申请后，状态立即进入 `打款中`
5. 运营线下打款后，手动标记 `已打款`

## 技术栈

- 后端：Go + Gin + PostgreSQL
- 前端：Vue 3 + Vite + TypeScript
- 主数据来源：复用 `sub2api` 主库中的 `users / user_affiliates / user_affiliate_ledger`
- 部署方式：`Dockerfile + docker compose + .env`

前端包管理说明：

- 前端使用 `pnpm`
- 为避免本地与 Docker / 服务器构建环境漂移，`frontend/package.json` 已固定：
  - `packageManager: pnpm@9.12.0`
- 如果服务器上直接运行 `pnpm install` 或在 Docker 中通过 `corepack` 安装依赖，请保持与该版本一致

## 协作工作流

- 本仓库已接入 OpenSpec for Codex
- 正式 change 默认从 `openspec-superpowers-orchestrator` 进入
- 最近一次已归档质量基线变更位于：
  - [2026-05-31-standardize-project-quality](/Users/lhl/Desktop/code/sub2api-distributor/openspec/changes/archive/2026-05-31-standardize-project-quality/proposal.md)
- 最近一次已归档的完整性修复变更位于：
  - [2026-05-31-fix-distributor-invite-and-settlement-integrity](/Users/lhl/Desktop/code/sub2api-distributor/openspec/changes/archive/2026-05-31-fix-distributor-invite-and-settlement-integrity/proposal.md)
- 当前正在进行中的 review findings 修复变更位于：
  - [fix-distributor-review-findings-round-2](/Users/lhl/Desktop/code/sub2api-distributor/openspec/changes/fix-distributor-review-findings-round-2/proposal.md)
  - 说明：本轮代码与验收已按该变更执行并通过本地验证，后续如确认收口，可继续归档该 change

## 目录结构

```text
backend/
  cmd/server
  internal/
  migrations/
frontend/
docs/
scripts/
openspec/
test-results/
```

## 环境变量

后端：

```bash
SERVER_PORT=8091
DATABASE_DSN=postgres://sub2api@localhost:5432/sub2api?sslmode=disable
JWT_SECRET=sub2api-distributor-dev-secret
```

前端：

```bash
VITE_API_BASE_URL=http://127.0.0.1:8091/api
VITE_MAIN_APP_BASE_URL=http://127.0.0.1:5173
```

生产部署补充：

```bash
APP_ENV=production
CORS_ALLOWED_ORIGINS=https://distributor.example.com
STATIC_DIR=/app/web
```

## 初始化演示数据

项目默认以本地 `sub2api` 数据库作为演示环境，初始化命令：

```bash
cd /Users/lhl/Desktop/code/sub2api-distributor
./scripts/seed_demo_data.sh
```

这个脚本会：

- 创建或更新 `distributor_` 相关表
- 写入运营演示账号
- 保留一个分销商和一个被邀请用户的最小演示链路
- 清理 demo 邀请链上的历史 `accrue / transfer` 流水
- 清理历史 review 遗留的测试分销商资料
- 重置并生成演示提现单和事件流

## 启动后端

```bash
cd /Users/lhl/Desktop/code/sub2api-distributor
SERVER_PORT=8091 \
DATABASE_DSN='postgres://sub2api@localhost:5432/sub2api?sslmode=disable' \
JWT_SECRET='sub2api-distributor-dev-secret' \
go run -mod=mod ./backend/cmd/server
```

默认地址：

- API：`http://127.0.0.1:8091/api`
- 健康检查：`http://127.0.0.1:8091/health`

## 启动前端

```bash
cd /Users/lhl/Desktop/code/sub2api-distributor/frontend
pnpm install
pnpm dev --host 127.0.0.1 --port 5177
```

推荐本地地址：

- 前端：`http://127.0.0.1:5177/`

说明：

- 后端 CORS 兼容 `5173 / 5176 / 5177`
- 当前文档和验收基线统一以 `5177` 作为推荐联调端口

## 演示账号

演示数据以 `./scripts/seed_demo_data.sh` 结果为准：

- 分销商：`dist_demo@example.com / Viewer123!`
- 被邀请人：`invitee_demo@example.com / Viewer123!`
- 运营：`operator_demo@local.dev / Distributor123!`

说明：

- 运营账号仅用于本地演示和验收
- 分销商是否能登录，取决于 `distributor_profiles` 中是否存在 `status = active` 的记录
- 被邀请人默认不登录本系统，它用于演示“被邀请后产生返利”的上游身份

邀请码链接说明：

- 分销商端展示的邀请码来自主库 `user_affiliates.aff_code`
- 邀请链接会拼接到 `VITE_MAIN_APP_BASE_URL/register?aff=<code>`
- 如果未配置 `VITE_MAIN_APP_BASE_URL`，前端会明确提示“未配置主系统注册地址”，不会再生成假链接

## 当前实现范围

已实现：

- 主库账号密码登录
- 分销商身份校验
- 分销商端后台路由、菜单、表格、弹窗
- 运营端后台路由、菜单、表格、弹窗
- 分销商邀请码查询 `GET /api/portal/invite-meta`
- 分销商邀请用户、返利明细、提现申请、收款信息
- 运营分销商管理、用户搜索开通、提现管理
- 线下打款后手动标记 `已打款`
- 本地 API 验收脚本
- 前端轻量单测基线
- 后端核心单元测试基线

本轮未实现：

- 自动打款
- 审核流
- 服务端分页/复杂筛选
- 完整数据库集成测试
- CI 平台集成

## 验证命令

后端测试：

```bash
cd /Users/lhl/Desktop/code/sub2api-distributor
go test ./backend/internal/app ./backend/internal/auth ./backend/internal/config ./backend/internal/distributor ./backend/internal/httpapi/...
```

前端测试：

```bash
cd /Users/lhl/Desktop/code/sub2api-distributor/frontend
pnpm test
```

前端构建：

```bash
cd /Users/lhl/Desktop/code/sub2api-distributor/frontend
pnpm build
```

API 验收：

```bash
cd /Users/lhl/Desktop/code/sub2api-distributor
./scripts/api_acceptance.sh
```

Docker 部署：

```bash
cd /Users/lhl/Desktop/code/sub2api-distributor/deploy
cp .env.example .env
docker compose up -d --build
```

如果构建时报类似下面的错误：

- `Lockfile failed supply-chain policy check`
- `minimumReleaseAge`
- `tinyglobby ... within the minimumReleaseAge cutoff`

优先检查：

1. 是否已经拉到最新代码
2. `frontend/package.json` 中是否包含：
   - `"packageManager": "pnpm@9.12.0"`
3. Docker 构建是否仍在使用旧缓存

常用处理方式：

```bash
cd /Users/lhl/Desktop/code/sub2api-distributor
git pull
cd deploy
docker compose build --no-cache
docker compose up -d
```

systemd 托管：

```bash
sudo bash /Users/lhl/Desktop/code/sub2api-distributor/deploy/install.sh
sudo systemctl enable --now sub2api-distributor
```

详细联调与验收说明见：

- [TESTING.md](/Users/lhl/Desktop/code/sub2api-distributor/TESTING.md)
- [distributor-system-overview.md](/Users/lhl/Desktop/code/sub2api-distributor/docs/distributor-system-overview.md)
- [distributor-workflows.md](/Users/lhl/Desktop/code/sub2api-distributor/docs/distributor-workflows.md)
- [deploy/README.md](/Users/lhl/Desktop/code/sub2api-distributor/deploy/README.md)
- [deploy/PRODUCTION_DEPLOYMENT.md](/Users/lhl/Desktop/code/sub2api-distributor/deploy/PRODUCTION_DEPLOYMENT.md)
- [deploy/SERVER_COMMAND_CHECKLIST.md](/Users/lhl/Desktop/code/sub2api-distributor/deploy/SERVER_COMMAND_CHECKLIST.md)

如果你要看“按页面怎么操作、每一步会发生什么、金额如何流转”，优先看：

- [distributor-workflows.md](/Users/lhl/Desktop/code/sub2api-distributor/docs/distributor-workflows.md)
