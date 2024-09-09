package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/KAA295/go-url-shortener/internal/config"
	"github.com/KAA295/go-url-shortener/internal/lib/logger/handlers/slogpretty"
	"github.com/KAA295/go-url-shortener/internal/lib/logger/sl"
	"github.com/KAA295/go-url-shortener/internal/storage/sqlite"

	del "github.com/KAA295/go-url-shortener/internal/http-server/handlers/delete"
	"github.com/KAA295/go-url-shortener/internal/http-server/handlers/redirect"
	"github.com/KAA295/go-url-shortener/internal/http-server/handlers/url/save"
	mwLogger "github.com/KAA295/go-url-shortener/internal/http-server/middleware/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	log.Info("starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID) //Встроенный мидлвейр в chi. Добавляет requestID к каждому поступающему запросу.
	router.Use(mwLogger.New(log))    //Будет логировать все входящие запросы
	router.Use(middleware.Recoverer) //Восстановление после паники
	router.Use(middleware.URLFormat) //Парсинг url

	router.Route("/url", func(r chi.Router) { //Авторизация
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.HttpServer.User: cfg.HttpServer.Password,
		}))

		r.Post("/", save.New(log, storage))
		r.Delete("/{alias}", del.New(log, storage))
	})

	router.Get("/{alias}", redirect.New(log, storage))

	log.Info("starting server", slog.String("adress", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	//В локальном окружении, логи выводятся текстом, в dev окружении - json, в prod - json и не выводятся сообщения дебага
	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log

}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
