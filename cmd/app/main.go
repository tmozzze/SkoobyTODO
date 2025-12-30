package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
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

	srv := &http.Server{
		Addr:    port,
		Handler: router,

		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}
	log.Info("server starting", "port", port)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("server failed to start", "err", err)
		os.Exit(1)
	}

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
