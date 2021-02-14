// Package endpoint implements different endpoint services required by the user service
package endpoint

import (
	"context"

	"github.com/decentralized-cloud/user/services/business"
	"github.com/go-kit/kit/endpoint"
	commonErrors "github.com/micro-business/go-core/system/errors"
)

type endpointCreatorService struct {
	businessService business.BusinessContract
}

// NewEndpointCreatorService creates new instance of the EndpointCreatorService, setting up all dependencies and returns the instance
// businessService: Mandatory. Reference to the instance of the User  service
// Returns the new service or error if something goes wrong
func NewEndpointCreatorService(
	businessService business.BusinessContract) (EndpointCreatorContract, error) {
	if businessService == nil {
		return nil, commonErrors.NewArgumentNilError("businessService", "businessService is required")
	}

	return &endpointCreatorService{
		businessService: businessService,
	}, nil
}

// CreateUserEndpoint creates Create User endpoint
// Returns the Create User endpoint
func (service *endpointCreatorService) CreateUserEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if ctx == nil {
			return &business.CreateUserResponse{
				Err: commonErrors.NewArgumentNilError("ctx", "ctx is required"),
			}, nil
		}

		if request == nil {
			return &business.CreateUserResponse{
				Err: commonErrors.NewArgumentNilError("request", "request is required"),
			}, nil
		}

		castedRequest := request.(*business.CreateUserRequest)
		if err := castedRequest.Validate(); err != nil {
			return &business.CreateUserResponse{
				Err: commonErrors.NewArgumentErrorWithError("request", "", err),
			}, nil
		}

		return service.businessService.CreateUser(ctx, castedRequest)
	}
}

// ReadUserEndpoint creates Read User endpoint
// Returns the Read User endpoint
func (service *endpointCreatorService) ReadUserEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if ctx == nil {
			return &business.ReadUserResponse{
				Err: commonErrors.NewArgumentNilError("ctx", "ctx is required"),
			}, nil
		}

		if request == nil {
			return &business.ReadUserResponse{
				Err: commonErrors.NewArgumentNilError("request", "request is required"),
			}, nil
		}

		castedRequest := request.(*business.ReadUserRequest)
		if err := castedRequest.Validate(); err != nil {
			return &business.ReadUserResponse{
				Err: commonErrors.NewArgumentErrorWithError("request", "", err),
			}, nil
		}

		return service.businessService.ReadUser(ctx, castedRequest)
	}
}

// UpdateUserEndpoint creates Update User endpoint
// Returns the Update User endpoint
func (service *endpointCreatorService) UpdateUserEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if ctx == nil {
			return &business.UpdateUserResponse{
				Err: commonErrors.NewArgumentNilError("ctx", "ctx is required"),
			}, nil
		}

		if request == nil {
			return &business.UpdateUserResponse{
				Err: commonErrors.NewArgumentNilError("request", "request is required"),
			}, nil
		}

		castedRequest := request.(*business.UpdateUserRequest)
		if err := castedRequest.Validate(); err != nil {
			return &business.UpdateUserResponse{
				Err: commonErrors.NewArgumentErrorWithError("request", "", err),
			}, nil
		}

		return service.businessService.UpdateUser(ctx, castedRequest)
	}
}

// DeleteUserEndpoint creates Delete User endpoint
// Returns the Delete User endpoint
func (service *endpointCreatorService) DeleteUserEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if ctx == nil {
			return &business.DeleteUserResponse{
				Err: commonErrors.NewArgumentNilError("ctx", "ctx is required"),
			}, nil
		}

		if request == nil {
			return &business.DeleteUserResponse{
				Err: commonErrors.NewArgumentNilError("request", "request is required"),
			}, nil
		}

		castedRequest := request.(*business.DeleteUserRequest)
		if err := castedRequest.Validate(); err != nil {
			return &business.DeleteUserResponse{
				Err: commonErrors.NewArgumentErrorWithError("request", "", err),
			}, nil
		}

		return service.businessService.DeleteUser(ctx, castedRequest)
	}
}
