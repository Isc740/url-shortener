package service

import (
	"context"
	"errors"
	"fmt"
	"time"

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

var (
	ErrLinkNotFound       = errors.New("link not found")
	ErrLinksNotFound      = errors.New("links not found")
	ErrLinkExpired        = errors.New("link has expired")
	ErrLinkInactive       = errors.New("link is inactive")
	ErrDuplicateTargetURL = errors.New("target url already shortened")
)

var linkCreateErrors = map[string]error{
	"23505": ErrDuplicateTargetURL,
	"23503": ErrUserNotFound,
}

func (s *LinkService) GetAll(ctx context.Context) ([]db.GetLinksRow, error) {
	links, err := s.queries.GetLinks(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting links: %w", err)
	}
	return links, err
}

func (s *LinkService) GetByStatus(ctx context.Context, status string) ([]db.GetLinksByStatusRow, error) {
	links, err := s.queries.GetLinksByStatus(ctx, status)
	if err != nil {
		return nil, fmt.Errorf("error getting target URL: %w", err)
	}
	return links, nil
}

func (s *LinkService) GetByID(ctx context.Context, id int) (db.GetLinkByIDRow, error) {
	link, err := s.queries.GetLinkByID(ctx, int64(id))
	if err != nil {
		return db.GetLinkByIDRow{}, wrapNotFound(err, ErrLinkNotFound, "error getting link")
	}
	return link, nil
}

func (s *LinkService) GetByTargetURL(ctx context.Context, targetUrl string) (db.GetLinkByTargetURLRow, error) {
	link, err := s.queries.GetLinkByTargetURL(ctx, targetUrl)
	if err != nil {
		return db.GetLinkByTargetURLRow{}, wrapNotFound(err, ErrLinkNotFound, "error getting link")
	}

	return link, nil
}

func (s *LinkService) GetByShortenedURL(ctx context.Context, shortenedURL string) (db.GetLinkByShortenedURLRow, error) {
	link, err := s.queries.GetLinkByShortenedURL(ctx, pgtype.Text{String: shortenedURL})
	if err != nil {
		return db.GetLinkByShortenedURLRow{}, wrapNotFound(err, ErrLinkNotFound, "error getting link")
	}

	return link, nil
}

func (s *LinkService) GetTargetURLByShortenedURL(ctx context.Context, shortenedURL string) (string, error) {
	link, err := s.queries.GetLinkByShortenedURL(ctx, pgtype.Text{String: shortenedURL, Valid: true})
	if err != nil {
		return "", wrapNotFound(err, ErrLinkNotFound, "error getting link")
	}
	if link.Status != "active" {
		return "", ErrLinkInactive
	}
	if link.ExpirationDate.Valid && time.Now().After(link.ExpirationDate.Time) {
		return "", ErrLinkExpired
	}

	return link.TargetUrl, nil
}

func (s *LinkService) Create(ctx context.Context, params CreateLinkDTO) (LinkDTO, error) {
	link, err := s.queries.CreateLink(ctx, db.CreateLinkParams{
		UserID:         pgtype.Int8{Int64: int64(params.UserID), Valid: true},
		TargetUrl:      params.TargetURL,
		Password:       pgtype.Text{String: params.Password, Valid: params.Password != ""},
		Status:         "active",
		ExpirationDate: pgtype.Timestamptz{Time: params.ExpirationDate, Valid: true},
		CreatedAt:      pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt:      pgtype.Timestamptz{Time: time.Now(), Valid: true},
	})

	if err != nil {
		return wrapCreateError[LinkDTO](err, "creating link", linkCreateErrors)
	}

	shortenedURL := shortenURL(link.ID)
	updatedLink, err := s.queries.UpdateLinkShortenedURL(ctx, db.UpdateLinkShortenedURLParams{
		ID:           link.ID,
		ShortenedUrl: pgtype.Text{String: shortenedURL, Valid: true},
	})
	if err != nil {
		return LinkDTO{}, fmt.Errorf("error shortening url: %w", err)
	}

	return toLinkDTO(
		updatedLink.ID,
		updatedLink.UserName,
		updatedLink.TargetUrl,
		updatedLink.ShortenedUrl.String,
		updatedLink.Status,
		updatedLink.ExpirationDate.Time,
		updatedLink.CreatedAt.Time,
		updatedLink.UpdatedAt.Time,
	), nil

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

	return toLinkDTO(
		link.ID,
		link.UserName,
		link.TargetUrl,
		link.ShortenedUrl.String,
		link.Status,
		link.ExpirationDate.Time,
		link.CreatedAt.Time,
		link.UpdatedAt.Time,
	), nil
}

func toLinkDTO(id int64, userName, targetURL, shortenedURL, status string, exp, created, updated time.Time) LinkDTO {
	return LinkDTO{
		ID:             id,
		UserName:       userName,
		TargetUrl:      targetURL,
		ShortenedUrl:   shortenedURL,
		Status:         status,
		ExpirationDate: exp,
		CreatedAt:      created,
		UpdatedAt:      updated,
	}
}
