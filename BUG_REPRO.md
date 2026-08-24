## Bug 是什么
本题对应的业务异常是：运维探针访问模块接口发生 panic 时，日志层的 defer 路径没有让恢复包装接住异常，调用方会直接收到未处理的崩溃。请修复中间件的退出顺序和恢复生命周期，并保持健康请求经过日志包装。

## 如何触发
在含 Bug 的工作区执行下面的目标复现指令，观察接口结果或并发运行时行为是否与业务预期不一致。

## 根因
文件：internal/api/middleware.go、internal/api/modules.go、internal/api/server.go；符号：requestLog、recoverPanic、(*Server).Handler；失效机制：日志和恢复中间件的 defer 生命周期与包装顺序不匹配，panic 发生时恢复包装没有稳定覆盖调用链。

## 运行指令
```bash
go test -v ./internal/api -run '^TestBug028RecoveryCoversRequestLog$' -count=1
```

## 错误信息
目标行为在含 Bug 快照中失败，正常回归行为应保持通过；以下保留本题实际运行产生的原始输出片段。

## 错误堆栈
```text
=== RUN   TestBug028RecoveryCoversRequestLog
    bug028_recovery_covers_request_log_test.go:17: panic escaped middleware: runtime error: invalid memory address or nil pointer dereference
--- FAIL: TestBug028RecoveryCoversRequestLog (0.00s)
FAIL
FAIL	keystone-asset-attestation/internal/api	0.554s
FAIL
```
