// Package business implements different business services required by the user service
package business

import (
	"github.com/decentralized-cloud/user/models"
	"github.com/micro-business/go-core/common"
)

// CreateUserRequest contains the request to create a new user
type CreateUserRequest struct {
	User models.User
}

// CreateUserResponse contains the result of creating a new user
type CreateUserResponse struct {
	Err    error
	UserID string
	User   models.User
	Cursor string
}

// ReadUserRequest contains the request to read an existing user
type ReadUserRequest struct {
	UserID string
}

// ReadUserResponse contains the result of reading an existing user
type ReadUserResponse struct {
	Err  error
	User models.User
}

// ReadUserByEmailRequest contains the request to read an existing user by email address
type ReadUserByEmailRequest struct {
	Email string
}

// ReadUserByEmailResponse contains the result of reading an existing user by email address
type ReadUserByEmailResponse struct {
	Err    error
	UserID string
	User   models.User
}

// UpdateUserRequest contains the request to update an existing user
type UpdateUserRequest struct {
	UserID string
	User   models.User
}

// UpdateUserResponse contains the result of updating an existing user
type UpdateUserResponse struct {
	Err    error
	User   models.User
	Cursor string
}

// DeleteUserRequest contains the request to delete an existing user
type DeleteUserRequest struct {
	UserID string
}

// DeleteUserResponse contains the result of deleting an existing user
type DeleteUserResponse struct {
	Err error
}

// SearchRequest contains the filter criteria to look for existing users
type SearchRequest struct {
	Pagination     common.Pagination
	SortingOptions []common.SortingOptionPair
	UserIDs        []string
}

// SearchResponse contains the list of the users that matched the result
type SearchResponse struct {
	Err             error
	HasPreviousPage bool
	HasNextPage     bool
	TotalCount      int64
	Users           []models.UserWithCursor
}
