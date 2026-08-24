## Bug 是什么
本题对应的业务异常是：审计页面拿到列表后会做本地展示处理，调用方修改返回结果后，下一次读取发现历史事件数量和内容都被外部代码污染。请修复审计列表的 slice 生命周期隔离，让页面处理只能作用于副本。

## 如何触发
在含 Bug 的工作区执行下面的目标复现指令，观察接口结果或并发运行时行为是否与业务预期不一致。

## 根因
文件：internal/ops/audit.go、internal/ops/metrics.go、internal/ops/queue.go；符号：(*AuditLog).List、(*Queue).worker、审计列表调用链；失效机制：审计列表直接返回内部 entries 的共享 slice，页面修改结果时会污染审计日志的内部状态。

## 运行指令
```bash
go test -v ./internal/ops -run '^TestBug030AuditListDoesNotPolluteLog$' -count=1
```

## 错误信息
目标行为在含 Bug 快照中失败，正常回归行为应保持通过；以下保留本题实际运行产生的原始输出片段。

## 错误堆栈
```text
=== RUN   TestBug030AuditListDoesNotPolluteLog
    bug030_audit_list_does_not_pollute_log_test.go:14: audit list mutation polluted log
--- FAIL: TestBug030AuditListDoesNotPolluteLog (0.00s)
FAIL
FAIL	keystone-asset-attestation/internal/ops	0.506s
FAIL
```
