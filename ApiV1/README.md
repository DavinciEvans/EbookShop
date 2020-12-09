# EbookShop 使用说明文档

启动项目：

``` shell script
docker-compose up -d
```

修改配置文件 `config.json`：

- **请先修改 `SecretKey` 和数据库密码**

- 若是在本机上部署，请先在 postgresql 新建 `ebookshop` 数据库，并将 `host` 修改为 `localhost`。其余配置项根据自己的数据库进行修改。

- docker-compose 部署，直接运行启动命令即可