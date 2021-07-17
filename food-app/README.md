# Food App

[原文地址](https://dev.to/stevensunflash/using-domain-driven-design-ddd-in-golang-3ee5)

## Note

DDD 由四层构建：

- 领域 Domain：用于定义领域与应用的业务逻辑

- 基础 Infrastructure：包含应用中的所有独立项：外部库，数据库引擎等等

- 应用 Application：服务于领域层与接口层之间的中间层。由接口层发送请求至领域层，在领域层处理完之后返回响应回接口层。

- 接口 Interface：维护所有与外界系统的交互，例如 web 服务，RMI 接口或者是 web 应用，批量处理前端。

![food-app](./food-app.jpg)

### 领域层

1. Entity：用于定义所有事物的 Schema。添加辅助函数，类似于验证，密码加密等。

1. Repository：定义基础层所需要实现的抽象方法的集合。

### 基础层

具体实现定义于领域层中 repository 的接口（本文仅需考虑数据库交互）。

由具体的 Repositories 结构体来存储所有应用的 repositories，以及一个数据库实例。
