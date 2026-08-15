package config

import "testing"

func TestLoad(t *testing.T) {
	t.Setenv("SCHED_WORKERS", "")
	t.Setenv("SCHED_BATCH_SIZE", "")
	c := Load()
	if c.Workers != 2 || c.BatchSize != 2 {
		t.Fatalf("%+v", c)
	}
}
