## Bug 是什么
本题对应的业务异常是：多个后台 worker 同时处理资产任务时，会并发增加同一个运行指标，最终计数小于实际处理数，监控曲线因此低估了负载。请调整指标更新的并发安全性，使累计值不会因为同时写入而丢失。

## 如何触发
在含 Bug 的工作区执行下面的目标复现指令，观察接口结果或并发运行时行为是否与业务预期不一致。

## 根因
文件：internal/ops/audit.go、internal/ops/metrics.go、internal/ops/queue.go；符号：(*Metrics).Inc、(*Queue).worker；失效机制：指标读取和累加不是一个受保护的原子操作，并发 worker 的更新会互相覆盖，累计值低于真实处理数。

## 运行指令
```bash
go test -v ./internal/ops -run '^TestBug017ConcurrentMetricIncrementsAreCounted$' -race -count=20
```

## 错误信息
目标行为在含 Bug 快照中失败，正常回归行为应保持通过；以下保留本题实际运行产生的原始输出片段。

## 错误堆栈
```text
=== RUN   TestBug017ConcurrentMetricIncrementsAreCounted
==================
WARNING: DATA RACE
Read at 0x00c000091020 by goroutine 16:
  runtime.mapdelete_fast64()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/internal/runtime/maps/runtime_fast64.go:505 +0x8c
  keystone-asset-attestation/internal/ops.(*Metrics).Inc()
      /Users/tog/Documents/ChatGPT/新建go/2026-08-24/keystone-asset-attestation__017/.base_snapshot/internal/ops/metrics.go:18 +0x98
  keystone-asset-attestation/internal/ops.TestBug017ConcurrentMetricIncrementsAreCounted.func1()
      /Users/tog/Documents/ChatGPT/新建go/2026-08-24/keystone-asset-attestation__017/.base_snapshot/internal/ops/bug017_concurrent_metric_increments_are_counted_test.go:15 +0x1c

Previous write at 0x00c000091020 by goroutine 8:
  runtime.mapaccess2_faststr()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/internal/runtime/maps/runtime_faststr.go:161 +0x2ac
  keystone-asset-attestation/internal/ops.(*Metrics).Inc()
      /Users/tog/Documents/ChatGPT/新建go/2026-08-24/keystone-asset-attestation__017/.base_snapshot/internal/ops/metrics.go:20 +0xe0
  keystone-asset-attestation/internal/ops.TestBug017ConcurrentMetricIncrementsAreCounted.func1()
      /Users/tog/Documents/ChatGPT/新建go/2026-08-24/keystone-asset-attestation__017/.base_snapshot/internal/ops/bug017_concurrent_metric_increments_are_counted_test.go:15 +0x1c

Goroutine 16 (running) created at:
  keystone-asset-attestation/internal/ops.TestBug017ConcurrentMetricIncrementsAreCounted()
      /Users/tog/Documents/ChatGPT/新建go/2026-08-24/keystone-asset-attestation__017/.base_snapshot/internal/ops/bug017_concurrent_metric_increments_are_counted_test.go:15 +0xd8
  testing.tRunner()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2036 +0x164
  testing.(*T).Run.gowrap1()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2101 +0x34

Goroutine 8 (finished) created at:
  keystone-asset-attestation/internal/ops.TestBug017ConcurrentMetricIncrementsAreCounted()
      /Users/tog/Documents/ChatGPT/新建go/2026-08-24/keystone-asset-attestation__017/.base_snapshot/internal/ops/bug017_concurrent_metric_increments_are_counted_test.go:15 +0xd8
  testing.tRunner()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2036 +0x164
  testing.(*T).Run.gowrap1()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2101 +0x34
==================
... output truncated for document readability ...
Goroutine 18 (running) created at:
  keystone-asset-attestation/internal/ops.TestBug017ConcurrentMetricIncrementsAreCounted()
      /Users/tog/Documents/ChatGPT/新建go/2026-08-24/keystone-asset-attestation__017/.base_snapshot/internal/ops/bug017_concurrent_metric_increments_are_counted_test.go:15 +0xd8
  testing.tRunner()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2036 +0x164
  testing.(*T).Run.gowrap1()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2101 +0x34
==================
==================
WARNING: DATA RACE
Write at 0x00c000091020 by goroutine 34:
  runtime.mapaccess2_faststr()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/internal/runtime/maps/runtime_faststr.go:161 +0x2ac
  keystone-asset-attestation/internal/ops.(*Metrics).Inc()
      /Users/tog/Documents/ChatGPT/新建go/2026-08-24/keystone-asset-attestation__017/.base_snapshot/internal/ops/metrics.go:20 +0xe0
  keystone-asset-attestation/internal/ops.TestBug017ConcurrentMetricIncrementsAreCounted.func1()
      /Users/tog/Documents/ChatGPT/新建go/2026-08-24/keystone-asset-attestation__017/.base_snapshot/internal/ops/bug017_concurrent_metric_increments_are_counted_test.go:15 +0x1c

Previous write at 0x00c000091020 by goroutine 23:
  runtime.mapaccess2_faststr()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/internal/runtime/maps/runtime_faststr.go:161 +0x2ac
  keystone-asset-attestation/internal/ops.(*Metrics).Inc()
      /Users/tog/Documents/ChatGPT/新建go/2026-08-24/keystone-asset-attestation__017/.base_snapshot/internal/ops/metrics.go:20 +0xe0
  keystone-asset-attestation/internal/ops.TestBug017ConcurrentMetricIncrementsAreCounted.func1()
      /Users/tog/Documents/ChatGPT/新建go/2026-08-24/keystone-asset-attestation__017/.base_snapshot/internal/ops/bug017_concurrent_metric_increments_are_counted_test.go:15 +0x1c

Goroutine 34 (running) created at:
  keystone-asset-attestation/internal/ops.TestBug017ConcurrentMetricIncrementsAreCounted()
      /Users/tog/Documents/ChatGPT/新建go/2026-08-24/keystone-asset-attestation__017/.base_snapshot/internal/ops/bug017_concurrent_metric_increments_are_counted_test.go:15 +0xd8
  testing.tRunner()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2036 +0x164
  testing.(*T).Run.gowrap1()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2101 +0x34

Goroutine 23 (finished) created at:
  keystone-asset-attestation/internal/ops.TestBug017ConcurrentMetricIncrementsAreCounted()
      /Users/tog/Documents/ChatGPT/新建go/2026-08-24/keystone-asset-attestation__017/.base_snapshot/internal/ops/bug017_concurrent_metric_increments_are_counted_test.go:15 +0xd8
  testing.tRunner()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2036 +0x164
  testing.(*T).Run.gowrap1()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2101 +0x34
==================
    bug017_concurrent_metric_increments_are_counted_test.go:19: metric count missing: processed 17
    testing.go:1712: race detected during execution of test
--- FAIL: TestBug017ConcurrentMetricIncrementsAreCounted (0.00s)
=== RUN   TestBug017ConcurrentMetricIncrementsAreCounted
    bug017_concurrent_metric_increments_are_counted_test.go:19: metric count missing: processed 30
--- FAIL: TestBug017ConcurrentMetricIncrementsAreCounted (0.00s)
=== RUN   TestBug017ConcurrentMetricIncrementsAreCounted
    bug017_concurrent_metric_increments_are_counted_test.go:19: metric count missing: processed 27
--- FAIL: TestBug017ConcurrentMetricIncrementsAreCounted (0.00s)
=== RUN   TestBug017ConcurrentMetricIncrementsAreCounted
fatal error: concurrent map writes

goroutine 146 [running]:
internal/runtime/maps.fatal({0x100b2190c?, 0xffffffffffffffff?})
	/opt/homebrew/Cellar/go/1.26.5/libexec/src/runtime/panic.go:1181 +0x20
keystone-asset-attestation/internal/ops.(*Metrics).Inc(...)
	/Users/tog/Documents/ChatGPT/新建go/2026-08-24/keystone-asset-attestation__017/.base_snapshot/internal/ops/metrics.go:20
keystone-asset-attestation/internal/ops.TestBug017ConcurrentMetricIncrementsAreCounted.func1()
	/Users/tog/Documents/ChatGPT/新建go/2026-08-24/keystone-asset-attestation__017/.base_snapshot/internal/ops/bug017_concurrent_metric_increments_are_counted_test.go:15 +0xe4
created by keystone-asset-attestation/internal/ops.TestBug017ConcurrentMetricIncrementsAreCounted in goroutine 145
	/Users/tog/Documents/ChatGPT/新建go/2026-08-24/keystone-asset-attestation__017/.base_snapshot/internal/ops/bug017_concurrent_metric_increments_are_counted_test.go:15 +0xdc
FAIL	keystone-asset-attestation/internal/ops	0.509s
FAIL
```
