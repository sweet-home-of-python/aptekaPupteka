package logger

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"log/slog"

	"github.com/go-chi/chi/middleware"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log = log.With(
			slog.String("component", "middleware/logger"),
		)

		log.Info("logger middleware enabled")

		// код самого обработчика
		fn := func(w http.ResponseWriter, r *http.Request) {
			// собираем исходную информацию о запросе
			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)

			// создаем обертку вокруг `http.ResponseWriter`
			// для получения сведений об ответе
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// Момент получения запроса, чтобы вычислить время обработки
			t1 := time.Now()

			// Запись отправится в лог в defer
			// в этот момент запрос уже будет обработан
			defer func() {
				entry.Info("request completed",
					slog.Int("status", ww.Status()),
					slog.Int("bytes", ww.BytesWritten()),
					slog.String("duration", time.Since(t1).String()),
				)
			}()

			// Передаем управление следующему обработчику в цепочке middleware
			next.ServeHTTP(ww, r)
		}

		// Возвращаем созданный выше обработчик, приведя его к типу http.HandlerFunc
		return http.HandlerFunc(fn)
	}
}

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger
	file, err := os.OpenFile("./logs/logs.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644) // |os.O_APPEND
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return nil
	}
	mw := io.MultiWriter(os.Stdout, file)
	switch env {
	case EnvLocal:
		log = slog.New(slog.NewTextHandler(mw, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case EnvDev:
		log = slog.New(slog.NewJSONHandler(mw, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case EnvProd:
		log = slog.New(slog.NewJSONHandler(mw, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
