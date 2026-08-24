## Bug 是什么
本题对应的业务异常是：外部系统向记录接口提交内容损坏的请求，接口把解码失败当成服务端错误，调用方无法区分是请求需要重发还是服务需要排查。请先分析错误状态为何被改写，代码不要修改。

## 如何触发
在含 Bug 的工作区执行下面的目标复现指令，观察接口结果或并发运行时行为是否与业务预期不一致。

## 根因
文件：internal/api/codec.go、internal/api/records.go、internal/api/server.go；符号：decode、(*Server).createRecord、writeError；失效机制：请求体解码错误离开解码层后被统一成服务端错误，调用方失去了区分请求问题和服务问题所需的状态。

## 运行指令
```bash
go test -v ./internal/api -run '^TestBug009MalformedBodyKeepsBadRequestStatus$' -count=1
```

## 错误信息
目标行为在含 Bug 快照中失败，正常回归行为应保持通过；以下保留本题实际运行产生的原始输出片段。

## 错误堆栈
```text
=== RUN   TestBug009MalformedBodyKeepsBadRequestStatus
2026/08/24 13:36:38 INFO request method=POST path=/v1/records duration=292.584µs
    bug009_malformed_body_keeps_bad_request_status_test.go:57: malformed body status = 500, want 400
--- FAIL: TestBug009MalformedBodyKeepsBadRequestStatus (0.00s)
FAIL
FAIL	keystone-asset-attestation/internal/api	0.509s
FAIL
```
