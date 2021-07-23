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
  ```

- clang-format：

  ```sh
  sudo apt-get install clang-format
  ```

- 依赖：

  ```sh
  go get -u google.golang.org/grpc
  go get -u google.golang.org/protobuf
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
