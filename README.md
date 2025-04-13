# 非常简单的个人导航网站

使用 [豆包](https://www.doubao.com/) 和 [copilot](https://github.com/copilot) 辅助开发。

- 体验地址： <https://nav.hanxi.cc/>
- 不需要账号密码，请不要删东西。

## 使用 Docker 运行

### 用 Docker compose 启动

新建 `docker-compose.yml` 文件，内容如下：

```yml
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

国内镜像：

```yml
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

### 用 Docker 启动

启动命令：

```bash
docker run -d \
  --name tiny-nav \
  -p 8080:58080 \
  -e NAV_USERNAME=admin \
  -e NAV_PASSWORD=123456 \
  -v /tiny-nav-data:/app/data \
  hanxi/tiny-nav
```

国内镜像启动：

```bash
docker run -d \
  --name tiny-nav \
  -p 8080:58080 \
  -e NAV_USERNAME=admin \
  -e NAV_PASSWORD=123456 \
  -v /tiny-nav-data:/app/data \
  docker.hanxi.cc/hanxi/tiny-nav
```

### 进入网站页面

使用浏览器访问 <http://ip:8080> 即可, ip 改成你机器的ip。

## 下载运行

1. 去 <https://github.com/hanxi/tiny-nav/releases> 下载对应平台的可执行文件
2. 以无用户密码的方式运行

```bash
./tiny-nav --port=58080 --no-auth
```

3. 打开浏览器访问 <http://localhost:58080> 即可。
4. 以有用户密码的方式运行

```bash
./tiny-nav --port=58080 --user=admin --password=123456
```

## 编译运行

### 编译

```
sh build.sh
```

这样会生成 `tiny-nav` 可执行文件。所有静态资源会被打包到 `tiny-nav` 可执行文件中。

### 启动

```
ENABLE_NO_AUTH=true LISTEN_PORT=58080 ./tiny-nav
```

网页访问 <http://localhost:58080> 即可。

## 技术栈

- 后端 Golang
- 前端 Vue

## 未来开发计划

- [ ] 查看模式：编辑需要账号密码，查看可以不用账号密码。
- [ ] 对比数据 md5 值，没变化则使用本地数据。
- [ ] 自动深色模式。
