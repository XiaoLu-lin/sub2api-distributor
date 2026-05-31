# sub2api-distributor frontend

这是 `sub2api-distributor` 的前端子项目，基于 `Vue 3 + Vite + TypeScript`，提供：

- 分销商后台
- 运营后台

## 常用命令

安装依赖：

```bash
pnpm install
```

本地开发：

```bash
pnpm dev --host 127.0.0.1 --port 5177
```

运行测试：

```bash
pnpm test
```

构建产物：

```bash
pnpm build
```

预览构建结果：

```bash
pnpm preview
```

## 目录说明

- `src/router`：路由和权限守卫
- `src/session`：登录态持久化
- `src/api`：请求封装和接口调用
- `src/views/portal`：分销商端页面
- `src/views/ops`：运营端页面
- `src/components/common`：轻量后台通用组件
- `src/components/layout`：登录布局和后台布局

## 当前测试范围

当前前端测试聚焦在轻量逻辑层：

- 格式化工具
- session 存储
- 菜单配置
- 基础 HTTP 请求封装

暂未覆盖页面渲染和浏览器交互级测试。
