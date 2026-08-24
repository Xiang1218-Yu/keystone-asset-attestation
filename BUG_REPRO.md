## Bug 是什么
本题对应的业务异常是：列表读取返回后，页面调用方会对其中一条记录的历史做展示处理，后续单条查询却出现相同改动，说明列表结果没有和存储状态隔开。需要修复列表返回的 slice 共享问题，保证列表编辑不会污染单条数据。

## 如何触发
在含 Bug 的工作区执行下面的目标复现指令，观察接口结果或并发运行时行为是否与业务预期不一致。

## 根因
文件：internal/api/records.go、internal/api/server.go、internal/core/engine.go；符号：(*Server).listRecords、(*Engine).List、cloneRecord；失效机制：列表查询返回的记录仍共享内部历史 slice，调用方的展示处理会污染之后的单条查询。

## 运行指令
```bash
go test -v ./internal/core -run '^TestBug023ListSliceDoesNotPolluteRecord$' -count=1
```

## 错误信息
目标行为在含 Bug 快照中失败，正常回归行为应保持通过；以下保留本题实际运行产生的原始输出片段。

## 错误堆栈
```text
=== RUN   TestBug023ListSliceDoesNotPolluteRecord
    bug023_list_slice_does_not_pollute_record_test.go:19: list result polluted stored record
--- FAIL: TestBug023ListSliceDoesNotPolluteRecord (0.00s)
FAIL
FAIL	keystone-asset-attestation/internal/core	0.487s
FAIL
```
