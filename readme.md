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

## 设置

配置文件为 `config.json`

- **请先修改 SecretKey**
- 使用 Docker 部署，若需要修改数据库用户和密码，请在 `config.json`中修改完后进入 `docker-compose.yml` 一并修改。
- 在本机上部署，请自行安装 Go 与 PostgreSQL，并在 PostgreSQL 当中创建名为 ebookshop 的数据库。之后还需修改配置文件中的 `Username` 与 `Password` 项，并将 `host` 项修改为 `localhost`。
- `Forge` 项为是否需要载入随机数据。在第一次启动之后会修改为 False，如果需要重新载入，请手动修改为 True。