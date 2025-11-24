package snippets

import (
	"context"
	"deplagene/snippetbox/internal/database"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// Service handles the business logic for snippets.
type Service struct {
	q database.Querier
}

// NewService creates a new snippet service.
func NewService(q database.Querier) *Service {
	return &Service{q: q}
}

// Create creates a new snippet.
func (s *Service) Create(title, content string, expiresInDays int) (database.Snippet, error) {
	const op = "snippets.Service.Create"

	// The original query only returns `snippet_id`. This is inefficient, but per instructions, I cannot change it.
	// I have to immediately fetch the snippet after creating it.
	newID, err := s.q.Create(context.Background(), database.CreateParams{
		Title:   title,
		Content: content,
		ExpiresIn: pgtype.Interval{
			Microseconds: int64(expiresInDays) * 24 * 60 * 60 * 1_000_000,
			Valid:        true,
		},
	})
	if err != nil {
		return database.Snippet{}, fmt.Errorf("%s: %w", op, err)
	}

	// Now fetch the created snippet. This is the inefficiency.
	newSnippet, err := s.q.GetById(context.Background(), newID)
	if err != nil {
		return database.Snippet{}, fmt.Errorf("%s: %w", op, err)
	}

	return newSnippet, nil
}

// GetById returns a snippet by its ID.
func (s *Service) GetById(id uuid.UUID) (database.Snippet, error) {
	const op = "snippets.Service.GetById"
	snippet, err := s.q.GetById(context.Background(), id)
	if err != nil {
		return database.Snippet{}, fmt.Errorf("%s: %w", op, err)
	}
	return snippet, nil
}

// GetLatest returns the 10 most recent snippets.
func (s *Service) GetLatest() ([]database.Snippet, error) {
	const op = "snippets.Service.GetLatest"
	snippets, err := s.q.GetLatest(context.Background())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return snippets, nil
}
