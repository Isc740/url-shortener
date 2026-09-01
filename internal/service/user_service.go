package service

import (
	"context"
	"time"

	"github.com/Isc740/url-shortener/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

func NewUserService(querier db.Querier) *UserService {
	return &UserService{
		queries: querier,
	}
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
		return UserResponse{}, err
	}

	return UserResponse{
		ID:        user.ID,
		UserName:  user.UserName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
	}, err
}

func (s *UserService) GetAll(ctx context.Context, limit int, offset int) ([]db.GetUsersRow, error) {
	return s.queries.GetUsers(ctx, db.GetUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
}

func (s *UserService) GetById(ctx context.Context, id int64) (db.GetUserByIDRow, error) {
	return s.queries.GetUserByID(ctx, id)

}

func (s *UserService) Update(ctx context.Context, i UpdateUserRequest) (db.User, error) {
	now := time.Now()
	return s.queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:        int64(i.ID),
		UserName:  i.UserName,
		Email:     i.Email,
		UpdatedAt: pgtype.Timestamptz{Time: now, Valid: true},
	})
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
