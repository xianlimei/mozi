Feature

这个平台其实可以演变以下几种服务：

1. 基于golang `plugin` 机制的任务系统。这是下面各个服务的基础。
2. Serverless。server 是通过 go `plugin` 机制载入的。
3. ApiGateway。一个请求相当于一个任务，这里就非常灵活了。可以实现 `group API`，即任务分割下分。