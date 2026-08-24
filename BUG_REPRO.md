## Bug 是什么
本题对应的业务异常是：多人协同审核同一批资产时，两个审核端可能同时把同一条记录推进到同一个阶段，返回成功的次数和版本号对不上，历史记录也会丢失。请把并发推进的状态更新处理好，保证一次有效转换只有一个结果。

## 如何触发
在含 Bug 的工作区执行下面的目标复现指令，观察接口结果或并发运行时行为是否与业务预期不一致。

## 根因
文件：internal/api/records.go、internal/api/server.go、internal/core/engine.go；符号：(*Server).advanceRecord、(*Engine).Advance；失效机制：推进的读取、阶段判断和历史写入之间存在竞态，多个审核请求可能基于同一旧快照提交相同转换。

## 运行指令
```bash
go test -v ./internal/api -run '^TestBug014ConcurrentAdvanceKeepsSingleTransition$' -race -count=20
```

## 错误信息
目标行为在含 Bug 快照中失败，正常回归行为应保持通过；以下保留本题实际运行产生的原始输出片段。

## 错误堆栈
```text
=== RUN   TestBug014ConcurrentAdvanceKeepsSingleTransition
2026/08/24 13:36:42 INFO request method=POST path=/v1/records duration=6.607917ms
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=170.416µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=158.542µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=317.583µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=266.25µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=366.625µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=35.042µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=351.5µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=92.584µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=97.792µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=366.75µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=147.625µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=457.917µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=255.917µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=138.708µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=456.583µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=537.834µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=379.042µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=281.917µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=1.901208ms
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=2.468875ms
    bug014_concurrent_advance_keeps_single_transition_test.go:75: concurrent advance successes = 2, want 1
--- FAIL: TestBug014ConcurrentAdvanceKeepsSingleTransition (0.01s)
=== RUN   TestBug014ConcurrentAdvanceKeepsSingleTransition
2026/08/24 13:36:42 INFO request method=POST path=/v1/records duration=5.758625ms
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=60.458µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=82.792µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=64.459µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=96.791µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=63.25µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=99.459µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=32.791µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=133.375µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=71.25µs
... output truncated for document readability ...
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=98.875µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=167.25µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=137.083µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=250.375µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=276.333µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=290µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=186.875µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=290.75µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=334.417µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=368µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=320.417µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=368.667µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=1.013333ms
--- PASS: TestBug014ConcurrentAdvanceKeepsSingleTransition (0.00s)
=== RUN   TestBug014ConcurrentAdvanceKeepsSingleTransition
2026/08/24 13:36:42 INFO request method=POST path=/v1/records duration=2.266958ms
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=34.334µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=51.292µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=125.875µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=53.458µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=28.75µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=41.709µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=82.959µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=16.75µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=42.333µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=96.709µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=146.208µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=204.833µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=169.667µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=207.875µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=198.75µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=226.917µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=176.5µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=200.709µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=711.25µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=741µs
    bug014_concurrent_advance_keeps_single_transition_test.go:75: concurrent advance successes = 2, want 1
--- FAIL: TestBug014ConcurrentAdvanceKeepsSingleTransition (0.00s)
=== RUN   TestBug014ConcurrentAdvanceKeepsSingleTransition
2026/08/24 13:36:42 INFO request method=POST path=/v1/records duration=2.609916ms
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=18.167µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=30.666µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=32.875µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=64.125µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=41.875µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=102.542µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=158.333µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=34.792µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=171.542µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=150.667µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=113.084µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=194µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=47.541µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=77.25µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=227.791µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=172.5µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=152.5µs
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=1.385459ms
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=1.451833ms
2026/08/24 13:36:42 INFO request method=POST path=/v1/records/same-stage/advance duration=1.547833ms
    bug014_concurrent_advance_keeps_single_transition_test.go:75: concurrent advance successes = 3, want 1
--- FAIL: TestBug014ConcurrentAdvanceKeepsSingleTransition (0.00s)
FAIL
FAIL	keystone-asset-attestation/internal/api	0.627s
FAIL
```
