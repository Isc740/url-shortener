package app

import (
	"log/slog"
	"net/http"

	"github.com/Isc740/url-shortener/internal/config"
	"github.com/Isc740/url-shortener/internal/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Config  config.Config
	DB      *pgxpool.Pool
	Queries *db.Queries
	Logger  *slog.Logger
	Router  *http.ServeMux
}

func NewApp(cfg config.Config, pool *pgxpool.Pool, logger *slog.Logger) *App {
	a := &App{
		Config:  cfg,
		DB:      pool,
		Queries: db.New(pool),
		Logger:  logger,
		Router:  http.NewServeMux(),
	}
	a.registerRoutes()
	return a
}

func (a *App) registerRoutes() {
	a.Router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
}

func (a *App) Run() error {
	addr := ":" + a.Config.Port
	a.Logger.Info("server starting", "addr", addr)
	return http.ListenAndServe(addr, a.Router)
}
