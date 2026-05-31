# 项目级规则

在遵循全局规则的基础上，`sub2api-distributor` 额外遵循以下约定。

## 项目形态

- 项目类型：全栈独立项目
- 后端技术栈：Go + Gin + PostgreSQL
- 前端技术栈：Vue 3 + Vite + TypeScript
- 包管理器：前端使用 `pnpm`

## 代码组织

- 后端代码位于 `backend/`
- 前端代码位于 `frontend/`
- 说明文档位于 `docs/`
- OpenSpec 变更工件位于 `openspec/`

## 项目约束

- 分销商结算逻辑保持独立，不回写主系统旧 `/affiliate` 页面逻辑
- 登录继续复用主库 `users`，但分销商后台权限由 `distributor_profiles` 决定
- 前端新增后台功能时，优先复用现有轻量组件，不引入额外 UI 大依赖

## 工作流入口

- 本仓库正式 change 默认入口：`openspec-superpowers-orchestrator`
- `openspec-explore`、`openspec-propose`、`openspec-apply-change`、`openspec-archive-change` 属于阶段技能
- `/opsx:*` 属于 OpenSpec 阶段命令入口，不作为默认顶层入口

## 验证

- 后端测试：`go test ./backend/internal/app ./backend/internal/distributor ./backend/internal/httpapi/...`
- 前端构建：`cd frontend && pnpm build`
