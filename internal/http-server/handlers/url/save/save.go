package save

import (
	"aptekaPupteka/lib/logger/sl"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-playground/validator"
	"github.com/go-text/render"
)

type Request struct {
	Drug string `json:"drug" validate:"required,drug"`
}

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
	Count  string `json:"count"`
}

type DrugSaver interface {
	SaveDrug(drugToSave string) (int64, error)
}

func New(log *slog.Logger, drugSaver DrugSaver) http.HandleFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, "suka")
			return
		}

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			render.JSON(w, r, validateErr)
			return
		}
	}
}
