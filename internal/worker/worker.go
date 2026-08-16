package worker

import (
	"context"
	"sync"

	"scheduler/internal/model"
	"scheduler/internal/service"
)

type Runner interface {
	RunJob(ctx context.Context, j *model.Job) error
}

type Pool struct {
	svc     *service.Service
	runner  Runner
	workers int
}

func New(svc *service.Service, r Runner, workers int) *Pool {
	if workers <= 0 {
		workers = 1
	}
	return &Pool{svc: svc, runner: r, workers: workers}
}

func (p *Pool) Run(ctx context.Context) model.Summary {
	batches := p.svc.ListBatches()

	var wg sync.WaitGroup
	ch := make(chan []*model.Job, len(batches))

	go func() {
		defer close(ch)
		for _, b := range batches {
			select {
			case <-ctx.Done():
				return
			case ch <- b:
			}
		}
	}()

	var (
		mu  sync.Mutex
		sum model.Summary
	)

	for i := 0; i < p.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var local model.Summary
			for batch := range ch {
				for _, j := range batch {
					if err := ctx.Err(); err != nil {
						break
					}
					if err := p.runner.RunJob(ctx, j); err != nil {
						_ = p.svc.MarkFailed(j.ID)
						local.Failed++
						continue
					}
					_ = p.svc.MarkDone(j.ID)
					local.Ran++
				}
			}
			mu.Lock()
			sum = model.MergeSummary(sum, local)
			mu.Unlock()
		}()
	}

	wg.Wait()
	return sum
}
