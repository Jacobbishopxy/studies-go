# PcBook

[原文地址](https://dev.to/techschoolguru/series/7311)

## Prerequisites

预安装 `protobuf-compiler`：

```sh
apt-get install protobuf-compiler
```

依赖：

```sh
go get -u google.golang.org/grpc
go get -u github.com/golang/protobuf/protoc-gen-go
```

## Note

由 `protoc` 生成 Golang 代码（已加入 Makefile）：

```sh
protoc --proto_path=proto proto/*.proto --go_out=plugins=grpc:pb
```
