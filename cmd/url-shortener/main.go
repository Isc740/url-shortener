package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/Isc740/url-shortener/internal/app"
	"github.com/Isc740/url-shortener/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.NewConfig()

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to db: ", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatal("unable to ping database", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	a := app.NewApp(cfg, pool, logger)

	if err := a.Run(); err != nil {
		log.Fatal("server failed: ", err)
	}

	// mux := http.NewServeMux()
}
