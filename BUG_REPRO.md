## Bug 是什么
本题对应的业务异常是：审计人员读取创建结果并在页面上补充证据备注时，修改其中一条证据，下一次查询发现历史记录里的证据也被改掉，审计页面无法复盘最初结果。需要修复证据 slice 的返回隔离，让展示层只修改自己的副本。

## 如何触发
在含 Bug 的工作区执行下面的目标复现指令，观察接口结果或并发运行时行为是否与业务预期不一致。

## 根因
文件：internal/api/codec.go、internal/api/records.go、internal/core/engine.go；符号：(*Server).createRecord、(*Engine).Create、cloneRecord；失效机制：创建结果中的证据集合仍与内部记录共享 slice，展示层修改返回值后会反向改写历史记录。

## 运行指令
```bash
go test -v ./internal/core -run '^TestBug021EvidenceSliceDoesNotPolluteRecord$' -count=1
```

## 错误信息
目标行为在含 Bug 快照中失败，正常回归行为应保持通过；以下保留本题实际运行产生的原始输出片段。

## 错误堆栈
```text
=== RUN   TestBug021EvidenceSliceDoesNotPolluteRecord
    bug021_evidence_slice_does_not_pollute_record_test.go:18: evidence mutation polluted stored record
--- FAIL: TestBug021EvidenceSliceDoesNotPolluteRecord (0.00s)
FAIL
FAIL	keystone-asset-attestation/internal/core	0.483s
FAIL
```
