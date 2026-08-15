# scheduler

一个用 Go 写的内存定时任务调度器，演示分层、并发派发、到点任务筛选与上下文取消。

## 功能
- 提交/查询任务、按运行时间筛选到点任务
- 任务分页查询
- 并发派发 worker 池，支持 context 取消

## 目录结构
```
cmd/scheduler/      程序入口
internal/config/    环境配置
internal/model/     模型与纯工具函数
internal/store/     内存存储（任务 + 锁）
internal/service/   业务逻辑
internal/worker/    派发 worker 池
```

## 运行与测试
```bash
go build ./...
go test ./...
go run ./cmd/scheduler
```

## 环境变量
| 变量 | 说明 | 默认值 |
|------|------|--------|
| `SCHED_WORKERS` | worker 数量 | `2` |
| `SCHED_BATCH_SIZE` | 分页大小 | `2` |
