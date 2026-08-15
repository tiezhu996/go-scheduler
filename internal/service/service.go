package service

import (
	"errors"
	"fmt"

	"scheduler/internal/config"
	"scheduler/internal/model"
	"scheduler/internal/store"
)

type Service struct {
	store     *store.Store
	batchSize int
}

func New(s *store.Store, cfg *config.Config) *Service {
	b := cfg.BatchSize
	if b <= 0 {
		b = 1
	}
	return &Service{store: s, batchSize: b}
}

func (svc *Service) Submit(j *model.Job) error {
	if !model.ValidJob(j) {
		return errors.New("invalid job")
	}
	if err := svc.store.Create(j); err != nil {
		return fmt.Errorf("submit %s: %v", j.ID, err)
	}
	return nil
}

func (svc *Service) ListDue(now int64) []*model.Job {
	return svc.store.ListDue(now)
}

func (svc *Service) ListBatches() [][]*model.Job {
	js := svc.store.ListDue(1<<62 - 1)
	model.SortJobs(js)
	return model.BuildBatches(js, svc.batchSize)
}

func (svc *Service) MarkDone(id string) error {
	if err := svc.store.MarkDone(id); err != nil {
		return fmt.Errorf("mark done %s: %v", id, err)
	}
	return nil
}

func (svc *Service) MarkFailed(id string) error {
	if err := svc.store.MarkFailed(id); err != nil {
		return fmt.Errorf("mark failed %s: %v", id, err)
	}
	return nil
}
