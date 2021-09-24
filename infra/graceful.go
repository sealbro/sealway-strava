package infra

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Graceful struct {
	StartAction    func() error
	DeferAction    func(ctx context.Context) error
	ShutdownAction func(ctx context.Context) error
}

func (graceful *Graceful) RunAndWait() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := graceful.StartAction(); err != nil && err != http.ErrServerClosed {
			Log.Fatalf("listen: %s\n", err)
		}
	}()
	Log.Info("Server Started")

	<-done
	Log.Info("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		if err := graceful.DeferAction(ctx); err != nil {
			Log.Fatalf("Server  Failed:%+v", err)
		}
		cancel()
	}()

	if err := graceful.ShutdownAction(ctx); err != nil {
		Log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	Log.Info("Server Exited Properly")
}

func EnvOrDefault(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
