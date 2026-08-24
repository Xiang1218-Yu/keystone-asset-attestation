## Bug 是什么
本题对应的业务异常是：资产校验客户端偶尔会发出请求体为空的请求，接口返回了与其他输入错误不同的状态，客户端重试策略因此走了另一条分支。请先定位错误码在解码层的变化，暂时不要改代码。

## 如何触发
在含 Bug 的工作区执行下面的目标复现指令，观察接口结果或并发运行时行为是否与业务预期不一致。

## 根因
文件：internal/api/codec.go、internal/api/server.go、internal/api/validate.go；符号：decode、(*Server).validate、writeError；失效机制：空请求体的解码错误与其他客户端输入错误没有保持一致的跨层状态语义，重试方因此走了错误分支。

## 运行指令
```bash
go test -v ./internal/api -run '^TestBug010EmptyBodyKeepsBadRequestStatus$' -count=1
```

## 错误信息
目标行为在含 Bug 快照中失败，正常回归行为应保持通过；以下保留本题实际运行产生的原始输出片段。

## 错误堆栈
```text
=== RUN   TestBug010EmptyBodyKeepsBadRequestStatus
2026/08/24 13:36:39 INFO request method=POST path=/v1/validate duration=245.833µs
    bug010_empty_body_keeps_bad_request_status_test.go:57: empty body status = 422, want 400
--- FAIL: TestBug010EmptyBodyKeepsBadRequestStatus (0.00s)
FAIL
FAIL	keystone-asset-attestation/internal/api	0.495s
FAIL
```
