package users

import (
	"github.com/deplagene/snippetbox/internal/database"
	"github.com/google/uuid"
)

type UserService interface {
	Register(name, email, password string) (uuid.UUID, error)
	Login(email, password string) (database.User, error)
}
