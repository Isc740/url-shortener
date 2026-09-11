package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Isc740/url-shortener/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserService struct {
	queries db.Querier
}

func NewUserService(querier db.Querier) *UserService {
	return &UserService{
		queries: querier,
	}
}

var (
	ErrUsersNotFound  = errors.New("links not found")
	ErrUserExpired    = errors.New("link has expired")
	ErrUserInactive   = errors.New("link is inactive")
	ErrUserNotFound   = errors.New("user does not exist")
	ErrDuplicateEmail = errors.New("email already in use")
)

var userCreateErrors = map[string]error{
	"23505": ErrDuplicateEmail,
}

func (s *UserService) GetAll(ctx context.Context, limit int, offset int) ([]db.GetUsersRow, error) {
	users, err := s.queries.GetUsers(ctx, db.GetUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("Error getting users: %w", err)
	}
	return users, nil
}

func (s *UserService) GetById(ctx context.Context, id int) (db.GetUserByIDRow, error) {
	user, err := s.queries.GetUserByID(ctx, int64(id))
	if err != nil {
		return db.GetUserByIDRow{}, wrapNotFound(err, ErrUserNotFound, "error getting user")
	}
	return user, nil
}

func (s *UserService) GetByUserName(ctx context.Context, name string) (db.GetUserByNameRow, error) {
	user, err := s.queries.GetUserByName(ctx, name)
	if err != nil {
		return db.GetUserByNameRow{}, wrapNotFound(err, ErrUserNotFound, "error getting user")
	}
	return user, nil
}

func (s *UserService) GetByEmailForUserAuth(ctx context.Context, email string) (db.GetUserByEmailForAuthRow, error) {
	user, err := s.queries.GetUserByEmailForAuth(ctx, email)
	if err != nil {
		return db.GetUserByEmailForAuthRow{}, wrapNotFound(err, ErrUserNotFound, "error getting user")
	}
	return user, nil
}

func (s *UserService) Create(ctx context.Context, i CreateUserRequest) (UserResponse, error) {
	now := time.Now()
	user, err := s.queries.CreateUser(ctx, db.CreateUserParams{
		UserName:  i.UserName,
		Email:     i.Email,
		Password:  i.Password,
		CreatedAt: pgtype.Timestamptz{Time: now, Valid: true},
		UpdatedAt: pgtype.Timestamptz{Time: now, Valid: true},
	})

	if err != nil {
		return wrapCreateError[UserResponse](err, "creating user", userCreateErrors)
	}

	return UserResponse{
		ID:        user.ID,
		UserName:  user.UserName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
	}, err
}

func (s *UserService) Update(ctx context.Context, i UpdateUserRequest) (UserResponse, error) {
	now := time.Now()
	user, err := s.queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:        int64(i.ID),
		UserName:  i.UserName,
		Email:     i.Email,
		UpdatedAt: pgtype.Timestamptz{Time: now, Valid: true},
	})
	if err != nil {
		return UserResponse{}, fmt.Errorf("error updating user: %w", err)
	}

	return UserResponse{
		ID:        user.ID,
		UserName:  user.UserName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
	}, nil
}

func (s *UserService) Delete(ctx context.Context, id int) error {
	now := time.Now()
	return s.queries.DeleteUser(ctx, db.DeleteUserParams{
		ID:        int64(id),
		DeletedAt: pgtype.Timestamptz{Time: now, Valid: true},
	})
}

func (s *UserService) Restore(ctx context.Context, id int) error {
	return s.queries.RestoreUser(ctx, int64(id))
}

func (s *UserService) Destroy(ctx context.Context, id int) error {
	return s.queries.DestroyUser(ctx, int64(id))
}
