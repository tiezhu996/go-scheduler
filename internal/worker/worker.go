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

	var sum model.Summary

	for i := 0; i < p.workers; i++ {
		go func() {
			wg.Add(1)
			defer wg.Done()
			for batch := range ch {
				var local model.Summary
				for _, j := range batch {
					if err := p.runner.RunJob(ctx, j); err != nil {
						_ = p.svc.MarkFailed(j.ID)
						local.Failed++
						continue
					}
					_ = p.svc.MarkDone(j.ID)
					local.Ran++
				}
				sum = model.MergeSummary(sum, local)
			}
		}()
	}

	wg.Wait()
	return sum
}
