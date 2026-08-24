## Bug 是什么
本题对应的业务异常是：资产配置页面读取模块列表后，运营人员会按页面需要调整展示内容；调用方修改了返回 slice，下一次创建记录时模块顺序被悄悄改变，评分和证据结果也跟着漂移。请修复模块列表的切片隔离，避免调用方改动内部状态。

## 如何触发
在含 Bug 的工作区执行下面的目标复现指令，观察接口结果或并发运行时行为是否与业务预期不一致。

## 根因
文件：internal/api/modules.go、internal/api/server.go、internal/core/engine.go；符号：(*Server).modules、(*Engine).Modules、(*Engine).Create；失效机制：模块列表接口直接暴露了引擎内部 slice 的共享底层数组，调用方修改返回值后会污染后续创建使用的模块顺序。

## 运行指令
```bash
go test -v ./internal/core -run '^TestBug020ModuleSliceDoesNotPolluteEngine$' -count=1
```

## 错误信息
目标行为在含 Bug 快照中失败，正常回归行为应保持通过；以下保留本题实际运行产生的原始输出片段。

## 错误堆栈
```text
=== RUN   TestBug020ModuleSliceDoesNotPolluteEngine
    bug020_module_slice_does_not_pollute_engine_test.go:10: module slice mutation polluted engine
--- FAIL: TestBug020ModuleSliceDoesNotPolluteEngine (0.00s)
FAIL
FAIL	keystone-asset-attestation/internal/core	0.479s
FAIL
```
