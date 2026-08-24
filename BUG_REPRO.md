## Bug 是什么
本题对应的业务异常是：监控系统探测模块接口遇到异常时，服务直接把 panic 冒泡到 HTTP 调用方，进程级恢复机制没有发挥作用。请先修复接口层的 panic 隔离，再确保正常模块请求仍能返回。

## 如何触发
在含 Bug 的工作区执行下面的目标复现指令，观察接口结果或并发运行时行为是否与业务预期不一致。

## 根因
文件：internal/api/middleware.go、internal/api/modules.go、internal/api/server.go；符号：recoverPanic、requestLog、(*Server).modules；失效机制：panic 恢复层没有在 handler 生命周期内隔离异常，未处理的 panic 直接穿过 HTTP 调用链。

## 运行指令
```bash
go test -v ./internal/api -run '^TestBug026ModulePanicStaysIsolated$' -count=1
```

## 错误信息
目标行为在含 Bug 快照中失败，正常回归行为应保持通过；以下保留本题实际运行产生的原始输出片段。

## 错误堆栈
```text
=== RUN   TestBug026ModulePanicStaysIsolated
    bug026_module_panic_stays_isolated_test.go:18: panic escaped handler: runtime error: invalid memory address or nil pointer dereference
--- FAIL: TestBug026ModulePanicStaysIsolated (0.00s)
FAIL
FAIL	keystone-asset-attestation/internal/api	0.595s
FAIL
```
