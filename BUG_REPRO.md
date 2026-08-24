## Bug 是什么
本题对应的业务异常是：记录正在推进阶段时客户端已经取消了请求，但服务端仍把阶段改掉，页面刷新后会看到一次没有完成的推进，提示报错：2026/08/24 10:33:35 INFO request method=POST path=/v1/records/cancel-advance/advance duration=60.708µs；cancelled advance status = 200, want 400；--- FAIL ---；FAIL；FAIL keystone-asset-attestation/internal/api 0.500s；FAIL。请说明取消生命周期在哪个环节失效，代码不要修改。

## 如何触发
在含 Bug 的工作区执行下面的目标复现指令，观察接口结果或并发运行时行为是否与业务预期不一致。

## 根因
文件：internal/api/records.go、internal/api/server.go、internal/core/engine.go；符号：(*Server).advanceRecord、(*Engine).Advance、contextError；失效机制：推进请求的取消状态在跨层调用中没有保持，阶段检查和状态写入继续使用未取消的执行路径。

## 运行指令
```bash
go test -v ./internal/api -run '^TestBug002AdvanceCancellationKeepsStage$' -count=1
```

## 错误信息
目标行为在含 Bug 快照中失败，正常回归行为应保持通过；以下保留本题实际运行产生的原始输出片段。

## 错误堆栈
```text
[cancel_advance] returncode=1
=== RUN   TestBug002AdvanceCancellationKeepsStage
2026/08/24 10:33:35 INFO request method=POST path=/v1/records duration=674.084µs
2026/08/24 10:33:35 INFO request method=POST path=/v1/records/cancel-advance/advance duration=60.708µs
    bug002_advance_cancellation_test.go:67: cancelled advance status = 200, want 400
--- FAIL: TestBug002AdvanceCancellationKeepsStage (0.00s)
FAIL
FAIL	keystone-asset-attestation/internal/api	0.500s
FAIL
```
