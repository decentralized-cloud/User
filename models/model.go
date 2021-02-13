// Package models defines the different object models used in User
package models

// User defines the user object
type User struct {
	Email string `bson:"email" json:"email"`
}

// UserWithCursor implements the pair of the user with a cursor that determines the
// location of the tennat in the repository.
type UserWithCursor struct {
	UserID string
	User   User
	Cursor string
}
