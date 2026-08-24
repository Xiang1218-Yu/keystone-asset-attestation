## Bug 是什么
本题对应的业务异常是：自动化导入任务没有携带操作者信息创建记录时，历史里的执行人变成空值，审计页面展示会触发空数据处理异常。需要修复操作者默认值从请求到记录历史的跨层传递，同时不能影响显式操作者。

## 如何触发
在含 Bug 的工作区执行下面的目标复现指令，观察接口结果或并发运行时行为是否与业务预期不一致。

## 运行指令
```bash
go test -v ./internal/api -run '^TestBug029AnonymousActorStaysSafe$' -count=1
```

## 错误信息
目标行为在含 Bug 快照中失败，正常回归行为应保持通过；以下保留本题实际运行产生的原始输出片段。

## 错误堆栈
```text
=== RUN   TestBug029AnonymousActorStaysSafe
2026/08/24 13:36:54 INFO request method=POST path=/v1/records duration=943.375µs
    bug029_anonymous_actor_stays_safe_test.go:58: anonymous actor missing: 
--- FAIL: TestBug029AnonymousActorStaysSafe (0.00s)
FAIL
FAIL	keystone-asset-attestation/internal/api	0.608s
FAIL
```
