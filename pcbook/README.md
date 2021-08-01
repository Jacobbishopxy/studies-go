# PcBook

[原文地址](https://dev.to/techschoolguru/series/7311)

## Prerequisites

- 预安装 `protoc`（Ubuntu）：

  ```sh
  PROTOC_ZIP=protoc-3.17.3-linux-x86_64.zip

  curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.17.3/$PROTOC_ZIP

  sudo unzip -o $PROTOC_ZIP -d /usr/local bin/protoc

  sudo unzip -o $PROTOC_ZIP -d /usr/local 'include/*'

  rm -f $PROTOC_ZIP

  sudo chown -R $USER /usr/local/include/google/

  protoc --version
  ```

- clang-format：

  ```sh
  sudo apt-get install clang-format
  ```

- 安装 gRPC 客户端 evans，用于运行时的 grpc reflection：

  ```sh
  curl -OL https://github.com/ktr0731/evans/releases/download/0.10.0/evans_linux_amd64.tar.gz

  tar zxvf evans_linux_amd64.tar.gz

  sudo mv evans /usr/local/bin/

  rm evans_linux_amd64.tar.gz

  evans -v
  ```

- 依赖：

  ```sh
  go get -u google.golang.org/grpc

  go get -u google.golang.org/protobuf

  # 非常重要！不直接作用于本项目（go mod tidy 会清除）
  # global 环境缺少该库会导致 protoc 命令无法生成 golang 代码
  go get -u github.com/golang/protobuf/protoc-gen-go
  ```

## Note

- 由 `protoc` 生成 Golang 代码（已加入 Makefile）：

  ```sh
  protoc --proto_path=proto proto/*.proto --go_out=plugins=grpc:pb
  ```

- 由于 proto 文件不位于项目根目录，因此 proto 文件中的 import 会显示错误。修改 `.vscode/settings.json` 如下：

  ```json
  "protoc": {
      "options": ["--proto_path=./pcbook/proto"]
  }
  ```

## GRPC Interceptor

本质上 GRPC 拦截器与一个中间件函数相似，可以被添加至服务端以及客户端。

- 服务端的拦截器是在调用服务端 RPC 方法前的一个函数。可被用于多个目的，例如 logging，tracing，rate-limiting，authentication 以及 authorization。
- 同样的，客户端拦截器是在调用客户端的 RPC 方法前的一个函数。

![grpc-interceptor](./grpc-interceptor.png)

本文中我们将：

- 首先，实现服务端拦截器通过 JWT 用于验证 gRPC 的 APIs。通过该拦截器可以确保特定职能的用户可以调用服务端的指定 API。
- 其次，实现客户端拦截器用于用户登录以及调用 gRPC 的 API 前的 JWT 绑定。
