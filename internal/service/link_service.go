package service

import (
	"context"
	"time"

	"github.com/Isc740/url-shortener/internal/base62"
	"github.com/Isc740/url-shortener/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type LinkService struct {
	queries db.Querier
}

func NewLinkService(querier db.Querier) *LinkService {
	return &LinkService{
		queries: querier,
	}
}

func (s *LinkService) Create(ctx context.Context, params CreateLinkDTO) (LinkDTO, error) {
	id, err := s.queries.NextLinkID(ctx)
	if err != nil {
		return LinkDTO{}, err
	}

	link, err := s.queries.CreateLink(ctx, db.CreateLinkParams{
		UserID:         pgtype.Int8{Int64: int64(params.UserID), Valid: true},
		TargetUrl:      params.TargetURL,
		ShortenedUrl:   base62.Encode(uint64(id)),
		Password:       pgtype.Text{String: params.Password, Valid: true},
		Status:         "active",
		ExpirationDate: pgtype.Timestamptz{Time: params.ExpirationDate, Valid: true},
		CreatedAt:      pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt:      pgtype.Timestamptz{Time: time.Now(), Valid: true},
	})
	if err != nil {
		return LinkDTO{}, err
	}

	return LinkDTO{
		ID:             link.ID,
		TargetUrl:      link.TargetUrl,
		ShortenedUrl:   link.ShortenedUrl,
		Status:         link.Status,
		ExpirationDate: link.ExpirationDate.Time,
		CreatedAt:      link.CreatedAt.Time,
		UpdatedAt:      link.UpdatedAt.Time,
	}, nil
}
