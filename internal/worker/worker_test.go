package worker

import (
	"context"
	"fmt"
	"testing"

	"scheduler/internal/config"
	"scheduler/internal/model"
	"scheduler/internal/service"
	"scheduler/internal/store"
)

type okRunner struct{}

func (okRunner) RunJob(ctx context.Context, j *model.Job) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return nil
}

func newPool() (*service.Service, *Pool) {
	s := store.New()
	svc := service.New(s, config.Load())
	return svc, New(svc, okRunner{}, 4)
}

func TestRunSummary(t *testing.T) {
	svc, p := newPool()
	for i := 0; i < 10; i++ {
		_ = svc.Submit(&model.Job{ID: fmt.Sprintf("j%d", i), Name: "x", RunAt: 0, Status: model.StatusPending})
	}
	sum := p.Run(context.Background())
	if sum.Ran != 10 {
		t.Fatalf("Ran=%d want 10", sum.Ran)
	}
}

type failRunner struct{ failID string }

func (f failRunner) RunJob(ctx context.Context, j *model.Job) error {
	if j.ID == f.failID {
		return fmt.Errorf("fail %s", j.ID)
	}
	return nil
}

func TestRunFailed(t *testing.T) {
	svc, _ := newPool()
	_ = svc.Submit(&model.Job{ID: "ok", Name: "x", RunAt: 0, Status: model.StatusPending})
	_ = svc.Submit(&model.Job{ID: "bad", Name: "x", RunAt: 0, Status: model.StatusPending})
	p := New(svc, failRunner{"bad"}, 2)
	sum := p.Run(context.Background())
	if sum.Failed != 1 || sum.Ran != 1 {
		t.Fatalf("sum=%+v", sum)
	}
}

func TestRunCancel(t *testing.T) {
	svc, p := newPool()
	_ = svc.Submit(&model.Job{ID: "j0", Name: "x", RunAt: 0, Status: model.StatusPending})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	sum := p.Run(ctx)
	if sum.Ran != 0 {
		t.Fatalf("Ran=%d want 0", sum.Ran)
	}
}
