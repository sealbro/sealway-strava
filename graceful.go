package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Graceful struct {
	startAction    func() error
	deferAction    func(ctx context.Context) error
	shutdownAction func(ctx context.Context) error
}

func (graceful *Graceful) RunAndWait() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := graceful.startAction(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Info("Server Started")

	<-done
	log.Info("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		if err := graceful.deferAction(ctx); err != nil {
			log.Fatalf("Server  Failed:%+v", err)
		}
		cancel()
	}()

	if err := graceful.shutdownAction(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Info("Server Exited Properly")
}
