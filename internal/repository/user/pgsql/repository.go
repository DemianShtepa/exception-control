package pgsql

import (
	"context"
	"errors"
	"github.com/DemianShtepa/exception-control/internal/app/database/sqlc"
	"github.com/DemianShtepa/exception-control/internal/domain"
	userReppository "github.com/DemianShtepa/exception-control/internal/repository/user"
	"github.com/DemianShtepa/exception-control/internal/services/auth"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db      *pgxpool.Pool
	queries *sqlc.Queries
}

func NewRepository(db *pgxpool.Pool, queries *sqlc.Queries) *Repository {
	return &Repository{db: db, queries: queries}
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, userReppository.ErrUserNotFound
	}

	return convertUser(user), nil
}

func (r *Repository) SaveUser(ctx context.Context, user *domain.User) error {
	var createdAtTimestamp pgtype.Timestamp
	err := createdAtTimestamp.Scan(user.CreatedAt)
	if err != nil {
		return err
	}

	dbUser, err := r.queries.CreateUser(ctx, sqlc.CreateUserParams{
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: createdAtTimestamp,
	})
	if err != nil {
		var psqlError *pgconn.PgError
		if errors.As(err, &psqlError) && psqlError.Code == "23505" {
			return userReppository.ErrUserExists
		}
		return err
	}

	user.ID = dbUser.ID

	return nil
}

func (r *Repository) Transaction(ctx context.Context, fn func(ctx context.Context, repository auth.Storage) error) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		rollbackError := tx.Rollback(ctx)

		err = errors.Join(err, rollbackError)
	}()

	repository := Repository{
		db:      r.db,
		queries: r.queries.WithTx(tx),
	}
	err = fn(ctx, &repository)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func convertUser(user sqlc.User) *domain.User {
	return &domain.User{
		ID:        user.ID,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt.Time,
	}
}
