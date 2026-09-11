package service

import (
	"errors"
	"fmt"

	"github.com/Isc740/url-shortener/internal/base62"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func shortenURL(id int64) string {
	return base62.Encode(uint64(id))
}

func wrapNotFound(err error, notFound error, msg string) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return notFound
	}
	return fmt.Errorf("%s: %w", msg, err)
}

func wrapCreateError[T any](err error, context string, codes map[string]error) (T, error) {
	var data T
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if errorType, ok := codes[pgErr.Code]; ok {
				return data, errorType
			}
		}
	}
	return data, fmt.Errorf("error %s: %w", context, err)
}
