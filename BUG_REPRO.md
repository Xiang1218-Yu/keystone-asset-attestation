## Bug 是什么
本题对应的业务异常是：运营人员在资产校验页面提交一份触发领域规则错误的内容，接口把调用方输入问题报成服务器故障，监控中会混入不该告警的失败。请定位这条错误传播链，代码不要动。

## 如何触发
在含 Bug 的工作区执行下面的目标复现指令，观察接口结果或并发运行时行为是否与业务预期不一致。

## 根因
文件：internal/api/codec.go、internal/api/validate.go、internal/core/engine.go；符号：decode、(*Server).validate、(*Engine).ValidatePayload；失效机制：领域规则拒绝错误在跨层传播中被重新包装，输入问题因此被错误归类为服务器故障。

## 运行指令
```bash
go test -v ./internal/api -run '^TestBug008ValidationErrorKeepsBadRequestStatus$' -count=1
```

## 错误信息
目标行为在含 Bug 快照中失败，正常回归行为应保持通过；以下保留本题实际运行产生的原始输出片段。

## 错误堆栈
```text
=== RUN   TestBug008ValidationErrorKeepsBadRequestStatus
2026/08/24 13:36:37 INFO request method=POST path=/v1/validate duration=170.667µs
    bug008_validation_error_keeps_bad_request_status_test.go:57: validation error status = 500, want 400
--- FAIL: TestBug008ValidationErrorKeepsBadRequestStatus (0.00s)
FAIL
FAIL	keystone-asset-attestation/internal/api	0.498s
FAIL
```
