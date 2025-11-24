package snippets

import (
	"context"
	"deplagene/snippetbox/internal/database"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	q database.Querier
}

func NewService(q database.Querier) *Service {
	return &Service{q: q}
}

func (s *Service) Create(title, content string, expiresInDays int) (uuid.UUID, error) {
	const op = "snippets.Service.Create"

	snippetID, err := s.q.Create(context.Background(), database.CreateParams{
		Title:   title,
		Content: content,
		ExpiresIn: pgtype.Interval{
			Microseconds: int64(expiresInDays) * 24 * 60 * 60 * 1_000_000,
			Valid:        true,
		},
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}
	return snippetID, nil
}

func (s *Service) GetById(id uuid.UUID) (database.Snippet, error) {
	const op = "snippets.Service.GetById"
	snippet, err := s.q.GetById(context.Background(), id)
	if err != nil {
		return database.Snippet{}, fmt.Errorf("%s: %w", op, err)
	}
	return snippet, nil
}

func (s *Service) GetLatest() ([]database.Snippet, error) {
	const op = "snippets.Service.GetLatest"
	snippets, err := s.q.GetLatest(context.Background())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return snippets, nil
}
