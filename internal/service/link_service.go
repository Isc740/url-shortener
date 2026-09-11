package service

import (
	"context"
	"fmt"
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

func (s *LinkService) GetAll(ctx context.Context) ([]db.GetLinksRow, error) {
	return s.queries.GetLinks(ctx)
}

func (s *LinkService) GetByID(ctx context.Context, id int) (db.GetLinkByIDRow, error) {
	return s.queries.GetLinkByID(ctx, int64(id))
}

func (s *LinkService) GetByTargetURL(ctx context.Context, targetUrl string) (db.GetLinkByTargetURLRow, error) {
	return s.queries.GetLinkByTargetURL(ctx, targetUrl)
}

func (s *LinkService) GetByShortenedURL(ctx context.Context, shortenedURL string) (db.GetLinkByShortenedURLRow, error) {
	return s.queries.GetLinkByShortenedURL(ctx, pgtype.Text{String: shortenedURL})
}

func (s *LinkService) GetByStatus(ctx context.Context, status string) ([]db.GetLinksByStatusRow, error) {
	return s.queries.GetLinksByStatus(ctx, status)
}

func (s *LinkService) Create(ctx context.Context, params CreateLinkDTO) (LinkDTO, error) {
	link, err := s.queries.CreateLink(ctx, db.CreateLinkParams{
		UserID:         pgtype.Int8{Int64: int64(params.UserID)},
		TargetUrl:      params.TargetURL,
		Password:       pgtype.Text{String: params.Password, Valid: true},
		Status:         "active",
		ExpirationDate: pgtype.Timestamptz{Time: params.ExpirationDate},
		CreatedAt:      pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt:      pgtype.Timestamptz{Time: time.Now(), Valid: true},
	})
	if err != nil {
		return LinkDTO{}, err
	}

	shortenedURL := base62.Encode(uint64(link.ID))
	updatedLink, err := s.queries.UpdateLinkShortenedURL(ctx, db.UpdateLinkShortenedURLParams{
		ID:           link.ID,
		ShortenedUrl: pgtype.Text{String: shortenedURL, Valid: true},
	})
	if err != nil {
		return LinkDTO{}, err
	}

	return LinkDTO{
		ID:             updatedLink.ID,
		UserName:       updatedLink.UserName,
		TargetUrl:      updatedLink.TargetUrl,
		ShortenedUrl:   updatedLink.ShortenedUrl.String,
		Status:         updatedLink.Status,
		ExpirationDate: updatedLink.ExpirationDate.Time,
		CreatedAt:      updatedLink.CreatedAt.Time,
		UpdatedAt:      updatedLink.UpdatedAt.Time,
	}, nil
}

func (s *LinkService) Update(ctx context.Context, params UpdateLinkDTO) (LinkDTO, error) {
	link, err := s.queries.UpdateLink(ctx, db.UpdateLinkParams{
		ID:             int64(params.ID),
		TargetUrl:      params.TargetURL,
		Status:         params.Status,
		ExpirationDate: pgtype.Timestamptz{Time: params.ExpirationDate, Valid: true},
		UpdatedAt:      pgtype.Timestamptz{Time: time.Now(), Valid: true},
	})
	if err != nil {
		return LinkDTO{}, fmt.Errorf("error updating link: %w", err)
	}

	return LinkDTO{
		ID:             link.ID,
		UserName:       link.UserName,
		TargetUrl:      link.TargetUrl,
		ShortenedUrl:   link.ShortenedUrl.String,
		Status:         link.Status,
		ExpirationDate: link.ExpirationDate.Time,
		CreatedAt:      link.CreatedAt.Time,
		UpdatedAt:      link.UpdatedAt.Time,
	}, nil
}
