## Bug 是什么
本题对应的业务异常是：供应商资产集中导入期间，多个采集端可能同时上报同一个资产标识；高峰时偶尔会看到两次创建都返回成功，后续读取只保留其中一份状态。这边需要修复并发创建的检查与写入一致性，让重复提交只能形成一份记录。

## 如何触发
在含 Bug 的工作区执行下面的目标复现指令，观察接口结果或并发运行时行为是否与业务预期不一致。

## 根因
文件：internal/api/records.go、internal/api/server.go、internal/core/engine.go；符号：(*Server).createRecord、(*Engine).Create；失效机制：重复检查与记录写入不在同一个互斥临界区，并发请求可以同时通过检查，最终产生多个成功响应和一次覆盖写入。

## 运行指令
```bash
go test -v ./internal/api -run '^TestBug013ConcurrentCreateKeepsSingleRecord$' -race -count=20
```

## 错误信息
目标行为在含 Bug 快照中失败，正常回归行为应保持通过；以下保留本题实际运行产生的原始输出片段。

## 错误堆栈
```text
=== RUN   TestBug013ConcurrentCreateKeepsSingleRecord
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=4.453417ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=4.445584ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=4.554042ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=4.4655ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=5.100917ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=4.49725ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=4.447584ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=4.548ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=4.515209ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=4.445125ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=4.535084ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=4.663ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=4.571333ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=4.410041ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=5.115417ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=4.641833ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=7.101875ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=12.061334ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=15.662917ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=18.608208ms
    bug013_concurrent_create_keeps_single_record_test.go:74: concurrent create successes = 4, want 1
--- FAIL: TestBug013ConcurrentCreateKeepsSingleRecord (0.02s)
=== RUN   TestBug013ConcurrentCreateKeepsSingleRecord
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=3.071875ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=3.075292ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=3.0535ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=3.19125ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=3.186417ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=3.157917ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=3.167833ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=3.236292ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=3.198667ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=3.148916ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=3.236583ms
... output truncated for document readability ...
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=2.183875ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=2.187833ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=2.078584ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=2.171875ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=1.983875ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=1.986084ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=1.959375ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=1.989167ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=2.01975ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=1.978292ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=1.96775ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=1.959666ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=1.983167ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=2.23375ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=2.003542ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=2.820667ms
--- PASS: TestBug013ConcurrentCreateKeepsSingleRecord (0.00s)
=== RUN   TestBug013ConcurrentCreateKeepsSingleRecord
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=2.013083ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=1.978958ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=2.0115ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=2.147542ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=1.96775ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=2.07875ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=2.0085ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=2.166042ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=1.965875ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=2.180667ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=2.122584ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=2.073875ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=2.1945ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=1.951625ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=2.044875ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=1.969959ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=2.028583ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=2.14275ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=2.160917ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=2.717541ms
--- PASS: TestBug013ConcurrentCreateKeepsSingleRecord (0.00s)
=== RUN   TestBug013ConcurrentCreateKeepsSingleRecord
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=1.61275ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=1.673375ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=1.678792ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=1.660542ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=1.683291ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=1.714959ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=1.667292ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=1.611833ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=1.6805ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=1.666916ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=1.648667ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=1.629458ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=1.641458ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=1.670041ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=1.664125ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=1.667959ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=1.708042ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=2.237541ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=4.484125ms
2026/08/24 13:36:41 INFO request method=POST path=/v1/records duration=6.202459ms
    bug013_concurrent_create_keeps_single_record_test.go:74: concurrent create successes = 3, want 1
--- FAIL: TestBug013ConcurrentCreateKeepsSingleRecord (0.01s)
FAIL
FAIL	keystone-asset-attestation/internal/api	0.649s
FAIL
```
