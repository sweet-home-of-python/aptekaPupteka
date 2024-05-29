package drug

import (
	"aptekaPupteka/internal/storage"
	"aptekaPupteka/internal/storage/sqlite"
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
	Count int32 `json:"count"`
}

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
	Count  string `json:"count"`
}

type DrugSaver interface {
	NewDrug(drugToSave string) (uint, error)
}
type DrugAdder interface {
	AddDrug(name string, count int32) (uint, error)
}
type DrugSubber interface {
	SubDrug(name string, count int32) (uint, error)
}
type DrugScrubber interface {
	GetAllDrugs() ([]sqlite.Drugs, error)
}
type DrugDeleter interface {
	DeleteDrug(name string) (uint, error)
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
		id, err := drugSaver.NewDrug(req.Drug)
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
		log.Info("drug added", "id", id)
		render.JSON(w, r, "tak to zaebis")
	}
}

func Add(log *slog.Logger, drugAdder DrugAdder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.Add"

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
		id, err := drugAdder.AddDrug(req.Drug, req.Count)
		_ = id
		if errors.Is(err, storage.ErrDrugExist) {
			log.Info("tam pisda s dobavkoi", slog.String("drug", req.Drug))
			render.JSON(w, r, "tam pisda s dobavkoi")
			return
		}
		if err != nil {
			log.Error("failed add drug", sl.Err(err))
			render.JSON(w, r, err.Error())
			return
		}
		log.Info("drug count added", "count", req.Count)
		render.JSON(w, r, "tak to zaebis dobavil")
	}
}

func Sub(log *slog.Logger, drugSubber DrugSubber) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.Sub"

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
		id, err := drugSubber.SubDrug(req.Drug, req.Count)
		_ = id
		if errors.Is(err, storage.ErrDrugExist) {
			log.Info("tam pisda s grabezhom", slog.String("drug", req.Drug))
			render.JSON(w, r, "tam pisda s grabezhom")
			return
		}
		if err != nil {
			log.Error("failed sub drug", sl.Err(err))
			render.JSON(w, r, err.Error())
			return
		}
		log.Info("drug count spisdil", "count", req.Count)
		render.JSON(w, r, "tak to zaebis pizdanul")
	}
}

func GetAll(log *slog.Logger, scrubber DrugScrubber) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.GetAll"

		// Добавляем к текущму объекту логгера поля op и request_id
		// Они могут очень упростить нам жизнь в будущем
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		list, err := scrubber.GetAllDrugs()
		if err != nil {
			log.Error("failed colect drugs", sl.Err(err))
			render.JSON(w, r, list)
			return
		}
		log.Info("drug posipalis")
		render.JSON(w, r, "tak to zaebis nagrebli")
	}
}
func Delete(log *slog.Logger, drugDeleter DrugDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.delete"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil{

			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, "failed to decode request")
			return
		}
		id, err := drugDeleter.DeleteDrug(req.Drug)

		if err  != nil{
			log.Error("failed to delete drug", sl.Err(err))
			render.JSON(w, r, err.Error())
			return
		}
		log.Info("drug  petuh", "id", id)
		render.JSON(w, r, "tak to zaebis udalil")
	}
}