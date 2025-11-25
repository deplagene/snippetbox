package snippets

import (
	"context"
	"fmt"

	"github.com/deplagene/snippetbox/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	q database.Querier
}

func NewService(q database.Querier) *Service {
	return &Service{q: q}
}

func (s *Service) Create(title, content string, userID uuid.UUID, expiresInDays int) (uuid.UUID, error) {
	const op = "snippets.Service.Create"

	snippetID, err := s.q.CreateSnippet(context.Background(), database.CreateSnippetParams{
		Title:   title,
		Content: content,
		UserID:  pgtype.UUID{Bytes: userID, Valid: true},
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
	snippet, err := s.q.GetSnippetByID(context.Background(), id)
	if err != nil {
		return database.Snippet{}, fmt.Errorf("%s: %w", op, err)
	}
	return database.Snippet{
		SnippetID: snippet.SnippetID,
		UserID:    snippet.UserID,
		Title:     snippet.Title,
		Content:   snippet.Content,
	}, nil
}

func (s *Service) GetLatestForUser(userID uuid.UUID) ([]database.Snippet, error) {
	const op = "snippets.Service.GetLatestForUser"
	snippets, err := s.q.GetLatestSnippetsForUser(context.Background(), pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	var result []database.Snippet
	for _, snip := range snippets {
		result = append(result, database.Snippet{
			SnippetID: snip.SnippetID,
			UserID:    pgtype.UUID{Bytes: userID, Valid: true},
			Title:     snip.Title,
			Content:   snip.Content,
		})
	}
	return result, nil
}

func (s *Service) GetLatest() ([]database.Snippet, error) {
	const op = "snippets.Service.GetLatest"
	snippets, err := s.q.GetLatestSnippets(context.Background())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	var result []database.Snippet
	for _, snip := range snippets {
		result = append(result, database.Snippet{
			SnippetID: snip.SnippetID,
			UserID:    snip.UserID,
			Title:     snip.Title,
			Content:   snip.Content,
		})
	}
	return result, nil
}

func (s *Service) Delete(id uuid.UUID) error {
	const op = "snippets.Service.Delete"

	if err := s.q.DeleteSnippet(context.Background(), id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}