package service

import (
	"time"
)

type CreateLinkDTO struct {
	UserID         int8      `json:"user_id"`
	TargetURL      string    `json:"target_url"`
	Password       string    `json:"password"`
	Status         string    `json:"status"`
	ExpirationDate time.Time `json:"expiration_date"`
}

type UpdateLinkDTO struct {
	ID             int64     `json:"id"`
	TargetURL      string    `json:"target_url"`
	Status         string    `json:"status"`
	ExpirationDate time.Time `json:"expiration_date"`
}

type LinkDTO struct {
	ID             int64     `json:"id"`
	UserName       string    `json:"user_name"`
	TargetUrl      string    `json:"target_url"`
	ShortenedUrl   string    `json:"shortened_url"`
	Status         string    `json:"status"`
	ExpirationDate time.Time `json:"expiration_date"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
