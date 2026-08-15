package main

import (
	"fmt"

	"scheduler/internal/config"
	"scheduler/internal/service"
	"scheduler/internal/store"
)

func main() {
	cfg := config.Load()
	st := store.New()
	svc := service.New(st, cfg)
	_ = svc
	fmt.Println("scheduler ready")
}
