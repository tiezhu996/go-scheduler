# BUG_REPRO

## Bug 是什么
BuildBatches / OrderIDs / ListBatches 返回共享底层数组子切片，worker 又用 batch[:len(batch)-1] 切掉最后一条，导致任务漏掉、列表串改。

## 如何触发
`go test ./...`

## 错误信息
- TestBuildBatchesFresh / TestOrderIDsFresh 失败。
- TestRunSummary 失败（Ran 计数不对）。
