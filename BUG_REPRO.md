# BUG_REPRO

## Bug 是什么
worker 生产者 goroutine 已有 `defer close(ch)`，ctx.Done 分支又显式 `close(ch)` 后 return，导致 channel 被关闭两次，触发 panic。

## 如何触发
`go test ./internal/worker -run TestRunCancel`

## 错误信息
`panic: close of closed channel`
