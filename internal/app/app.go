package app

import (
	"log/slog"
	"net/http"

	"github.com/Isc740/url-shortener/internal/config"
	"github.com/Isc740/url-shortener/internal/db"
	"github.com/Isc740/url-shortener/internal/handler"
	"github.com/Isc740/url-shortener/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Config  config.Config
	Queries *db.Queries
	Logger  *slog.Logger
	Router  *http.ServeMux
}

func NewApp(cfg config.Config, pool *pgxpool.Pool, logger *slog.Logger) *App {
	queries := db.New(pool)

	userHandler := handler.NewUserHandler(service.NewUserService(queries), logger)
	linkHandler := handler.NewLinkHandler(service.NewLinkService(queries))

	app := &App{
		Config:  cfg,
		Queries: queries,
		Logger:  logger,
		Router:  http.NewServeMux(),
	}

	userHandler.RegisterRoutes(app.Router)
	linkHandler.RegisterRoutes(app.Router)

	app.Router.HandleFunc("GET /", healthCheck)

	return app
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (a *App) Run() error {
	addr := ":" + a.Config.Port
	a.Logger.Info("server starting", "addr", addr)
	return http.ListenAndServe(addr, a.Router)
}
