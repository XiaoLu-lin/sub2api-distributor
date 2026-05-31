# sub2api-distributor 部署说明

这套部署方案尽量对齐主项目 `sub2api` 的交付方式，采用：

- `Dockerfile`
- `docker compose`
- `.env` 环境变量
- 可选 `Caddy` 反向代理
- 可选 `systemd` 托管
- 可选安装脚本
- 自动部署脚本

## 1. 部署前提

- 已有可访问的 `sub2api` PostgreSQL 数据库
- 目标服务器已安装 Docker 和 Docker Compose
- 你已经准备好一个固定域名，例如 `distributor.example.com`

## 2. 目录说明

- 根目录镜像构建文件：[Dockerfile](/Users/lhl/Desktop/code/sub2api-distributor/Dockerfile)
- Compose 配置：
  - [docker-compose.yml](/Users/lhl/Desktop/code/sub2api-distributor/deploy/docker-compose.yml)
  - [docker-compose.local.yml](/Users/lhl/Desktop/code/sub2api-distributor/deploy/docker-compose.local.yml)
- 环境变量示例：[.env.example](/Users/lhl/Desktop/code/sub2api-distributor/deploy/.env.example)
- Caddy 示例：[Caddyfile](/Users/lhl/Desktop/code/sub2api-distributor/deploy/Caddyfile)
- 自动部署脚本：
  - [deploy.sh](/Users/lhl/Desktop/code/sub2api-distributor/deploy/deploy.sh)
  - [redeploy.sh](/Users/lhl/Desktop/code/sub2api-distributor/deploy/redeploy.sh)
  - [check.sh](/Users/lhl/Desktop/code/sub2api-distributor/deploy/check.sh)

## 3. 启动步骤

```bash
cd /Users/lhl/Desktop/code/sub2api-distributor/deploy
cp .env.example .env
```

编辑 `.env`，至少确认这些值：

- `DATABASE_DSN`
- `JWT_SECRET`
- `VITE_MAIN_APP_BASE_URL`
- `CORS_ALLOWED_ORIGINS`
- `VITE_API_BASE_URL`

然后启动：

```bash
docker compose up -d --build
```

查看日志：

```bash
docker compose logs -f sub2api-distributor
```

如果你要直接走自动部署脚本：

```bash
bash /opt/sub2api-distributor/deploy/deploy.sh
```

如果你要拉最新代码并重部署：

```bash
bash /opt/sub2api-distributor/deploy/redeploy.sh
```

## 3.1 使用 systemd 托管

服务文件：

- [sub2api-distributor.service](/Users/lhl/Desktop/code/sub2api-distributor/deploy/sub2api-distributor.service)

安装脚本：

- [install.sh](/Users/lhl/Desktop/code/sub2api-distributor/deploy/install.sh)

如果代码已经位于服务器目标目录，例如：

```bash
/opt/sub2api-distributor
```

那么可以直接执行：

```bash
sudo bash /opt/sub2api-distributor/deploy/install.sh
sudo systemctl enable --now sub2api-distributor
```

说明：

- 脚本会检测当前项目目录是否已经等于安装目录
- 如果已经在目标目录中，只会补 `.env`、安装 `systemd` 服务和重载 `systemd`
- 如果不在目标目录中，才会执行文件复制

## 4. 健康检查

应用启动后可直接检查：

```bash
curl http://127.0.0.1:8091/health
```

预期返回：

```json
{"status":"ok"}
```

## 5. 和主系统保持一致的地方

- 都采用 `deploy/` 目录集中管理部署文件
- 都采用 `docker compose + .env`
- 都提供 `Caddyfile` 反向代理示例
- 都通过固定 `JWT_SECRET` 保证会话稳定

## 6. 当前部署边界

当前版本部署时：

- 分销商后台前端已内置到镜像中，由后端统一托管
- 后端提供 `/api/*` 和前端页面路由
- 数据库直接连接主系统 `sub2api` 主库
- 不额外创建 Redis、PostgreSQL 容器
- 默认启动时自动执行 `backend/migrations/*.sql`

## 7. 上线前建议

- 使用独立生产域名
- `DATABASE_DSN` 使用真实数据库账号密码
- `JWT_SECRET` 使用随机强密钥
- `CORS_ALLOWED_ORIGINS` 只保留真实前台域名
- 如与主系统同机部署，优先通过内网地址或容器网络连接主库

## 8. 发布检查清单

上线前请逐项确认：

- `.env` 中的 `DATABASE_DSN` 已改成生产数据库
- `.env` 中的 `JWT_SECRET` 已替换成随机强密钥
- `.env` 中的 `VITE_MAIN_APP_BASE_URL` 指向主系统真实注册地址
- `.env` 中的 `CORS_ALLOWED_ORIGINS` 只包含真实分销商后台域名
- 目标数据库账号对 `distributor_*` 表有建表/更新权限
- `docker compose up -d --build` 已执行成功
- `curl http://127.0.0.1:8091/health` 返回 `{"status":"ok"}`
- 浏览器能正常访问登录页
- 分销商登录、运营登录至少各手工验证一次
- Caddy 或其他反代已配置 HTTPS
