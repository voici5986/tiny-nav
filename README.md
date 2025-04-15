# TinyNav · 非常简单的个人导航网站

[![GitHub Stars](https://img.shields.io/github/stars/hanxi/tiny-nav?style=flat-square)](https://github.com/hanxi/tiny-nav/stargazers)
[![Docker Pulls](https://img.shields.io/docker/pulls/hanxi/tiny-nav?style=flat-square)](https://hub.docker.com/r/hanxi/tiny-nav)
[![Docker Image Size](https://img.shields.io/docker/image-size/hanxi/tiny-nav?style=flat-square)](https://hub.docker.com/r/hanxi/tiny-nav)

> ✨ 一款极简、自托管的个人导航网站，基于 Go + Vue 开发。

**在线体验地址** 👉 [https://nav.hanxi.cc](https://nav.hanxi.cc)  
无需账号密码即可访问，请勿修改或删除公共数据 🙏

---

## 🐳 使用 Docker 快速部署

### 使用 Docker Compose

#### 国际镜像：

```yaml
services:
  tiny-nav:
    image: hanxi/tiny-nav
    container_name: tiny-nav
    restart: unless-stopped
    ports:
      - 8080:58080
    environment:
      NAV_USERNAME: admin
      NAV_PASSWORD: 123456
    volumes:
      - /tiny-nav-data:/app/data
```

#### 国内镜像：

```yaml
services:
  tiny-nav:
    image: docker.hanxi.cc/hanxi/tiny-nav
    container_name: tiny-nav
    restart: unless-stopped
    ports:
      - 8080:58080
    environment:
      NAV_USERNAME: admin
      NAV_PASSWORD: 123456
    volumes:
      - /tiny-nav-data:/app/data
```

启动命令：

```bash
docker compose up -d
```

### 使用 Docker 运行

#### 国际镜像：

```bash
docker run -d \
  --name tiny-nav \
  -p 8080:58080 \
  -e NAV_USERNAME=admin \
  -e NAV_PASSWORD=123456 \
  -v /tiny-nav-data:/app/data \
  hanxi/tiny-nav
```

### 国内镜像：

```bash
docker run -d \
  --name tiny-nav \
  -p 8080:58080 \
  -e NAV_USERNAME=admin \
  -e NAV_PASSWORD=123456 \
  -v /tiny-nav-data:/app/data \
  docker.hanxi.cc/hanxi/tiny-nav
```

访问页面：打开浏览器访问 http://<你的IP>:8080

## 🧩 本地运行（非 Docker）

1. 前往 Releases 页面 下载对应平台的可执行文件
2. 无认证启动：
```bash
./tiny-nav --port=58080 --no-auth
```
3. 有账号密码启动：
```bash
./tiny-nav --port=58080 --user=admin --password=123456
````
4. 访问地址：http://localhost:58080

## 🔧 从源码编译

```bash
sh build.sh
```

将生成 tiny-nav 可执行文件，所有前端资源已打包至其中。运行示例：

```bash
ENABLE_NO_AUTH=true LISTEN_PORT=58080 ./tiny-nav
```

访问：http://localhost:58080

## 🧱 技术栈

- 后端：Golang
- 前端：Vue 3

## 📌 开发计划

- [ ] 支持只读模式：查看免登录，编辑需登录
- [ ] 数据 MD5 对比，避免重复加载
- [ ] 自动深色模式
- [ ] 支持书签导入
- [ ] 支持站内搜索

