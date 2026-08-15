package model

import "sort"

type Job struct {
	ID       string
	Name     string
	RunAt    int64
	Status   string
	Attempts int
}

const (
	StatusPending = "pending"
	StatusRunning = "running"
	StatusDone    = "done"
	StatusFailed  = "failed"
)

type Summary struct {
	Ran    int
	Failed int
	Skipped int
}

func ValidJob(j *Job) bool {
	return j != nil && j.ID != "" && j.Name != ""
}

func SortJobs(js []*Job) []*Job {
	sort.SliceStable(js, func(i, j int) bool { return js[i].RunAt < js[j].RunAt })
	return js
}

func BuildBatches(js []*Job, size int) [][]*Job {
	if size <= 0 {
		size = 1
	}
	out := make([][]*Job, 0, (len(js)+size-1)/size)
	for i := 0; i < len(js); i += size {
		end := i + size
		if end > len(js) {
			end = len(js)
		}
		b := make([]*Job, end-i)
		copy(b, js[i:end])
		out = append(out, b)
	}
	return out
}

func MergeSummary(dst Summary, src Summary) Summary {
	dst.Ran += src.Ran
	dst.Skipped += src.Skipped
	return dst
}
