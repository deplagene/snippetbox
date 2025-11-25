package snippets

import (
	"github.com/deplagene/snippetbox/internal/database"
	"github.com/google/uuid"
)

type SnippetService interface {
	Create(title, content string, userID uuid.UUID, expiresInDays int) (uuid.UUID, error)
	GetById(id uuid.UUID) (database.Snippet, error)
	GetLatest() ([]database.Snippet, error)
	GetLatestForUser(userID uuid.UUID) ([]database.Snippet, error)
	Delete(id uuid.UUID) error
}
