package main

import (
	"aptekaPupteka/internal/config"
	"aptekaPupteka/internal/http-server/handlers/drug"
	"aptekaPupteka/internal/http-server/middleware/basicAuth"
	"aptekaPupteka/internal/http-server/middleware/logger"
	"aptekaPupteka/internal/storage/sqlite"
	"aptekaPupteka/lib/logger/sl"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	cfg := config.MustLoad() // Загружаем конфиг
	log := logger.SetupLogger(cfg.Env) // Настраиваем логгер
	log = log.With(slog.String("env", cfg.Env)) // Добавляем информацию об окружении
	log.Info("starting url", slog.String("adress", cfg.Address)) 
	log.Debug("debug message enabled")
	storage, err := sqlite.New(cfg.StoragePath) // Инициализируем хранилище

	if err != nil {
		log.Error("failed to initialize storage:", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter() // Инициализируем роутер
	router.Use(middleware.RequestID) // Идентификатор запроса
	router.Use(middleware.Logger)    // Логирование всех запросов
	router.Use(middleware.Recoverer) // Если где-то внутри сервера (обработчика запроса) произойдет паника, приложение не должно упасть
	router.Use(middleware.URLFormat) // Парсер URLов поступающих запросов

	 // Маршрутизация

	router.With(basicAuth.BasicAuth).Post("/api/newDrug", drug.New(log, storage))
	router.With(basicAuth.BasicAuth).Post("/api/addDrug", drug.Add(log, storage)) 
	router.With(basicAuth.BasicAuth).Post("/api/subDrug", drug.Sub(log, storage))
	router.With(basicAuth.BasicAuth).Get("/api/getDrugs", drug.GetAll(log, storage))
	router.With(basicAuth.BasicAuth).Post("/api/deleteDrug", drug.Delete(log, storage))
	router.With(basicAuth.BasicAuth).Post("/api/getPage", drug.GetPage(log, storage))
	router.With(basicAuth.BasicAuth).Post("/api/searchDrug", drug.Search(log, storage))
	
	// Статические файлы

	router.Handle("/*", http.FileServer(http.Dir("./cmd/apteka/static/"))) 
	router.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./cmd/apteka/static/index.html")
	}))


	log.Info("starting server", slog.String("adress", cfg.Address))

	srv := &http.Server{ // Запускаем сервер
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
