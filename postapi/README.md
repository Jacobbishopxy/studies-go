# Write a Simple REST API in Golang

[原文地址](https://dev.to/lucasnevespereira/write-a-rest-api-in-golang-following-best-practices-pe9)

## 依赖

- mux：路由
- sqlx：数据库连接
- pq：postgres driver

## 设计

- [database](./app/database)：数据库相关，包含配置、业务接口、业务逻辑、以及数据库 schema（表初始化）
- [models](./app/models)：业务数据结构
- [handlers](./app/handlers.go)：API 接口方法
- [helpers](./app/helpers.go)：util
