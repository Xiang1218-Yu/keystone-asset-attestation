## Bug 是什么
本题对应的业务异常是：运营看板刷新资产快照时，阶段统计数量会随着记录版本号放大，页面显示的记录数和实际列表数量对不上。请修复快照状态汇总的一致性，让版本推进不会重复放大统计。

## 如何触发
在含 Bug 的工作区执行下面的目标复现指令，观察接口结果或并发运行时行为是否与业务预期不一致。

## 根因
文件：internal/api/modules.go、internal/api/server.go、internal/core/engine.go；符号：(*Server).snapshot、(*Engine).Snapshot；失效机制：快照统计把记录版本推进产生的状态重复计入汇总，统计对象没有按记录实体维度保持一致。

## 运行指令
```bash
go test -v ./internal/api -run '^TestBug024SnapshotCountsRecordsNotVersions$' -count=1
```

## 错误信息
目标行为在含 Bug 快照中失败，正常回归行为应保持通过；以下保留本题实际运行产生的原始输出片段。

## 错误堆栈
```text
=== RUN   TestBug024SnapshotCountsRecordsNotVersions
2026/08/24 13:36:50 INFO request method=POST path=/v1/records duration=990.959µs
2026/08/24 13:36:50 INFO request method=POST path=/v1/records/snapshot-state/advance duration=109.708µs
2026/08/24 13:36:50 INFO request method=GET path=/v1/snapshot duration=9.917µs
    bug024_snapshot_counts_records_not_versions_test.go:60: snapshot count was version-weighted: 
--- FAIL: TestBug024SnapshotCountsRecordsNotVersions (0.00s)
FAIL
FAIL	keystone-asset-attestation/internal/api	0.564s
FAIL
```
