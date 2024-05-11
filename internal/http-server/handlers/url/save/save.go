package save

import (
	"aptekaPupteka/internal/storage"
	"aptekaPupteka/lib/logger/sl"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
)

type Request struct {
	Drug string `json:"drug"`
}

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
	Count  string `json:"count"`
}

type DrugSaver interface {
	SaveDrug(drugToSave string) (int64, error)
}

func New(log *slog.Logger, drugSaver DrugSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		// Добавляем к текущму объекту логгера поля op и request_id
		// Они могут очень упростить нам жизнь в будущем
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// Создаем объект запроса и анмаршаллим в него запрос
		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			// Такую ошибку встретим, если получили запрос с пустым телом
			// Обработаем её отдельно
			log.Error("request body is empty")

			render.JSON(w, r, "empty request")
			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, "failed to decode request")
			return
		}

		// Лучше больше логов, чем меньше - лишнее мы легко сможем почистить,
		// при необходимости. А вот недостающую информацию мы уже не получим.
		log.Info("request body decoded", slog.Any("req", req))

		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid request", sl.Err(err))
			render.JSON(w, r, "invalid request")
			return
		}
		id, err := drugSaver.SaveDrug(req.Drug)
		if errors.Is(err, storage.ErrDrugExist) {
			log.Info("drug already exist", slog.String("drug", req.Drug))
			render.JSON(w, r, "drug already exist")
			return
		}
		if err != nil {
			log.Error("failed save drug", sl.Err(err))
			render.JSON(w, r, "HUIIIII")
			return
		}
		log.Info("drug added", slog.Int64("id", id))
		render.JSON(w, r, "tak to zaebis")
	}
}

// func New(log *slog.Logger, drugSaver DrugSaver) http.HandleFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		const op = "handlers.url.save.New"

// 		log = log.With(
// 			slog.String("op", op),
// 			slog.String("request_id", middleware.GetReqID(r.Context())),
// 		)
// 		var req Request

// 		err := render.DecodeJSON(r.Body, &req)
// 		if err != nil {
// 			log.Error("failed to decode request body", sl.Err(err))
// 			render.JSON(w, r, "suka")
// 			return
// 		}

// 		if err := validator.New().Struct(req); err != nil {
// 			validateErr := err.(validator.ValidationErrors)
// 			log.Error("invalid request", sl.Err(err))
// 			render.JSON(w, r, validateErr)
// 			return
// 		}
// 	}
// }
