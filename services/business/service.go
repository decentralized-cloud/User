// Package business implements different business services required by the user service
package business

import (
	"context"

	"github.com/decentralized-cloud/user/services/repository"
	commonErrors "github.com/micro-business/go-core/system/errors"
)

type businessService struct {
	repositoryService repository.RepositoryContract
}

// NewBusinessService creates new instance of the BusinessService, setting up all dependencies and returns the instance
// repositoryService: Mandatory. Reference to the repository service that can persist the user related data
// Returns the new service or error if something goes wrong
func NewBusinessService(
	repositoryService repository.RepositoryContract) (BusinessContract, error) {
	if repositoryService == nil {
		return nil, commonErrors.NewArgumentNilError("repositoryService", "repositoryService is required")
	}

	return &businessService{
		repositoryService: repositoryService,
	}, nil
}

// CreateUser creates a new user.
// ctx: Mandatory The reference to the context
// request: Mandatory. The request to create a new user
// Returns either the result of creating new user or error if something goes wrong.
func (service *businessService) CreateUser(
	ctx context.Context,
	request *CreateUserRequest) (*CreateUserResponse, error) {
	response, err := service.repositoryService.CreateUser(ctx, &repository.CreateUserRequest{
		User: request.User,
	})

	if err != nil {
		return &CreateUserResponse{
			Err: mapRepositoryError(err, "", ""),
		}, nil
	}

	return &CreateUserResponse{
		UserID: response.UserID,
		User:   response.User,
		Cursor: response.Cursor,
	}, nil
}

// ReadUser read an existing user
// ctx: Mandatory The reference to the context
// request: Mandatory. The request to read an existing user
// Returns either the result of reading an existing user or error if something goes wrong.
func (service *businessService) ReadUser(
	ctx context.Context,
	request *ReadUserRequest) (*ReadUserResponse, error) {
	response, err := service.repositoryService.ReadUser(ctx, &repository.ReadUserRequest{
		UserID: request.UserID,
	})

	if err != nil {
		return &ReadUserResponse{
			Err: mapRepositoryError(err, request.UserID, ""),
		}, nil
	}

	return &ReadUserResponse{
		User: response.User,
	}, nil
}

// ReadUserByEmail read an existing user by email address
// ctx: Mandatory The reference to the context
// request: Mandatory. The request to read an existing user by email address
// Returns either the result of reading an existing user by email address or error if something goes wrong.
func (service *businessService) ReadUserByEmail(
	ctx context.Context,
	request *ReadUserByEmailRequest) (*ReadUserByEmailResponse, error) {
	response, err := service.repositoryService.ReadUserByEmail(ctx, &repository.ReadUserByEmailRequest{
		Email: request.Email,
	})

	if err != nil {
		return &ReadUserByEmailResponse{
			Err: mapRepositoryError(err, "", request.Email),
		}, nil
	}

	return &ReadUserByEmailResponse{
		UserID: response.UserID,
		User:   response.User,
	}, nil
}

// UpdateUser update an existing user
// ctx: Mandatory The reference to the context
// request: Mandatory. The request to update an existing user
// Returns either the result of updateing an existing user or error if something goes wrong.
func (service *businessService) UpdateUser(
	ctx context.Context,
	request *UpdateUserRequest) (*UpdateUserResponse, error) {
	response, err := service.repositoryService.UpdateUser(ctx, &repository.UpdateUserRequest{
		UserID: request.UserID,
		User:   request.User,
	})

	if err != nil {
		return &UpdateUserResponse{
			Err: mapRepositoryError(err, request.UserID, ""),
		}, nil
	}

	return &UpdateUserResponse{
		User:   response.User,
		Cursor: response.Cursor,
	}, nil
}

// DeleteUser delete an existing user
// ctx: Mandatory The reference to the context
// request: Mandatory. The request to delete an existing user
// Returns either the result of deleting an existing user or error if something goes wrong.
func (service *businessService) DeleteUser(
	ctx context.Context,
	request *DeleteUserRequest) (*DeleteUserResponse, error) {
	_, err := service.repositoryService.DeleteUser(ctx, &repository.DeleteUserRequest{
		UserID: request.UserID,
	})

	if err != nil {
		return &DeleteUserResponse{
			Err: mapRepositoryError(err, request.UserID, ""),
		}, nil
	}

	return &DeleteUserResponse{}, nil
}

// Search returns the list of users that matched the criteria
// ctx: Mandatory The reference to the context
// request: Mandatory. The request contains the search criteria
// Returns the list of users that matched the criteria
func (service *businessService) Search(
	ctx context.Context,
	request *SearchRequest) (*SearchResponse, error) {
	result, err := service.repositoryService.Search(ctx, &repository.SearchRequest{
		Pagination:     request.Pagination,
		SortingOptions: request.SortingOptions,
		UserIDs:        request.UserIDs,
	})

	if err != nil {
		return &SearchResponse{
			Err: mapRepositoryError(err, "", ""),
		}, nil
	}

	return &SearchResponse{
		HasPreviousPage: result.HasPreviousPage,
		HasNextPage:     result.HasNextPage,
		TotalCount:      result.TotalCount,
		Users:           result.Users,
	}, nil
}

func mapRepositoryError(err error, userID string, email string) error {
	if repository.IsUserAlreadyExistsError(err) {
		return NewUserAlreadyExistsErrorWithError(err)
	}

	if repository.IsUserNotFoundError(err) {
		return NewUserNotFoundErrorWithError(userID, err)
	}

	if repository.IsUserByEmailNotFoundError(err) {
		return NewUserByEmailNotFoundErrorWithError(email, err)
	}

	return NewUnknownErrorWithError("", err)
}
