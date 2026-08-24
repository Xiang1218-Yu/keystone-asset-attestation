## Bug 是什么
本题对应的业务异常是：外部导入服务提交记录时，请求体在解码前已经被关闭，正常数据被当成空请求，批量导入流程全部失败。需要调整请求体资源的 defer 生命周期，让解码完成前保持请求体可读。

## 如何触发
在含 Bug 的工作区执行下面的目标复现指令，观察接口结果或并发运行时行为是否与业务预期不一致。

## 根因
文件：internal/api/codec.go、internal/api/records.go、internal/api/server.go；符号：decode、(*Server).createRecord、request body 生命周期；失效机制：请求体关闭的 defer 生命周期早于解码过程，读取链路拿到空内容，正常导入数据因此被判为无效请求。

## 运行指令
```bash
go test -v ./internal/api -run '^TestBug027RequestBodyStaysReadable$' -count=1
```

## 错误信息
目标行为在含 Bug 快照中失败，正常回归行为应保持通过；以下保留本题实际运行产生的原始输出片段。

## 错误堆栈
```text
=== RUN   TestBug027RequestBodyStaysReadable
2026/08/24 13:36:52 INFO request method=POST path=/v1/records duration=237.875µs
    bug027_request_body_stays_readable_test.go:75: readable body status = 400, want 201
--- FAIL: TestBug027RequestBodyStaysReadable (0.00s)
FAIL
FAIL	keystone-asset-attestation/internal/api	0.525s
FAIL
```
