## Bug 是什么
本题对应的业务异常是：审计页面查询一条不存在的记录时，响应状态虽然看起来正常，但错误内容已经变成笼统提示，调用方丢失了资源不存在的上下文。请说明错误信息在哪一层被截断，代码不要改。

## 如何触发
在含 Bug 的工作区执行下面的目标复现指令，观察接口结果或并发运行时行为是否与业务预期不一致。

## 根因
文件：internal/api/records.go、internal/api/server.go、internal/core/engine.go；符号：(*Server).getRecord、(*Engine).Get、writeError；失效机制：记录不存在的错误信息在 API 返回阶段被截断，状态虽可返回但资源上下文没有穿过错误传播链。

## 运行指令
```bash
go test -v ./internal/api -run '^TestBug011MissingRecordPreservesErrorBody$' -count=1
```

## 错误信息
目标行为在含 Bug 快照中失败，正常回归行为应保持通过；以下保留本题实际运行产生的原始输出片段。

## 错误堆栈
```text
=== RUN   TestBug011MissingRecordPreservesErrorBody
2026/08/24 13:36:39 INFO request method=GET path=/v1/records/missing-body duration=212.459µs
    bug011_missing_record_preserves_error_body_test.go:57: missing body lost domain error: 
--- FAIL: TestBug011MissingRecordPreservesErrorBody (0.00s)
FAIL
FAIL	keystone-asset-attestation/internal/api	0.509s
FAIL
```
