package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/DemianShtepa/exception-control/internal/domain"
	userRepository "github.com/DemianShtepa/exception-control/internal/repository/user"
	"log/slog"
	"time"
)

const TokenTTL = time.Hour * 12

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
)

type Storage interface {
	GetUserByEmail(context.Context, string) (*domain.User, error)
	SaveUser(context.Context, *domain.User) error
	Transaction(ctx context.Context, fn func(ctx context.Context, repository Storage) error) error
}

type Hasher interface {
	HashPassword(plainPassword string) (string, error)
	VerifyHashPassword(passwordHash, plainPassword string) error
}

type TokenGenerator interface {
	GenerateToken(domain.User, time.Duration) (string, error)
}

type Auth struct {
	log            *slog.Logger
	storage        Storage
	hasher         Hasher
	tokenGenerator TokenGenerator
}

func NewAuth(log *slog.Logger, storage Storage, hasher Hasher, tokenGenerator TokenGenerator) *Auth {
	return &Auth{log: log, storage: storage, hasher: hasher, tokenGenerator: tokenGenerator}
}

func (a *Auth) Register(ctx context.Context, email, password string) error {
	passwordHash, err := a.hasher.HashPassword(password)
	if err != nil {
		a.log.Error(fmt.Sprintf("failed to hash password: %s", err))

		return err
	}

	user := domain.User{
		Email:     email,
		Password:  passwordHash,
		CreatedAt: time.Now(),
	}
	err = a.storage.Transaction(ctx, func(ctx context.Context, repository Storage) error {
		_, err := repository.GetUserByEmail(ctx, email)
		if err != nil {
			if !errors.Is(err, userRepository.ErrUserNotFound) {
				return err
			}
		}

		return repository.SaveUser(ctx, &user)
	})
	if err != nil {
		if errors.Is(err, userRepository.ErrUserExists) {
			return ErrUserExists
		}

		a.log.Error(fmt.Sprintf("failed to register user: %s", err))

		return err
	}

	return nil
}

func (a *Auth) Login(ctx context.Context, email, password string) (string, error) {
	user, err := a.storage.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, userRepository.ErrUserNotFound) {
			return "", ErrInvalidCredentials
		}

		a.log.Error(fmt.Sprintf("failed to get user: %s", err))

		return "", err
	}

	if err = a.hasher.VerifyHashPassword(user.Password, password); err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := a.tokenGenerator.GenerateToken(*user, TokenTTL)
	if err != nil {
		a.log.Error(fmt.Sprintf("failed to generate token: %s", err))

		return "", err
	}

	return token, nil
}
