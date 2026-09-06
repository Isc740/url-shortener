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

	userService := service.NewUserService(queries)
	userHandler := handler.NewUserHandler(userService, logger)

	linkService := service.NewLinkService(queries)
	linkHandler := handler.NewLinkHandler(linkService)

	a := &App{
		Config:  cfg,
		Queries: queries,
		Logger:  logger,
		Router:  http.NewServeMux(),
	}
	a.registerRoutes(userHandler, linkHandler)
	return a
}

func (a *App) registerRoutes(userHandler *handler.UserHandler, linkHandler *handler.LinkHandler) {
	a.Router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	a.Router.HandleFunc("GET /users", userHandler.GetAll)
	a.Router.HandleFunc("GET /users/{$}", userHandler.GetAll)
	a.Router.HandleFunc("GET /users/{id}", userHandler.GetByID)
	a.Router.HandleFunc("POST /users", userHandler.Create)

	a.Router.HandleFunc("POST /links", linkHandler.Create)
}

func (a *App) Run() error {
	addr := ":" + a.Config.Port
	a.Logger.Info("server starting", "addr", addr)
	return http.ListenAndServe(addr, a.Router)
}
