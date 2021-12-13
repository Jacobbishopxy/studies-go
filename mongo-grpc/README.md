# Mongo Grpc

[原文地址](https://itnext.io/learning-go-mongodb-crud-with-grpc-98e425aeaae6)

## Note

- 数据类型转换： Protobuf Message (Request) → Regular Go Struct → Convert to BSON + Mongo Action → Protobuf Message (Response)

- 客户端使用 cobra 生成模板代码： `cobra init --pkg-name=client`

  - 使用 cobra 创建新命令文件，例如： `cobra add create`

- 客户端 gRPC 命令

  - 创建 post： `go run . create -a "Nico Vergauwen" -t "Learning Go" -c "The quick brown fox jumps over the lazy dog"`

  - 获取所有 posts：`go run . list`

  - 获取 post（根据 MongoDB 生成的 id 进行查询）：`go run . read -i "61b737ed901cf4202a71fb47"`

  - 更新 post（根据 MongoDB 生成的 id 进行更新）：`go run . update -i "61b737ed901cf4202a71fb47" -a "Nico Vergauwen" -t "Learning Go" -c "Updated content"`

  - 删除 post（根据 MongoDB 生成的 id 进行删除）：`go run . delete -i "61b737ed901cf4202a71fb47"`

- Makefile

  - `proto-gen`: 由 `proto` 文件下的 `*.proto` 生成 `*.go` 文件至 `pb` 目录下

  - `proto-clean`：清除 `pb` 目录下的所有 go 文件

  - `mongo-init`：初始化 MongoDB

  - `server-start`：启动 gRPC 服务端

  - `client-create`：启动 gRPC 客户端，创建 post

  - `client-list`：启动 gRPC 客户端，获取所有 posts
