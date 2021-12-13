# Mongo Grpc

[原文地址](https://itnext.io/learning-go-mongodb-crud-with-grpc-98e425aeaae6)

## Note

- 数据类型转换： Protobuf Message (Request) → Regular Go Struct → Convert to BSON + Mongo Action → Protobuf Message (Response)

- 客户端使用 cobra 生成模板代码： `cobra init --pkg-name=client`

  - 使用 cobra 创建新命令文件，例如： `cobra add create`
