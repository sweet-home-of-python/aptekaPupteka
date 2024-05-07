package main

import (
	"aptekaPupteka/internal/config"
	"aptekaPupteka/internal/storage/sqlite"
	"aptekaPupteka/lib/logger/sl"
	"log/slog"
	"os"
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
	// log.Info("saved drug", slog.Int64("id", id))
	// id, err := storage.TakeDrugCount("banan2", 22)
	// if err != nil {
	// 	log.Error("failed to add drugs", sl.Err(err))
	// 	os.Exit(1)
	// }
	id, err := storage.DeleteDrug("banan")
	if err != nil {
		log.Error("failed to delete drug", sl.Err(err))
		os.Exit(1)
	}
	log.Info("take drugs", slog.Int64("id", id))
	// id, err := storage.SaveDrug("banan")
	// if err != nil {
	// 	log.Error("failed to save url", sl.Err(err))
	// 	os.Exit(1)
	// }

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
