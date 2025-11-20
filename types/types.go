package types

import (
	"time"
)

type SnippetDTO struct {
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Created time.Time `json:"created"`
}
