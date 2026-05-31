# sub2api-distributor 服务器上线文档

这份文档用于指导 `sub2api-distributor` 上线到正式服务器。

目标：

- 让分销商后台以 Docker 方式运行
- 尽量和主项目 `sub2api` 的部署方式保持一致
- 明确上线前还缺哪些生产信息

---

## 1. 部署前提

服务器需要满足：

- Linux 服务器，建议 Ubuntu 22.04 及以上
- 已安装 `git`
- 已安装 `curl`
- 已安装 `bash`
- 已安装 Docker
- 已安装 Docker Compose
- 已开放 `80` 和 `443` 端口

说明：

- 如果你要先不走域名、直接验证服务，也可以临时开放 `8091`
- 正式环境建议通过 Caddy 或其他反向代理走 `443`

---

## 2. 部署架构

当前项目部署时：

- 前端会在构建时打包进镜像
- 后端统一提供：
  - `/api/*` 接口
  - 前端页面路由
  - `/health` 健康检查
- 数据库直接连接主系统 `sub2api` 的 PostgreSQL
- 不额外创建 Redis、PostgreSQL 容器
- 启动时会自动执行 `backend/migrations/*.sql`

---

## 3. 上线前必须准备的信息

在正式部署前，你至少要准备好这些值：

### 3.1 域名信息

- 分销商系统正式域名  
  例如：`distributor.example.com`

- 主系统正式注册地址  
  例如：`https://main.example.com/register`

说明：

- 分销商后台域名用于浏览器访问和 CORS 配置
- 主系统注册地址用于生成分销商邀请链接

### 3.2 数据库信息

- PostgreSQL 主机地址
- PostgreSQL 端口
- 数据库名
- 数据库用户名
- 数据库密码
- 是否需要 SSL

最终会整理成：

```bash
DATABASE_DSN=postgres://user:password@host:5432/sub2api?sslmode=disable
```

### 3.3 鉴权信息

- 一个固定的 `JWT_SECRET`

建议生成方式：

```bash
openssl rand -hex 32
```

### 3.4 反向代理方式

二选一：

- 使用 Caddy
- 使用你现有的 Nginx / 其他网关

当前项目里已提供 Caddy 示例配置。

### 3.5 服务托管方式

二选一：

- 直接使用 `docker compose`
- 使用 `systemd` 托管 `docker compose`

建议正式环境使用 `systemd`

---

## 4. 服务器部署步骤

### 4.1 准备服务器目录

建议部署目录：

```bash
/opt/sub2api-distributor
```

操作：

```bash
sudo mkdir -p /opt/sub2api-distributor
sudo chown -R $USER:$USER /opt/sub2api-distributor
cd /opt/sub2api-distributor
```

### 4.2 拉取项目代码

```bash
cd /opt/sub2api-distributor
git clone <你的仓库地址> .
```

如果已经有代码仓库：

```bash
cd /opt/sub2api-distributor
git pull
```

### 4.3 准备环境变量

进入部署目录：

```bash
cd /opt/sub2api-distributor/deploy
cp .env.example .env
```

然后编辑 `.env`：

```bash
vim /opt/sub2api-distributor/deploy/.env
```

至少要修改这些值：

- `DATABASE_DSN`
- `JWT_SECRET`
- `VITE_MAIN_APP_BASE_URL`
- `CORS_ALLOWED_ORIGINS`
- `VITE_API_BASE_URL`

推荐值示例：

```bash
BIND_HOST=0.0.0.0
SERVER_PORT=8091
APP_ENV=production
TZ=Asia/Shanghai
VITE_API_BASE_URL=/api
RUN_MIGRATIONS_ON_STARTUP=true
DATABASE_DSN=postgres://sub2api:change_this_password@10.0.0.10:5432/sub2api?sslmode=disable
JWT_SECRET=replace_with_a_real_secret
VITE_MAIN_APP_BASE_URL=https://main.example.com
CORS_ALLOWED_ORIGINS=https://distributor.example.com
```

### 4.4 检查数据库权限

部署前确认这个数据库账号具备：

- 可以连接 `sub2api` 主库
- 可以执行 `backend/migrations/*.sql`
- 可以创建或更新 `distributor_*` 表
- 可以读取：
  - `users`
  - `user_affiliates`
  - `user_affiliate_ledger`

如果数据库有限制来源 IP，请先把服务器 IP 加到白名单。

### 4.5 启动服务

先手工启动一版，确认镜像和容器都正常：

```bash
cd /opt/sub2api-distributor/deploy
docker compose up -d --build
```

查看容器状态：

```bash
docker compose ps
```

查看日志：

```bash
docker compose logs -f sub2api-distributor
```

### 4.6 验证健康检查

先在服务器本机执行：

```bash
curl http://127.0.0.1:8091/health
```

预期返回：

```json
{"status":"ok"}
```

如果这里失败，先不要继续配域名，优先查：

- `.env` 是否配置正确
- 数据库是否能连通
- 迁移是否执行失败
- 容器日志是否报错

### 4.7 验证页面

如果你临时暴露了 `8091`，可以先用浏览器直接打开：

```text
http://服务器IP:8091/
```

预期：

- 能看到登录页
- 无白屏
- 无接口 500

---

## 5. 配置域名和 HTTPS

### 5.1 使用 Caddy

项目已提供示例文件：

- [Caddyfile](/Users/lhl/Desktop/code/sub2api-distributor/deploy/Caddyfile)

你需要做的事：

1. 把 `distributor.example.com` 改成你的真实域名
2. 反向代理目标保持为：

```text
localhost:8091
```

3. 重载或重启 Caddy

### 5.2 域名验证

域名配置完成后，浏览器访问：

```text
https://你的分销商域名/
```

需要确认：

- 能打开登录页
- HTTPS 正常
- 浏览器没有 CORS 报错
- 静态资源能正常加载

---

## 6. 配置 systemd 托管

如果你希望服务器重启后自动拉起服务，建议用 `systemd`。

项目已提供：

- [sub2api-distributor.service](/Users/lhl/Desktop/code/sub2api-distributor/deploy/sub2api-distributor.service)
- [install.sh](/Users/lhl/Desktop/code/sub2api-distributor/deploy/install.sh)

执行步骤：

```bash
sudo bash /opt/sub2api-distributor/deploy/install.sh
sudo systemctl enable --now sub2api-distributor
```

查看状态：

```bash
sudo systemctl status sub2api-distributor
```

重启服务：

```bash
sudo systemctl restart sub2api-distributor
```

停止服务：

```bash
sudo systemctl stop sub2api-distributor
```

---

## 7. 上线后验收步骤

建议按下面顺序验收：

1. 验证 `/health`
2. 打开正式域名首页
3. 用运营账号登录
4. 检查 `分销商管理`
5. 用分销商账号登录
6. 检查：
   - 邀请码
   - 邀请链接
   - 邀请用户页
   - 提现申请页
7. 确认邀请链接跳转到主系统正式注册地址

---

## 8. 发布检查清单

上线前逐项确认：

- `.env` 中 `DATABASE_DSN` 已替换成生产库
- `.env` 中 `JWT_SECRET` 已替换成真实强密钥
- `.env` 中 `VITE_MAIN_APP_BASE_URL` 已指向主系统真实注册地址
- `.env` 中 `CORS_ALLOWED_ORIGINS` 只保留真实分销商域名
- 数据库账号具备建表和读写所需权限
- `docker compose up -d --build` 已成功执行
- `curl http://127.0.0.1:8091/health` 返回 `{"status":"ok"}`
- 浏览器可正常打开登录页
- 运营账号和分销商账号至少各登录验证一次
- 域名 HTTPS 已正确配置

---

## 9. 当前还缺什么

如果你现在就准备正式上线，还缺的是这些真实生产信息：

- 分销商正式域名
- 主系统正式注册地址域名
- PostgreSQL 生产连接串
- 生产 `JWT_SECRET`
- 是否使用 Caddy
- 是否使用 `systemd`

---

## 10. 当前已知限制

当前这轮本地已经完成的验证：

- 后端测试通过
- 前端测试通过
- 前端生产构建通过
- 部署脚本语法检查通过

但还有一个限制：

- 当前本机环境没有 `docker` 命令，因此未完成真实 `docker build` 验证

所以正式上线前，建议你在目标服务器上优先执行：

```bash
docker compose up -d --build
docker compose logs -f sub2api-distributor
curl http://127.0.0.1:8091/health
```
