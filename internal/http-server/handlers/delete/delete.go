package delete

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/KAA295/go-url-shortener/internal/storage"

	resp "github.com/KAA295/go-url-shortener/internal/lib/api/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type UrlDeleter interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlDeleter UrlDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		err := urlDeleter.DeleteURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", alias)

			render.JSON(w, r, resp.Error("not found"))

			return
		}

		if err != nil {
			log.Info("couldn't delete url")

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("deleted url", slog.String("alias", alias))

		render.JSON(w, r, resp.OK())
	}
}
