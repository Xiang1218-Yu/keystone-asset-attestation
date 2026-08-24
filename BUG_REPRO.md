## Bug 是什么
本题对应的业务异常是：任务消费期间管理员动态注册处理器，队列会出现并发读写，race 检查直接报错并可能让任务丢失。需要修复处理器注册和 worker 读取之间的同步，同时保留已注册任务的正常执行。

## 如何触发
在含 Bug 的工作区执行下面的目标复现指令，观察接口结果或并发运行时行为是否与业务预期不一致。

## 根因
文件：internal/ops/audit.go、internal/ops/metrics.go、internal/ops/queue.go；符号：(*Queue).Register、(*Queue).worker；失效机制：处理器注册与 worker 读取共享处理器表时缺少一致的并发同步，动态注册期间会出现竞态并影响任务分发。

## 运行指令
```bash
go test -v ./internal/ops -run '^TestBug019QueueRegistrationStaysRaceFree$' -race -count=20
```

## 错误信息
目标行为在含 Bug 快照中失败，正常回归行为应保持通过；以下保留本题实际运行产生的原始输出片段。

## 错误堆栈
```text
=== RUN   TestBug019QueueRegistrationStaysRaceFree
==================
WARNING: DATA RACE
Write at 0x00c000119020 by goroutine 11:
  runtime.mapaccess2_faststr()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/internal/runtime/maps/runtime_faststr.go:161 +0x2ac
  keystone-asset-attestation/internal/ops.(*Queue).Register()
      /Users/tog/Documents/ChatGPT/新建go/2026-08-24/keystone-asset-attestation__019/.base_snapshot/internal/ops/queue.go:42 +0x6c
  keystone-asset-attestation/internal/ops.TestBug019QueueRegistrationStaysRaceFree.func1()
      /Users/tog/Documents/ChatGPT/新建go/2026-08-24/keystone-asset-attestation__019/.base_snapshot/internal/ops/bug019_queue_registration_stays_race_free_test.go:17 +0x80
  keystone-asset-attestation/internal/ops.TestBug019QueueRegistrationStaysRaceFree.gowrap1()
      /Users/tog/Documents/ChatGPT/新建go/2026-08-24/keystone-asset-attestation__019/.base_snapshot/internal/ops/bug019_queue_registration_stays_race_free_test.go:19 +0x38

Previous read at 0x00c000119020 by goroutine 7:
  runtime.mapdelete_fast64()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/internal/runtime/maps/runtime_fast64.go:505 +0x8c
  keystone-asset-attestation/internal/ops.(*Queue).worker()
      /Users/tog/Documents/ChatGPT/新建go/2026-08-24/keystone-asset-attestation__019/.base_snapshot/internal/ops/queue.go:68 +0x154
  keystone-asset-attestation/internal/ops.NewQueue.gowrap1()
      /Users/tog/Documents/ChatGPT/新建go/2026-08-24/keystone-asset-attestation__019/.base_snapshot/internal/ops/queue.go:35 +0x2c

Goroutine 11 (running) created at:
  keystone-asset-attestation/internal/ops.TestBug019QueueRegistrationStaysRaceFree()
      /Users/tog/Documents/ChatGPT/新建go/2026-08-24/keystone-asset-attestation__019/.base_snapshot/internal/ops/bug019_queue_registration_stays_race_free_test.go:15 +0x74
  testing.tRunner()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2036 +0x164
  testing.(*T).Run.gowrap1()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2101 +0x34

Goroutine 7 (running) created at:
  keystone-asset-attestation/internal/ops.NewQueue()
      /Users/tog/Documents/ChatGPT/新建go/2026-08-24/keystone-asset-attestation__019/.base_snapshot/internal/ops/queue.go:35 +0x1bc
  keystone-asset-attestation/internal/ops.TestBug019QueueRegistrationStaysRaceFree()
      /Users/tog/Documents/ChatGPT/新建go/2026-08-24/keystone-asset-attestation__019/.base_snapshot/internal/ops/bug019_queue_registration_stays_race_free_test.go:11 +0x28
  testing.tRunner()
... output truncated for document readability ...
      /Users/tog/Documents/ChatGPT/新建go/2026-08-24/keystone-asset-attestation__019/.base_snapshot/internal/ops/queue.go:68 +0x15c
  keystone-asset-attestation/internal/ops.NewQueue.gowrap1()
      /Users/tog/Documents/ChatGPT/新建go/2026-08-24/keystone-asset-attestation__019/.base_snapshot/internal/ops/queue.go:35 +0x2c

Goroutine 11 (running) created at:
  keystone-asset-attestation/internal/ops.TestBug019QueueRegistrationStaysRaceFree()
      /Users/tog/Documents/ChatGPT/新建go/2026-08-24/keystone-asset-attestation__019/.base_snapshot/internal/ops/bug019_queue_registration_stays_race_free_test.go:15 +0x74
  testing.tRunner()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2036 +0x164
  testing.(*T).Run.gowrap1()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2101 +0x34

Goroutine 7 (running) created at:
  keystone-asset-attestation/internal/ops.NewQueue()
      /Users/tog/Documents/ChatGPT/新建go/2026-08-24/keystone-asset-attestation__019/.base_snapshot/internal/ops/queue.go:35 +0x1bc
  keystone-asset-attestation/internal/ops.TestBug019QueueRegistrationStaysRaceFree()
      /Users/tog/Documents/ChatGPT/新建go/2026-08-24/keystone-asset-attestation__019/.base_snapshot/internal/ops/bug019_queue_registration_stays_race_free_test.go:11 +0x28
  testing.tRunner()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2036 +0x164
  testing.(*T).Run.gowrap1()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2101 +0x34
==================
    testing.go:1712: race detected during execution of test
--- FAIL: TestBug019QueueRegistrationStaysRaceFree (0.01s)
=== RUN   TestBug019QueueRegistrationStaysRaceFree
--- PASS: TestBug019QueueRegistrationStaysRaceFree (0.01s)
=== RUN   TestBug019QueueRegistrationStaysRaceFree
--- PASS: TestBug019QueueRegistrationStaysRaceFree (0.01s)
=== RUN   TestBug019QueueRegistrationStaysRaceFree
--- PASS: TestBug019QueueRegistrationStaysRaceFree (0.01s)
=== RUN   TestBug019QueueRegistrationStaysRaceFree
--- PASS: TestBug019QueueRegistrationStaysRaceFree (0.01s)
=== RUN   TestBug019QueueRegistrationStaysRaceFree
--- PASS: TestBug019QueueRegistrationStaysRaceFree (0.01s)
=== RUN   TestBug019QueueRegistrationStaysRaceFree
--- PASS: TestBug019QueueRegistrationStaysRaceFree (0.01s)
=== RUN   TestBug019QueueRegistrationStaysRaceFree
--- PASS: TestBug019QueueRegistrationStaysRaceFree (0.01s)
=== RUN   TestBug019QueueRegistrationStaysRaceFree
--- PASS: TestBug019QueueRegistrationStaysRaceFree (0.01s)
=== RUN   TestBug019QueueRegistrationStaysRaceFree
--- PASS: TestBug019QueueRegistrationStaysRaceFree (0.02s)
=== RUN   TestBug019QueueRegistrationStaysRaceFree
--- PASS: TestBug019QueueRegistrationStaysRaceFree (0.01s)
=== RUN   TestBug019QueueRegistrationStaysRaceFree
--- PASS: TestBug019QueueRegistrationStaysRaceFree (0.01s)
=== RUN   TestBug019QueueRegistrationStaysRaceFree
--- PASS: TestBug019QueueRegistrationStaysRaceFree (0.01s)
=== RUN   TestBug019QueueRegistrationStaysRaceFree
--- PASS: TestBug019QueueRegistrationStaysRaceFree (0.01s)
=== RUN   TestBug019QueueRegistrationStaysRaceFree
--- PASS: TestBug019QueueRegistrationStaysRaceFree (0.01s)
=== RUN   TestBug019QueueRegistrationStaysRaceFree
--- PASS: TestBug019QueueRegistrationStaysRaceFree (0.01s)
=== RUN   TestBug019QueueRegistrationStaysRaceFree
--- PASS: TestBug019QueueRegistrationStaysRaceFree (0.01s)
=== RUN   TestBug019QueueRegistrationStaysRaceFree
--- PASS: TestBug019QueueRegistrationStaysRaceFree (0.01s)
=== RUN   TestBug019QueueRegistrationStaysRaceFree
--- PASS: TestBug019QueueRegistrationStaysRaceFree (0.01s)
=== RUN   TestBug019QueueRegistrationStaysRaceFree
--- PASS: TestBug019QueueRegistrationStaysRaceFree (0.01s)
FAIL
FAIL	keystone-asset-attestation/internal/ops	0.742s
FAIL
```
