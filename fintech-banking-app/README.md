# Fintech Banking App

[原文地址](https://dev.to/duomly/series/6782)

## 工具

- migrate: migration cli

  ```sh
  go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
  ```

- sqlc: sql generate cli

  ```sh
  go get github.com/kyleconroy/sqlc/cmd/sqlc
  ```

- gomock: mock db

  ```sh
  go get github.com/golang/mock/mockgen
  ```

  ```sh
  vi ~/.zshrc
  ```

  文件内容：

  ```sh
  export PATH=$PATH:~/go/bin
  ```

  执行

  ```sh
  source ~/.zshrc
  which mockgen
  ```

  ```sh
  mockgen fintech-banking-app/db/sqlc Store
  ```

## 依赖

- pq: postgres driver

  ```sh
  go get -u github.com/lib/pq
  ```

- testify: test

  ```sh
  go get github.com/stretchr/testify
  ```

- gin: web framework

  ```sh
  go get github.com/gin-gonic/gin
  ```

- viper: env file

  ```sh
  go get github.com/spf13/viper
  ```

## 隔离等级与读取现象

1. 数据库事务的 ACID 属性：

   - `Atomicity` 原子性
   - `Consistency` 持续性
   - `Isolation` 隔离性：作为数据库事务属性的最高等级，用于确保所有的并发事务不会被互相影响
   - `Durability` 持久性

1. 4 种读取现象

   - `dirty read`：发生在当一个事务读取一些由其它并发事务所写的数据，并且还未被提交时。这很不好，因为我们不知道其它事务最终是否会提交或者回滚。因此我们可能会读取到错误的数据，当回滚出现时。
   - `non-repeatable read`：当一个事务读取相同的记录两次并发现不同的值，因为该行在第一次读取后被其它事务所修改了。
   - `phantom read`：与前者类似，影响的是查找多行而不是单行。这种情况下，相同的查询被重复执行后，返回的却是不同的行，这是因为其它的事务近期做了一些改变，例如插入新行或是删除已有行。
   - `serialization anomaly`：一组并发提交的事务，如果我们尝试以任何顺序运行并且它们之间没有重叠，它们不能被完成。

1. 4 种隔离等级

   - `read uncommitted`：最低等级的隔离。该等级下的事务可以看到其它事务的未被提交的数据，因此允许 `dirty read` 现象出现。
   - `read committed`：事务只能看到其它事务提交后的数据。因此 `dirty read` 不可能出现。
   - `repeatable read`：确保相同的查询请求将会返回同样的结果，无论它被执行了多少次，即使有其它的并发事务提交了新的改动。
   - `serializable`：最高等级的隔离。并发事务运行在这个等级被保证了：如果它们是以某种顺序运行并且之间没用重叠，那么返回的结果就是相同的。基本上就意味着存在至少一种方法来排序这些并发事务，使得它们的结果与一个接着一个的运行的结果相同。

1. 结论：

   - MySql 的 隔离性：![isolation_levels_in_mysql](./img/isolation_levels_in_mysql.png)
   - PostgreSql 的 隔离性：![isolation_levels_in_postgresql](./img/isolation_levels_in_postgresql.png)
   - 两者间的比较：![compare_mysql_vs_postgres](./img/compare_mysql_vs_postgres.png)

1. 总结：需要记住的最重要的一点就是当使用高隔离等级可能会出现错误、超时、或者甚至死锁。因此我们需要小心的为我们的事务实现一个重试机制。另外不同的数据库引擎可能有不同的隔离等级实现。所以确保你小心的阅读文档，并在写代码之前独自尝试一下。

## 设置 Github Actions 使 Go + Postgres 运行自动测试

**Continuos integration (CI)**是软件开发过程中的一个重要部分，它集成了由团队合作所造成共享代码仓库中的持续性的改动。为了确保高质量的代码并且减少潜在错误，每个集成通常都会被自动构建和测试的工具所验证。

### Workflow

为了使用 Github Actions，我们需要定义一个工作流 workflow。它基本上是由一个或者多个工作而组成的一个自动化过程。它可以由以下三种方式触发：

- 通过一个 Github 仓库的事件
- 通过设定一个重复的时间表
- 通过 UI 手动点击

![github_actions_workflow](./img/github_actions_workflow.png)

### Runner

为了运行任务，我们必须为每个任务都指定一个运行器 runner。一个 runner 就是一个监听可行任务的简单的服务，并且它每次只会运行一个任务。我们可以直接使用 Github 的 runner，或者指定你自己的 runner。

![github_actions_runner](./img/github_actions_runner.png)

Runners 会运行任务，接着报告进程，日志以及返回的结果至 Github，这样我们便可以在 repo 的 UI 上轻松地查看。

### Job

一个任务 job 是一个将会被执行在同一个 runner 上的步骤的集合。通常所有在 workflow 的任务会并行执行，除非当你某些任务依赖于其它的任务，这时就会顺序运行。

![github_actions_job](./img/github_actions_job.png)

通过 `needs` 关键字告诉 `test` 任务是依赖于 `build` 任务的，这样就能在成功构建后再进行进行测试。

### Step

步骤 steps 是独立的顺序运行的任务 tasks，在一个 Job 中一个接着一个。一个 step 可以包含一个或者多个行动 actions。

![github_actions_step](./img/github_actions_step.png)

Action 基本上是一个独立的命令，像是一个 `test_server.sh` 的脚本。如果一个步骤中包含多个 actions，它们则会被顺序运行。

一个有趣的事情就是 action 是可以被重复使用的。因此如果某人已经编写了一个我们所需的 github action，我们实际上是可以将其使用于我们的工作流中。

### Summary

![github_actions_summary](./img/github_actions_summary.png)

## Mock DB

使用 Mock DB 的几个原因：

- 首先，帮助我们可以更简单的编写出独立的测试，因为每个测试将会用到隔离的 mock DB 来存储数据，即消除了测试之间的相互影响。如果使用的是真实 DB，所有的测试都会在同一个地方读写数据，这样就很难避免冲突，特别是一个大的项目中携带者巨大的基础代码。
- 其次，测试将会更快的运行，因为它们不再需要花时间与 DB 沟通并等待回应。所有的行为将在内存与统一过程中被执行。
- 第三也是最重要的原因就是：它让我们编写 100% 覆盖率的测试。通过一个 mock DB，我们可以轻松的创造并测试一些边缘情况，例如一个意外的错误，或者一个连接失败，这些都是在使用真实 DB 情况中不能被达到的。

### 生成 Mock DB

`mockgen` 给予用户两种生成 mocks 的方式。`source mode` 将会从一个源文件生成 mock 接口。

当源文件引入了其它包的时候这将变得更复杂，这也是我们在真实项目中会遇到的情况。

这种情况下，更好地办法是使用 `reflect mode`，即我们仅需提供包名称以及接口，便可以让 mockgen 自动使用反射。

执行：

```sh
mockgen fintech-banking-app/db/sqlc Store
```

第一个参数为 `Store` 接口的导入路径。第二个参数为接口的名称。

```sh
mockgen -destination db/mock/store.go fintech-banking-app/db/sqlc Store
```

使用 `-destination` 用于指定生成的文件。

生成的 `db/mock/store.go` 文件中有两个重要的结构体：`MockStore` 与 `MockStoreMockRecorder`。

前者实现了 `Store` 接口所需要的所有函数。后者可以指定函数被调用的次数，以及带有何种参数。

由于生成的 package 名称为 `mock_sqlc` 并不符合惯用法，加上 `-package` 可以自定义命名。

```sh
mockgen -package mockdb -destination db/mock/store.go fintech-banking-app/db/sqlc Store
```

## 如何存储密码

添加 `cost` 和 `salt` 用于生成最终的哈希字符串：

![securely_store_password](./img/securely_store_password.png)

该哈希字符串由 4 个部分组成：

- 第一部分为 `hash algorithm identifier` 哈希算法标识。这里 `2A` 即 `bcrypt` 算法。
- 第二部分为 `cost`。本例中为 `10`，意味着将会有 `2^10 = 1024` 种钥匙解释。
- 第三部分为 `salt` 拥有 `16 bytes` 长度，或者 `128 bits`。它是由 `base64` 格式来进行编码，将会生成 `22` 个字符的字符串。
- 最后部分为 `24 bytes` 的哈希值，编码为 `31` 个字符。

这四部分会串联在一起形成一个哈希字符串，并最终将会在数据库中储存该哈希字符串。

![securely_store_password2](./img/securely_store_password2.png)

如何验证用户输入的密码是正确的呢？首先我们通过 `username` 在数据库中找到 `hashed_password`。接着使用 `hashed_password` 中的 `cost` 和 `salt` 作为入参来哈希用户输入进 `bcrypt` 中的 `naked_password`。其输出结果将会变为另一个哈希值。接着比较这两个哈希值，如果一致，则密码正确。

![securely_store_password3](./img/securely_store_password3.png)
