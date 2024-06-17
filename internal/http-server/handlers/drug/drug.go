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

// еденица наркотика
type DrugUnit struct {
	Drug  string `json:"name"`
	Count int32  `json:"quantity"`
}

// единица страницы
type PageUnit struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// структура ответа
type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
	Count  string `json:"count"`
}

// структура ответа с наркотиком
type ResponseDrug struct {
	Name     string `json:"name"`
	Quantity int32  `json:"quantity"`
}

// интерфейс для создания нового наркотика
type DrugSaver interface {
	NewDrug(drugToSave string, count int32) (uint, error)
}

// интерфейс для добавления наркотика
type DrugAdder interface {
	AddDrug(name string, count int32) (uint, error)
}

// интерфейс для вычитания наркотика
type DrugSubber interface {
	SubDrug(name string, count int32) (uint, error)
}

// интерфейс для получения всех наркотиков
type DrugScrubber interface {
	GetAllDrugs() ([]sqlite.Drugs, error)
}

// интерфейс для удаления наркотика
type DrugDeleter interface {
	DeleteDrug(name string) (uint, error)
}

// интерфейс для получения страниц наркотиков
type PageGetter interface {
	GetPage(page int, limit int) ([]sqlite.Drugs, error)
}
// интерфейс для поиска наркотиков
type DrugSeacher interface {
	SearchDrug(drug string) ([]sqlite.Drugs, error)
}



// метод new для интерфейса создания нового наркотика
func New(log *slog.Logger, drugSaver DrugSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.drug.New"

		
		// Добавляем к текущму объекту логгера поля op и request_id
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req DrugUnit

		err := render.DecodeJSON(r.Body, &req)

		

		log.Info("request body decoded:", slog.Any("req", req))

		if err := validator.New().Struct(req); err != nil {
			log.Error("validation error:", sl.Err(err))
			render.JSON(w, r, "invalid request") //отвечаем дегенератам что запрос неверный
			return
		}

		id, err := drugSaver.NewDrug(req.Drug, req.Count)

		if err != nil {
			if errors.Is(err, storage.ErrDrugExist) {
				log.Info("drug already exist:", slog.String("drug", req.Drug))
				render.JSON(w, r, "drug already exist")
				return
			}else{
				log.Error("failed save drug", sl.Err(err))
				render.JSON(w, r, "save drug error")
				return
			}
		}
		
		log.Info("drug added", "id", id)
		render.JSON(w, r, "drug saved")

	}
}


// метод add для интерфейса добавления наркотика
func Add(log *slog.Logger, drugAdder DrugAdder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.drug.Add"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req DrugUnit

		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Error("request body is empty")
				render.JSON(w, r, "empty request")
				return
			}else{
				log.Error("failed to decode request body", sl.Err(err))
				render.JSON(w, r, "failed to decode request")
				return
			}
		}
		
		log.Info("request body decoded", slog.Any("req", req))

		
		if err := validator.New().Struct(req); err != nil {
			log.Error("validation error:", sl.Err(err))
			render.JSON(w, r, "invalid request")
			return
		}

		_ , err = drugAdder.AddDrug(req.Drug, req.Count)

		if err != nil {
			log.Error("failed add drug:", sl.Err(err))
			render.JSON(w, r, err.Error())
			return
		}

		log.Info("drug count added:", "count", req.Count)
		render.JSON(w, r, "drug count added")

	}
}

// метод sub для интерфейса вычитания наркотика
func Sub(log *slog.Logger, drugSubber DrugSubber) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.Sub"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)


		var req DrugUnit

		err := render.DecodeJSON(r.Body, &req)
		
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Error("request body is empty")
				render.JSON(w, r, "empty request")
				return
			}else{
				log.Error("failed to decode request body", sl.Err(err))
				render.JSON(w, r, "failed to decode request")
				return
			}
			
		}


		log.Info("request body decoded:", slog.Any("req", req))

		if err := validator.New().Struct(req); err != nil {
			log.Error("validation error:", sl.Err(err))
			render.JSON(w, r, "invalid request")
			return
		}

		_, err = drugSubber.SubDrug(req.Drug, req.Count)

		if err != nil {
			log.Error("failed to sub drug:", sl.Err(err))
			render.JSON(w, r, err.Error())
			return
		}

		log.Info("drug count substracted:", "count", req.Count)
		render.JSON(w, r, "drug count substracted")
	}
}
// метод get для интерфейса получения списка наркотиков

func GetAll(log *slog.Logger, scrubber DrugScrubber) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.GetAll"


		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		list, err := scrubber.GetAllDrugs()

		if err != nil {
			log.Error("failed colect drugs:", sl.Err(err))
			render.JSON(w, r, list)
			return
		}
		log.Info("drug collected:", slog.Any("list", list))
		render.JSON(w, r, "tak to zaebis no hui ti che poluchish")
	}
}

// метод delete для интерфейса удаления наркотика

func Delete(log *slog.Logger, drugDeleter DrugDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.delete"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req DrugUnit
		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("failed to decode request body:", sl.Err(err))
			render.JSON(w, r, "failed to decode request")
			return
		}

		id, err := drugDeleter.DeleteDrug(req.Drug)

		if err != nil {
			log.Error("failed to delete drug:", sl.Err(err))
			render.JSON(w, r, err.Error())
			return
		}
		log.Info("drug  deleted:", "id", id)
		render.JSON(w, r, "drug deleted")
	}
}

// метод get для интерфейса получения страницы наркотиков
func GetPage(log *slog.Logger, pageGetter PageGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.GetPage"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req PageUnit

		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("failed to decode request body:", sl.Err(err))
			render.JSON(w, r, "failed to decode request")
			return
		}

		list, err := pageGetter.GetPage(req.Page, req.Limit) // получаем страницу

		var resp []ResponseDrug

		if err != nil {
			log.Error("failed get page:", sl.Err(err))
			render.JSON(w, r, "failed get page")
			return
		}

		for _, v := range list {
			resp = append(resp, ResponseDrug{Name: v.Name, Quantity: v.Count}) // формируем ответ
		}

		log.Info("page collected:", slog.Any("list", list))
		render.JSON(w, r, resp)
	}
}

func Search(log *slog.Logger, seacher DrugSeacher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.drug.Search"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req DrugUnit

		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("failed to decode request body:", sl.Err(err))
			render.JSON(w, r, "failed to decode request")
			return
		}

		list, err := seacher.SearchDrug(req.Drug) // получаем страницу

		var resp []ResponseDrug

		if err != nil {
			log.Error("failed search drug:", sl.Err(err))
			render.JSON(w, r, "failed search drug")
			return
		}

		for _, v := range list {
			resp = append(resp, ResponseDrug{Name: v.Name, Quantity: v.Count}) // формируем ответ
		}

		log.Info("search drug collected:", slog.Any("list", list))
		render.JSON(w, r, resp)
	}
}