package store

import (
	"testing"

	"scheduler/internal/model"
)

func TestCreateGet(t *testing.T) {
	s := New()
	if err := s.Create(&model.Job{ID: "a", Name: "x", Status: model.StatusPending}); err != nil {
		t.Fatal(err)
	}
	if err := s.Create(&model.Job{ID: "a"}); err != ErrJobExists {
		t.Fatalf("dup err=%v", err)
	}
	if _, err := s.Get("a"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Get("nope"); err != ErrJobNotFound {
		t.Fatalf("err=%v", err)
	}
}

func TestListDue(t *testing.T) {
	s := New()
	_ = s.Create(&model.Job{ID: "a", RunAt: 10, Status: model.StatusPending})
	_ = s.Create(&model.Job{ID: "b", RunAt: 20, Status: model.StatusPending})
	_ = s.Create(&model.Job{ID: "c", RunAt: 30, Status: model.StatusDone})
	got := s.ListDue(15)
	if len(got) != 1 || got[0].ID != "a" {
		t.Fatalf("ListDue=%v", got)
	}
}

func TestOrderIDsFresh(t *testing.T) {
	s := New()
	_ = s.Create(&model.Job{ID: "b", Status: model.StatusPending})
	_ = s.Create(&model.Job{ID: "a", Status: model.StatusPending})
	ids := s.OrderIDs()
	ids[0] = "z"
	if s.OrderIDs()[0] != "b" {
		t.Fatal("OrderIDs aliased")
	}
}

func TestMarkDoneFailed(t *testing.T) {
	s := New()
	_ = s.Create(&model.Job{ID: "a", Status: model.StatusPending})
	if err := s.MarkDone("a"); err != nil {
		t.Fatal(err)
	}
	if err := s.MarkFailed("nope"); err != ErrJobNotFound {
		t.Fatalf("err=%v", err)
	}
}
