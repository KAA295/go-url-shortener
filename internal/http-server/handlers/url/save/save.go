package save

import (
	"errors"
	"log/slog"
	"net/http"

	resp "github.com/KAA295/go-url-shortener/internal/lib/api/response"
	"github.com/KAA295/go-url-shortener/internal/lib/logger/sl"
	"github.com/KAA295/go-url-shortener/internal/lib/random"
	"github.com/KAA295/go-url-shortener/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct { //Структура запроса
	URL   string `json:"url" validate:"required,url"` //validate - поля для валидации - required - обязательное поле. url - должно являться url
	Alias string `json:"alias,omitempty"`
}

type Response struct { //Структура ответа
	resp.Response
	Alias string `json:"alias,omitempty"`
}

// TODO: move to config
const aliasLength = 4

// Mock иммитирует Storage
type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc { //Передаем сюда Storage, который автоматически будет удовлетворять интерфейсу
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With( //То есть логи пишем с названием операции и с реквест айди
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		//Анмаршалинг запроса
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request")) //Отправляем ошибку с помощью метода, описанного в lib

			return //Обязательно прерываем выполнение хэндлера
		}

		log.Info("request body decoded", slog.Any("request", req))

		//Валидация запроса
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors) //приводим ошибку к нужному типу

			log.Error("invalide request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		alias := req.Alias

		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URL))

			render.JSON(w, r, resp.Error("url already exists"))

			return
		}
		if err != nil {
			log.Error("failed to add url", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to add url"))

			return
		}

		log.Info("url added", slog.Int64("id", id))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Alias:    alias,
		})

	}
}
