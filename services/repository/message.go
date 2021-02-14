// Package repository implements different repository services required by the user service
package repository

import (
	"github.com/decentralized-cloud/user/models"
)

// CreateUserRequest contains the request to create a new user
type CreateUserRequest struct {
	Email string
	User  models.User
}

// CreateUserResponse contains the result of creating a new user
type CreateUserResponse struct {
	User   models.User
	Cursor string
}

// ReadUserRequest contains the request to read an existing user
type ReadUserRequest struct {
	Email string
}

// ReadUserResponse contains the result of reading an existing user
type ReadUserResponse struct {
	User models.User
}

// UpdateUserRequest contains the request to update an existing user
type UpdateUserRequest struct {
	Email string
	User  models.User
}

// UpdateUserResponse contains the result of updating an existing user
type UpdateUserResponse struct {
	User   models.User
	Cursor string
}

// DeleteUserRequest contains the request to delete an existing user
type DeleteUserRequest struct {
	Email string
}

// DeleteUserResponse contains the result of deleting an existing user
type DeleteUserResponse struct {
}
