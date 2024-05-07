package main

import (
	"aptekaPupteka/internal/config"
	"aptekaPupteka/internal/storage/sqlite"
	"aptekaPupteka/lib/logger/sl"
	"log/slog"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env))
	log.Info("starting url", slog.String("adress", cfg.Address))
	log.Debug("debug message enabled")
	storage, err := sqlite.New(cfg.StoragePath)
	_ = storage
	if err != nil {
		log.Error("failed to initialize storage", sl.Err(err))
		os.Exit(1)
	}
	// id, err := storage.SaveDrug("banan2")
	// if err != nil {
	// 	log.Error("failed to save drug", sl.Err(err))
	// 	os.Exit(1)
	// }
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)    // Логирование всех запросов
	router.Use(middleware.Recoverer) // Если где-то внутри сервера (обработчика запроса) произойдет паника, приложение не должно упасть
	router.Use(middleware.URLFormat) // Парсер URLов поступающих запросов

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
