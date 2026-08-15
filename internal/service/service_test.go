package service

import (
	"errors"
	"testing"

	"scheduler/internal/config"
	"scheduler/internal/model"
	"scheduler/internal/store"
)

func newSvc() (*store.Store, *Service) {
	s := store.New()
	return s, New(s, config.Load())
}

func TestSubmitListDue(t *testing.T) {
	_, svc := newSvc()
	if err := svc.Submit(&model.Job{ID: "a", Name: "x", RunAt: 10, Status: model.StatusPending}); err != nil {
		t.Fatal(err)
	}
	if err := svc.Submit(&model.Job{}); err == nil {
		t.Fatal("invalid job should error")
	}
	if len(svc.ListDue(20)) != 1 {
		t.Fatal("ListDue len")
	}
}

func TestSubmitWraps(t *testing.T) {
	_, svc := newSvc()
	_ = svc.Submit(&model.Job{ID: "a", Name: "x", Status: model.StatusPending})
	if err := svc.Submit(&model.Job{ID: "a", Name: "x", Status: model.StatusPending}); !errors.Is(err, store.ErrJobExists) {
		t.Fatalf("errors.Is=false err=%v", err)
	}
}

func TestListBatchesOrder(t *testing.T) {
	_, svc := newSvc()
	_ = svc.Submit(&model.Job{ID: "b", Name: "x", RunAt: 20, Status: model.StatusPending})
	_ = svc.Submit(&model.Job{ID: "a", Name: "x", RunAt: 10, Status: model.StatusPending})
	batches := svc.ListBatches()
	if len(batches) == 0 || batches[0][0].ID != "a" {
		t.Fatalf("batches=%v", batches)
	}
}
