# sub2api-distributor 服务器执行命令版清单

这份文档只保留服务器实际执行步骤，适合部署时直接照着执行。

---

## 1. 准备部署目录

```bash
sudo mkdir -p /opt/sub2api-distributor
sudo chown -R $USER:$USER /opt/sub2api-distributor
cd /opt/sub2api-distributor
```

---

## 2. 拉取代码

首次部署：

```bash
git clone <你的仓库地址> .
```

如果已经部署过，更新代码：

```bash
cd /opt/sub2api-distributor
git pull
```

---

## 3. 进入部署目录并生成配置

```bash
cd /opt/sub2api-distributor/deploy
cp .env.example .env
```

---

## 4. 编辑生产环境变量

```bash
vim /opt/sub2api-distributor/deploy/.env
```

至少修改这些值：

```bash
DATABASE_DSN=postgres://user:password@host:5432/sub2api?sslmode=disable
JWT_SECRET=替换成真实强密钥
VITE_MAIN_APP_BASE_URL=https://你的主系统域名
CORS_ALLOWED_ORIGINS=https://你的分销商域名
VITE_API_BASE_URL=/api
RUN_MIGRATIONS_ON_STARTUP=true
```

---

## 5. 构建并启动容器

```bash
cd /opt/sub2api-distributor/deploy
docker compose up -d --build
```

或者直接执行自动部署脚本：

```bash
bash /opt/sub2api-distributor/deploy/deploy.sh
```

---

## 6. 查看容器状态

```bash
docker compose ps
```

---

## 7. 查看运行日志

```bash
docker compose logs -f sub2api-distributor
```

---

## 8. 验证健康检查

```bash
curl http://127.0.0.1:8091/health
```

预期返回：

```json
{"status":"ok"}
```

---

## 9. 本机验证页面

```bash
curl -I http://127.0.0.1:8091/
```

---

## 10. 安装 systemd 托管

```bash
sudo bash /opt/sub2api-distributor/deploy/install.sh
sudo systemctl daemon-reload
sudo systemctl enable --now sub2api-distributor
```

---

## 11. 查看 systemd 状态

```bash
sudo systemctl status sub2api-distributor
```

---

## 12. 常用维护命令

拉最新代码并重部署：

```bash
bash /opt/sub2api-distributor/deploy/redeploy.sh
```

重启服务：

```bash
sudo systemctl restart sub2api-distributor
```

停止服务：

```bash
sudo systemctl stop sub2api-distributor
```

查看服务状态：

```bash
sudo systemctl status sub2api-distributor
```

查看容器日志：

```bash
docker compose -f /opt/sub2api-distributor/deploy/docker-compose.yml logs -f sub2api-distributor
```

再次检查健康接口：

```bash
curl http://127.0.0.1:8091/health
```

---

## 13. 使用 Caddy 时的配置命令

编辑配置文件：

```bash
vim /opt/sub2api-distributor/deploy/Caddyfile
```

把里面的：

```text
distributor.example.com
```

改成你的真实域名。

重载 Caddy：

```bash
sudo systemctl reload caddy
```

---

## 14. 最终浏览器验收

打开：

```text
https://你的分销商域名/
```

需要人工确认：

- 登录页正常打开
- 运营账号可登录
- 分销商账号可登录
- 邀请码和邀请链接正确

---

## 15. 部署前还缺的真实信息

正式执行前，你还需要准备：

- 仓库地址
- 分销商正式域名
- 主系统正式域名
- 生产 `DATABASE_DSN`
- 生产 `JWT_SECRET`
