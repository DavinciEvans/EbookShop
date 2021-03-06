# EBookShop 网上电子书书店

使用 Go + PostgreSQL 搭建的网上电子书书店的数据库后端，是数据库的课程设计。

部署使用 Docker Compose。

## 使用

在使用之前请先确保安装了 Docker 和 Docker Compose。

``` shell
git clone https://github.com/DavinciEvans/EbookShop.git
cd EbookShop
docker-compose up -d
```

### Logs

- run.log：Go 的响应日志
- db.log：数据库日志

### API List

API 的具体使用说明已放在 README 文件夹当中。同时还附有配套的 postman  示例文件，直接导入 postman 当中即可。

### 获取新版本

仅针对不熟悉 git 的同学使用：

```shell
cd /PATH/TO/EbookShop
git fetch origin
git pull origin
```

## 设置

配置文件为 `config.json`

- **请先修改 SecretKey**
- 使用 Docker 部署，若需要修改数据库用户和密码，请在 `config.json`中修改完后进入 `docker-compose.yml` 一并修改。
- 在本机上部署，请自行安装 Go 与 PostgreSQL，并在 PostgreSQL 当中创建名为 ebookshop 的数据库。之后还需修改配置文件中的 `Username` 与 `Password` 项，并将 `host` 项修改为 `localhost`。
- `Forge` 项为是否需要载入随机数据。在第一次启动之后会修改为 False，如果需要重新载入，请手动修改为 True。

## 更新日志

- 6.8：
  1. 增加了 User 的 Get 接口，使用 /api/v1/auth 可以获取当前登录用户的信息
- 12.17:
  1. 付款时，会自动删除购物车里的东西
  2. 修复了删除书籍时不存在也返回 200 的错误
  3. 在 nginx 上增加前端挂载点
- 12.16:
  1. 修复了登陆路由不正确导致无法登录的问题
  2. 修改了获取单个书籍信息时的查询语句
- 12.14:
  1. 增加了 Nginx，以后使用可以直接从 localhost 访问了。
- 12.13:
  1. 修复无法启动的 bug
- 12.10：
  1. 修复了关于用户登录的 bug。现在用户登录统一使用 /auth POST 接口进行登录，注册用户统一使用 /auth/newAuth POST 接口进行注册。
  2. 修改书籍返回信息中的 category 为名称而不是 ID。
