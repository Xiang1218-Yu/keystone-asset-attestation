## Bug 是什么
本题对应的业务异常是：业务方在审核页面整理查询结果里的阶段历史后，存储中的历史也被同步改写，下一次页面打开会看到一条不存在的推进记录。请调整历史 slice 的生命周期隔离，避免查询结果反向改变内部记录。

## 如何触发
在含 Bug 的工作区执行下面的目标复现指令，观察接口结果或并发运行时行为是否与业务预期不一致。

## 根因
文件：internal/api/records.go、internal/api/server.go、internal/core/engine.go；符号：(*Server).getRecord、(*Engine).Get、cloneRecord；失效机制：记录历史返回值没有和存储对象完全隔离，页面整理查询结果时会写入内部历史集合。

## 运行指令
```bash
go test -v ./internal/core -run '^TestBug022HistorySliceDoesNotPolluteRecord$' -count=1
```

## 错误信息
目标行为在含 Bug 快照中失败，正常回归行为应保持通过；以下保留本题实际运行产生的原始输出片段。

## 错误堆栈
```text
=== RUN   TestBug022HistorySliceDoesNotPolluteRecord
    bug022_history_slice_does_not_pollute_record_test.go:18: history mutation polluted stored record
--- FAIL: TestBug022HistorySliceDoesNotPolluteRecord (0.00s)
FAIL
FAIL	keystone-asset-attestation/internal/core	0.740s
FAIL
```
