## Bug 是什么
本题对应的业务异常是：审核员推进记录失败后，记录推进失败的响应只剩下一个笼统提示，调用方无法知道是阶段不被领域规则接受还是记录本身不存在。请先定位错误传播被截断的位置，代码不要动。

## 如何触发
在含 Bug 的工作区执行下面的目标复现指令，观察接口结果或并发运行时行为是否与业务预期不一致。

## 根因
文件：internal/api/records.go、internal/api/server.go、internal/core/engine.go；符号：(*Server).advanceRecord、(*Engine).Advance、writeError；失效机制：推进失败的领域错误在服务层被压缩成笼统提示，阶段非法和记录不存在两类结果无法被调用方区分。

## 运行指令
```bash
go test -v ./internal/api -run '^TestBug012AdvancePreservesDomainError$' -count=1
```

## 错误信息
目标行为在含 Bug 快照中失败，正常回归行为应保持通过；以下保留本题实际运行产生的原始输出片段。

## 错误堆栈
```text
=== RUN   TestBug012AdvancePreservesDomainError
2026/08/24 13:36:40 INFO request method=POST path=/v1/records duration=1.032625ms
2026/08/24 13:36:40 INFO request method=POST path=/v1/records/advance-body/advance duration=15.417µs
    bug012_advance_preserves_domain_error_test.go:59: domain error lost: 
--- FAIL: TestBug012AdvancePreservesDomainError (0.00s)
FAIL
FAIL	keystone-asset-attestation/internal/api	0.494s
FAIL
```
