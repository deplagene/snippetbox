package users

import (
	"context"
	"errors"
	"fmt"
	"github.com/deplagene/snippetbox/internal/database"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	q database.Querier
}

func NewService(q database.Querier) *Service {
	return &Service{q: q}
}

func (s *Service) Register(name, email, password string) (uuid.UUID, error) {
	const op = "users.Service.Register"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	userID, err := s.q.CreateUser(context.Background(), database.CreateUserParams{
		Name:          name,
		Email:         email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}
	return userID, nil
}

func (s *Service) Login(email, password string) (database.User, error) {
	const op = "users.Service.Login"
	user, err := s.q.GetUserByEmail(context.Background(), email)
	if err != nil {
		return database.User{}, fmt.Errorf("%s: %w", op, err)
	}

	err = bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return database.User{}, fmt.Errorf("%s: %w", op, errors.New("invalid credentials"))
		}
		return database.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return database.User{
		UserID: user.UserID,
		Name:   user.Name,
		Email:  user.Email,
	}, nil
}
