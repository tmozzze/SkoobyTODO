package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tmozzze/SkoobyTODO/internal/config"
	"github.com/tmozzze/SkoobyTODO/internal/handlers"
	"github.com/tmozzze/SkoobyTODO/internal/service"
	"github.com/tmozzze/SkoobyTODO/internal/storage/inmemory"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// Load config
	cfg := config.New()
	err := cfg.Load(".env")
	if err != nil {
		fmt.Fprintf(os.Stdout, "Config error: %v", err)
		os.Exit(1)
	}

	// Init logger
	log := setupLogger(cfg.Env)
	log.Info("starting SkoobyTODO", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// Init storage
	storage := inmemory.NewMemStorage(log)
	log.Info("storage is initialized")

	// Init service
	svc := service.NewService(storage, log)
	log.Info("service is initialized")

	// Init handlers
	handler := handlers.NewHandler(svc, log)
	log.Info("handler is initialized")

	// Init router
	router := handler.InitRoutes()

	// Server
	port := ":" + cfg.ServerPort
	readTimeout := time.Duration(cfg.ReadTimeout) * time.Second
	writeTimeout := time.Duration(cfg.WriteTimeout) * time.Second
	idleTimeout := time.Duration(cfg.IdleTimeout) * time.Second
	handlerTimeout := time.Duration(cfg.HandlerTimeout) * time.Second
	routerWithTimeout := http.TimeoutHandler(router, handlerTimeout, "Timeout!\n")

	srv := &http.Server{
		Addr:         port,
		Handler:      routerWithTimeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}
	log.Info("server starting", "port", port)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server failed to start", "err", err)
			os.Exit(1)
		}

	}()

	// graceful shutdown

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit // wait for signal

	log.Info("server is shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server forced to shutdown", "err", err)
	}

	log.Info("server exited properly")

}

func setupLogger(env string) *slog.Logger {
	switch env {
	case envLocal: // Text Debug
		return slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev: // JSON Debug
		return slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd: // JSON Info
		return slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}
