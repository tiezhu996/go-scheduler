package model

import "testing"

func TestValidJob(t *testing.T) {
	if !ValidJob(&Job{ID: "a", Name: "x"}) || ValidJob(nil) || ValidJob(&Job{}) {
		t.Fatal("ValidJob wrong")
	}
}

func TestSortJobs(t *testing.T) {
	in := []*Job{{ID: "c", RunAt: 3}, {ID: "a", RunAt: 1}, {ID: "b", RunAt: 2}}
	got := SortJobs(in)
	for i, id := range []string{"a", "b", "c"} {
		if got[i].ID != id {
			t.Fatalf("order=%v", got)
		}
	}
}

func TestBuildBatchesFresh(t *testing.T) {
	in := []*Job{{ID: "1"}, {ID: "2"}, {ID: "3"}}
	b := BuildBatches(in, 2)
	if len(b) != 2 {
		t.Fatalf("len=%d", len(b))
	}
	b[0][0] = &Job{ID: "x"}
	if in[0].ID != "1" {
		t.Fatal("mutating batch corrupted input")
	}
}

func TestMergeSummary(t *testing.T) {
	got := MergeSummary(Summary{Ran: 1}, Summary{Ran: 2, Failed: 3})
	if got.Ran != 3 || got.Failed != 3 {
		t.Fatalf("%+v", got)
	}
}
