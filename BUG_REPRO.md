## Bug 是什么
本题对应的业务异常是：服务发布准备阶段需要停止后台队列，但正在处理的任务没有收到取消信号，关闭流程只能等待超时，发布期间会拖住服务退出。请修复任务处理和队列 context 的生命周期传递，让停止操作能及时结束任务。

## 如何触发
在含 Bug 的工作区执行下面的目标复现指令，观察接口结果或并发运行时行为是否与业务预期不一致。

## 根因
文件：cmd/keystone-asset-attestation/main.go、internal/ops/queue.go、internal/ops/metrics.go；符号：(*Queue).Stop、(*Queue).worker、main；失效机制：队列 worker 使用了不能被停止流程取消的上下文，关闭时正在执行的任务无法及时退出，服务只能等待超时。

## 运行指令
```bash
go test -v ./internal/ops -run '^TestBug018QueueStopCancelsRunningHandler$' -count=1
```

## 错误信息
目标行为在含 Bug 快照中失败，正常回归行为应保持通过；以下保留本题实际运行产生的原始输出片段。

## 错误堆栈
```text
=== RUN   TestBug018QueueStopCancelsRunningHandler
    bug018_queue_stop_cancels_running_handler_test.go:29: running handler did not receive queue cancellation
--- FAIL: TestBug018QueueStopCancelsRunningHandler (0.10s)
FAIL
FAIL	keystone-asset-attestation/internal/ops	0.592s
FAIL
```
