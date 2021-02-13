// Package endpoint implements different endpoint services required by the user service
package endpoint

import "github.com/go-kit/kit/endpoint"

// EndpointCreatorContract declares the contract that creates endpoints to create new user,
// read, update and delete existing users.
type EndpointCreatorContract interface {
	// CreateUserEndpoint creates Create User endpoint
	// Returns the Create User endpoint
	CreateUserEndpoint() endpoint.Endpoint

	// ReadUserEndpoint creates Read User endpoint
	// Returns the Read User endpoint
	ReadUserEndpoint() endpoint.Endpoint

	// ReadUserByEmailEndpoint creates Read User By Email endpoint
	// Returns the Read User By Email endpoint
	ReadUserByEmailEndpoint() endpoint.Endpoint

	// UpdateUserEndpoint creates Update User endpoint
	// Returns the Update User endpoint
	UpdateUserEndpoint() endpoint.Endpoint

	// DeleteUserEndpoint creates Delete User endpoint
	// Returns the Delete User endpoint
	DeleteUserEndpoint() endpoint.Endpoint

	// SearchEndpoint creates Search User endpoint
	// Returns the Search User endpoint
	SearchEndpoint() endpoint.Endpoint
}
