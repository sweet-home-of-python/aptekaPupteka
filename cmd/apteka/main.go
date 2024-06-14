package main

import (
	"aptekaPupteka/internal/config"
	"aptekaPupteka/internal/http-server/handlers/drug"
	"aptekaPupteka/internal/storage/sqlite"
	"aptekaPupteka/lib/logger/sl"
	"encoding/base64"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

var users = map[string]string{
    "yar": "asslove",
    "dim": "ass",
}

func BasicAuth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        auth := r.Header.Get("Authorization")
        if auth == "" {
            w.WriteHeader(http.StatusUnauthorized)
            w.Write([]byte("Unauthorized"))
            return
        }

        parts := strings.SplitN(auth, " ", 2)
        if len(parts) != 2 || parts[0] != "Basic" {
            w.WriteHeader(http.StatusUnauthorized)
            w.Write([]byte("Unauthorized"))
            return
        }

        payload, err := base64.StdEncoding.DecodeString(parts[1])
        if err != nil {
            w.WriteHeader(http.StatusUnauthorized)
            w.Write([]byte("Unauthorized"))
            return
        }

        pair := strings.SplitN(string(payload), ":", 2)
        if len(pair) != 2 {
            w.WriteHeader(http.StatusUnauthorized)
            w.Write([]byte("Unauthorized"))
            return
        }

        user, pass := pair[0], pair[1]

        if password, ok := users[user]; !ok || password != pass {
            w.WriteHeader(http.StatusUnauthorized)
            w.Write([]byte("Unauthorized"))
            return
        }

        next.ServeHTTP(w, r)
    })
}

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env))
	log.Info("starting url", slog.String("adress", cfg.Address))
	log.Debug("debug message enabled")
	storage, err := sqlite.New(cfg.StoragePath)

	if err != nil {
		log.Error("failed to initialize storage:", sl.Err(err))
		os.Exit(1)
	}

	// authenticator := middleware.BasicAuth("realm", map[string]string{
	// 	"yarik": "30sm",
	// })

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	// router.Use(authenticator)
	router.Use(middleware.Logger)    // Логирование всех запросов
	router.Use(middleware.Recoverer) // Если где-то внутри сервера (обработчика запроса) произойдет паника, приложение не должно упасть
	router.Use(middleware.URLFormat) // Парсер URLов поступающих запросов
	router.With(BasicAuth).Post("/api/newDrug", drug.New(log, storage))
	router.With(BasicAuth).Post("/api/addDrug", drug.Add(log, storage))
	router.With(BasicAuth).Post("/api/subDrug", drug.Sub(log, storage))
	router.With(BasicAuth).Get("/api/getDrugs", drug.GetAll(log, storage))
	router.With(BasicAuth).Post("/api/deleteDrug", drug.Delete(log, storage))
	router.With(BasicAuth).Post("/api/getPage", drug.GetPage(log, storage))
	router.With(BasicAuth).Post("/api/searchDrug", drug.Search(log, storage))

	router.Handle("/*", http.FileServer(http.Dir("./cmd/apteka/static/")))
	router.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./cmd/apteka/static/index.html")
	}))

	log.Info("starting server", slog.String("adress", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}
	log.Error("server stopped!")

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
