# gojober

job worker

## 设计思路

Jober 服务化后端启动，监控特定的目录下的`源文件`或者 `.so`文件，动态载入新的 `job` 逻辑。

对外提供统一的 `job` 上传接口，`HTTP` 或者 `gRPC`。上传的数据（**约定**）对应的 `struct` 为：

``` golang
type JobArgs struct {
    Name string `json:"name"` // job name --> .so 处理单元
    Args []byte `json:"args"` // job 处理逻辑的输入参数，是一个 json encoded 二进制数据
}
```

`Jober` 解析出 `Name` 并查找对应的 `plugin`，将 `[]byte` 传递给该 `plugin` 的执行单元 `Run` 函数，开始 `job` 的处理。

> 这里先简化逻辑，实际上是接受到 `job` 后，先将参数存到 redis。另一个线程去取这些 job 并处理。


在这其中，`Args` 的 `[]byte` 中包含的有效 `struct` 数据是业务逻辑中已经设计好的，通常 `.so` 中会有相应的 `struct` 类型用于 `decode`。


## TODO
[ ] 编写 server
[ ] job 本地化 (参考 nsq 的方案)
[ ] job 统计
[ ] 补全单元测试
[ ] 编写 client