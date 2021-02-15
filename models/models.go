// Package models defines the different object models used in User
package models

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

var (
	// ContextKeyParsedToken var
	ContextKeyParsedToken = contextKey("ParsedToken")
)

// ParsedToken contains details that are encoded in the received JWT token
type ParsedToken struct {
	Email string
}

// User defines the user object
type User struct {
}

// UserWithCursor implements the pair of the user with a cursor that determines the
// location of the tennat in the repository.
type UserWithCursor struct {
	UserID string
	User   User
	Cursor string
}
