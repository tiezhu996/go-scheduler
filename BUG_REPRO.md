# BUG_REPRO

## Bug 是什么
Store 构造时未初始化 jobs map，service.Submit 内部走 store.Create，首次向 nil map 写入触发 panic。

## 如何触发
`go test ./internal/store -run TestCreateGet`

## 错误信息
`panic: assignment to entry in nil map`
