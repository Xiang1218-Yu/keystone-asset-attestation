## Bug 是什么
本题对应的业务异常是：测试环境的服务启动探针没有配置引擎，访问模块接口时直接触发空指针并被包装成通用错误，监控只能看到模糊失败。需要修复 nil 引擎场景下的运行时处理，同时保留正常引擎下的模块响应。

## 如何触发
在含 Bug 的工作区执行下面的目标复现指令，观察接口结果或并发运行时行为是否与业务预期不一致。

## 运行指令
```bash
go test -v ./internal/api -run '^TestBug025NilEngineIsRecovered$' -count=1
```

## 错误信息
目标行为在含 Bug 快照中失败，正常回归行为应保持通过；以下保留本题实际运行产生的原始输出片段。

## 错误堆栈
```text
=== RUN   TestBug025NilEngineIsRecovered
2026/08/24 13:36:51 INFO request method=GET path=/v1/modules duration=546.958µs
    bug025_nil_engine_is_recovered_test.go:20: nil engine status = 200, want 500
--- FAIL: TestBug025NilEngineIsRecovered (0.00s)
FAIL
FAIL	keystone-asset-attestation/internal/api	0.587s
FAIL
```
