package store

import (
	"errors"
	"sync"

	"scheduler/internal/model"
)

var (
	ErrJobNotFound = errors.New("job not found")
	ErrJobExists   = errors.New("job already exists")
)

type Store struct {
	mu    sync.RWMutex
	jobs  map[string]*model.Job
	order []string
}

func New() *Store {
	return &Store{
		jobs:  make(map[string]*model.Job),
		order: []string{},
	}
}

func (s *Store) Create(j *model.Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.jobs[j.ID]; ok {
		return ErrJobExists
	}
	s.jobs[j.ID] = j
	s.order = append(s.order, j.ID)
	return nil
}

func (s *Store) Get(id string) (*model.Job, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	j, ok := s.jobs[id]
	if !ok {
		return nil, ErrJobNotFound
	}
	return j, nil
}

func (s *Store) ListDue(now int64) []*model.Job {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*model.Job, 0)
	for _, id := range s.order {
		j := s.jobs[id]
		if j.Status == model.StatusPending && j.RunAt <= now {
			out = append(out, j)
		}
	}
	return out
}

func (s *Store) OrderIDs() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]string, len(s.order))
	copy(out, s.order)
	return out
}

func (s *Store) MarkDone(id string) error {
	j, ok := s.jobs[id]
	if !ok {
		return ErrJobNotFound
	}
	j.Status = model.StatusDone
	return nil
}

func (s *Store) MarkFailed(id string) error {
	j, ok := s.jobs[id]
	if !ok {
		return ErrJobNotFound
	}
	j.Attempts++
	j.Status = model.StatusFailed
	return nil
}
