# Blog

[原文地址](https://dev.to/umschaudhary/blog-project-with-go-gin-mysql-and-docker-part-1-3cg1)

## 依赖

- gin: web 框架

  ```sh
  go get github.com/gin-gonic/gin
  ```

- gorm: ORM

  ```sh
  go get gorm.io/driver/mysql gorm.io/gorm
  ```

- godotenv: env 文件变量

  ```sh
  go get github.com/joho/godotenv
  ```

## 设计

- api：业务相关

  - controller：web 逻辑层
  - repository：数据库逻辑
  - routes：api 路由集（）
  - service：连接层（controller & repository）

- infrastructure

- models

- util

- main.go
