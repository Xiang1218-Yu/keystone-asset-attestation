## Bug 是什么
本题对应的业务异常是：调用方取消创建请求后，接口仍然返回成功，后续查询还能看到一条本不该写入的记录，提示报错：2026/08/24 10:33:33 INFO request method=POST path=/v1/records duration=550.042µs；cancelled create status = 201, want 400；--- FAIL ---；FAIL；FAIL keystone-asset-attestation/internal/api 1.117s；FAIL。请定位取消信号为什么没有阻止这次跨层处理，代码不要动。

## 如何触发
在含 Bug 的工作区执行下面的目标复现指令，观察接口结果或并发运行时行为是否与业务预期不一致。

## 根因
文件：internal/api/records.go、internal/api/server.go、internal/core/engine.go；符号：(*Server).createRecord、(*Engine).Create、contextError；失效机制：请求上下文没有沿 HTTP、服务和引擎链路持续传递，取消后的请求仍进入写入路径，导致状态变更没有被及时阻断。

## 运行指令
```bash
go test -v ./internal/api -run '^TestBug001CreateCancellationStopsMutation$' -count=1
```

## 错误信息
目标行为在含 Bug 快照中失败，正常回归行为应保持通过；以下保留本题实际运行产生的原始输出片段。

## 错误堆栈
```text
[cancel_create] returncode=1
=== RUN   TestBug001CreateCancellationStopsMutation
2026/08/24 10:33:33 INFO request method=POST path=/v1/records duration=550.042µs
    bug001_create_cancellation_test.go:64: cancelled create status = 201, want 400
--- FAIL: TestBug001CreateCancellationStopsMutation (0.00s)
FAIL
FAIL	keystone-asset-attestation/internal/api	1.117s
FAIL
```
