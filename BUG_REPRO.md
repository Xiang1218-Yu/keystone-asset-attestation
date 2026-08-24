## Bug 是什么
本题对应的业务异常是：审核事件高峰期，多个事件同时写入审计日志，编号偶尔重复，后续按编号追踪会把不同事件误认为同一条。请先修复审计写入的并发分配逻辑，再保证普通审计查询不受影响。

## 如何触发
在含 Bug 的工作区执行下面的目标复现指令，观察接口结果或并发运行时行为是否与业务预期不一致。

## 根因
文件：internal/ops/audit.go、internal/ops/metrics.go、internal/ops/queue.go；符号：(*AuditLog).Add、(*Metrics).Inc、(*Queue).worker；失效机制：审计编号分配在并发事件写入时没有形成原子状态更新，多个 worker 可以读取同一编号并写入重复事件。

## 运行指令
```bash
go test -v ./internal/ops -run '^TestBug016ConcurrentAuditIdsStayUnique$' -race -count=20
```

## 错误信息
目标行为在含 Bug 快照中失败，正常回归行为应保持通过；以下保留本题实际运行产生的原始输出片段。

## 错误堆栈
```text
=== RUN   TestBug016ConcurrentAuditIdsStayUnique
    bug016_concurrent_audit_ids_stay_unique_test.go:23: duplicate audit id: audit-1000000002-1
--- FAIL: TestBug016ConcurrentAuditIdsStayUnique (0.00s)
=== RUN   TestBug016ConcurrentAuditIdsStayUnique
    bug016_concurrent_audit_ids_stay_unique_test.go:23: duplicate audit id: audit-1000000002-1
--- FAIL: TestBug016ConcurrentAuditIdsStayUnique (0.00s)
=== RUN   TestBug016ConcurrentAuditIdsStayUnique
    bug016_concurrent_audit_ids_stay_unique_test.go:23: duplicate audit id: audit-1000000002-2
--- FAIL: TestBug016ConcurrentAuditIdsStayUnique (0.00s)
=== RUN   TestBug016ConcurrentAuditIdsStayUnique
    bug016_concurrent_audit_ids_stay_unique_test.go:23: duplicate audit id: audit-1000000002-4
--- FAIL: TestBug016ConcurrentAuditIdsStayUnique (0.00s)
=== RUN   TestBug016ConcurrentAuditIdsStayUnique
    bug016_concurrent_audit_ids_stay_unique_test.go:23: duplicate audit id: audit-1000000002-2
--- FAIL: TestBug016ConcurrentAuditIdsStayUnique (0.00s)
=== RUN   TestBug016ConcurrentAuditIdsStayUnique
    bug016_concurrent_audit_ids_stay_unique_test.go:23: duplicate audit id: audit-1000000002-3
--- FAIL: TestBug016ConcurrentAuditIdsStayUnique (0.00s)
=== RUN   TestBug016ConcurrentAuditIdsStayUnique
    bug016_concurrent_audit_ids_stay_unique_test.go:23: duplicate audit id: audit-1000000002-2
--- FAIL: TestBug016ConcurrentAuditIdsStayUnique (0.00s)
=== RUN   TestBug016ConcurrentAuditIdsStayUnique
    bug016_concurrent_audit_ids_stay_unique_test.go:23: duplicate audit id: audit-1000000002-2
--- FAIL: TestBug016ConcurrentAuditIdsStayUnique (0.00s)
=== RUN   TestBug016ConcurrentAuditIdsStayUnique
    bug016_concurrent_audit_ids_stay_unique_test.go:23: duplicate audit id: audit-1000000002-2
--- FAIL: TestBug016ConcurrentAuditIdsStayUnique (0.00s)
=== RUN   TestBug016ConcurrentAuditIdsStayUnique
    bug016_concurrent_audit_ids_stay_unique_test.go:23: duplicate audit id: audit-1000000002-2
--- FAIL: TestBug016ConcurrentAuditIdsStayUnique (0.00s)
=== RUN   TestBug016ConcurrentAuditIdsStayUnique
    bug016_concurrent_audit_ids_stay_unique_test.go:23: duplicate audit id: audit-1000000002-1
--- FAIL: TestBug016ConcurrentAuditIdsStayUnique (0.00s)
=== RUN   TestBug016ConcurrentAuditIdsStayUnique
    bug016_concurrent_audit_ids_stay_unique_test.go:23: duplicate audit id: audit-1000000002-4
--- FAIL: TestBug016ConcurrentAuditIdsStayUnique (0.00s)
=== RUN   TestBug016ConcurrentAuditIdsStayUnique
    bug016_concurrent_audit_ids_stay_unique_test.go:23: duplicate audit id: audit-1000000002-4
--- FAIL: TestBug016ConcurrentAuditIdsStayUnique (0.00s)
=== RUN   TestBug016ConcurrentAuditIdsStayUnique
    bug016_concurrent_audit_ids_stay_unique_test.go:23: duplicate audit id: audit-1000000002-2
--- FAIL: TestBug016ConcurrentAuditIdsStayUnique (0.00s)
=== RUN   TestBug016ConcurrentAuditIdsStayUnique
    bug016_concurrent_audit_ids_stay_unique_test.go:23: duplicate audit id: audit-1000000002-2
--- FAIL: TestBug016ConcurrentAuditIdsStayUnique (0.00s)
=== RUN   TestBug016ConcurrentAuditIdsStayUnique
    bug016_concurrent_audit_ids_stay_unique_test.go:23: duplicate audit id: audit-1000000002-3
--- FAIL: TestBug016ConcurrentAuditIdsStayUnique (0.00s)
=== RUN   TestBug016ConcurrentAuditIdsStayUnique
    bug016_concurrent_audit_ids_stay_unique_test.go:23: duplicate audit id: audit-1000000002-2
--- FAIL: TestBug016ConcurrentAuditIdsStayUnique (0.00s)
=== RUN   TestBug016ConcurrentAuditIdsStayUnique
    bug016_concurrent_audit_ids_stay_unique_test.go:23: duplicate audit id: audit-1000000002-3
--- FAIL: TestBug016ConcurrentAuditIdsStayUnique (0.00s)
=== RUN   TestBug016ConcurrentAuditIdsStayUnique
    bug016_concurrent_audit_ids_stay_unique_test.go:23: duplicate audit id: audit-1000000002-2
--- FAIL: TestBug016ConcurrentAuditIdsStayUnique (0.00s)
=== RUN   TestBug016ConcurrentAuditIdsStayUnique
    bug016_concurrent_audit_ids_stay_unique_test.go:23: duplicate audit id: audit-1000000002-3
--- FAIL: TestBug016ConcurrentAuditIdsStayUnique (0.00s)
FAIL
FAIL	keystone-asset-attestation/internal/ops	0.523s
FAIL
```
