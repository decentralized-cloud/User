// Package business implements different business services required by the user service
package business

import "context"

// BusinessContract declares the service that can create new user, read, update
// and delete existing users.
type BusinessContract interface {
	// CreateUser creates a new user.
	// ctx: Mandatory The reference to the context
	// request: Mandatory. The request to create a new user
	// Returns either the result of creating new user or error if something goes wrong.
	CreateUser(
		ctx context.Context,
		request *CreateUserRequest) (*CreateUserResponse, error)

	// ReadUser read an existing user
	// ctx: Mandatory The reference to the context
	// request: Mandatory. The request to read an existing user
	// Returns either the result of reading an existing user or error if something goes wrong.
	ReadUser(
		ctx context.Context,
		request *ReadUserRequest) (*ReadUserResponse, error)

	// ReadUserByEmail read an existing user by email address
	// ctx: Mandatory The reference to the context
	// request: Mandatory. The request to read an existing user by email address
	// Returns either the result of reading an existing user by email address or error if something goes wrong.
	ReadUserByEmail(
		ctx context.Context,
		request *ReadUserByEmailRequest) (*ReadUserByEmailResponse, error)

	// UpdateUser update an existing user
	// ctx: Mandatory The reference to the context
	// request: Mandatory. The request to update an existing user
	// Returns either the result of updateing an existing user or error if something goes wrong.
	UpdateUser(
		ctx context.Context,
		request *UpdateUserRequest) (*UpdateUserResponse, error)

	// DeleteUser delete an existing user
	// ctx: Mandatory The reference to the context
	// request: Mandatory. The request to delete an existing user
	// Returns either the result of deleting an existing user or error if something goes wrong.
	DeleteUser(
		ctx context.Context,
		request *DeleteUserRequest) (*DeleteUserResponse, error)

	// Search returns the list of users that matched the criteria
	// ctx: Mandatory The reference to the context
	// request: Mandatory. The request contains the search criteria
	// Returns the list of users that matched the criteria
	Search(
		ctx context.Context,
		request *SearchRequest) (*SearchResponse, error)
}
