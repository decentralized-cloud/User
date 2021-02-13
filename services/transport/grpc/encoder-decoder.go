// Package grpc implements functions to expose user service endpoint using GRPC protocol.
package grpc

import (
	"context"

	userGRPCContract "github.com/decentralized-cloud/user/contract/grpc/go"
	"github.com/decentralized-cloud/user/models"
	"github.com/decentralized-cloud/user/services/business"
	"github.com/micro-business/go-core/common"
	commonErrors "github.com/micro-business/go-core/system/errors"
	"github.com/thoas/go-funk"
)

// decodeCreateUserRequest decodes CreateUser request message from GRPC object to business object
// context: Mandatory The reference to the context
// request: Mandatory. The reference to the GRPC request
// Returns either the decoded request or error if something goes wrong
func decodeCreateUserRequest(
	ctx context.Context,
	request interface{}) (interface{}, error) {
	castedRequest := request.(*userGRPCContract.CreateUserRequest)

	return &business.CreateUserRequest{
		User: models.User{
			Email: castedRequest.User.Email,
		}}, nil
}

// encodeCreateUserResponse encodes CreateUser response from business object to GRPC object
// context: Optional The reference to the context
// request: Mandatory. The reference to the business response
// Returns either the decoded response or error if something goes wrong
func encodeCreateUserResponse(
	ctx context.Context,
	response interface{}) (interface{}, error) {
	castedResponse := response.(*business.CreateUserResponse)

	if castedResponse.Err == nil {
		return &userGRPCContract.CreateUserResponse{
			Error:  userGRPCContract.Error_NO_ERROR,
			UserID: castedResponse.UserID,
			User: &userGRPCContract.User{
				Email: castedResponse.User.Email,
			},
			Cursor: castedResponse.Cursor,
		}, nil
	}

	return &userGRPCContract.CreateUserResponse{
		Error:        mapError(castedResponse.Err),
		ErrorMessage: castedResponse.Err.Error(),
	}, nil
}

// decodeReadUserRequest decodes ReadUser request message from GRPC object to business object
// context: Optional The reference to the context
// request: Mandatory. The reference to the GRPC request
// Returns either the decoded request or error if something goes wrong
func decodeReadUserRequest(
	ctx context.Context,
	request interface{}) (interface{}, error) {
	castedRequest := request.(*userGRPCContract.ReadUserRequest)

	return &business.ReadUserRequest{
		UserID: castedRequest.UserID,
	}, nil
}

// encodeReadUserResponse encodes ReadUser response from business object to GRPC object
// context: Optional The reference to the context
// request: Mandatory. The reference to the business response
// Returns either the decoded response or error if something goes wrong
func encodeReadUserResponse(
	ctx context.Context,
	response interface{}) (interface{}, error) {
	castedResponse := response.(*business.ReadUserResponse)

	if castedResponse.Err == nil {
		return &userGRPCContract.ReadUserResponse{
			Error: userGRPCContract.Error_NO_ERROR,
			User: &userGRPCContract.User{
				Email: castedResponse.User.Email,
			},
		}, nil
	}

	return &userGRPCContract.ReadUserResponse{
		Error:        mapError(castedResponse.Err),
		ErrorMessage: castedResponse.Err.Error(),
	}, nil
}

// decodeReadUserByEmailRequest decodes ReadUserByEmail request message from GRPC object to business object
// context: Optional The reference to the context
// request: Mandatory. The reference to the GRPC request
// Returns either the decoded request or error if something goes wrong
func decodeReadUserByEmailRequest(
	ctx context.Context,
	request interface{}) (interface{}, error) {
	castedRequest := request.(*userGRPCContract.ReadUserByEmailRequest)

	return &business.ReadUserByEmailRequest{
		Email: castedRequest.Email,
	}, nil
}

// encodeReadUserByEmailResponse encodes ReadUserByEmail response from business object to GRPC object
// context: Optional The reference to the context
// request: Mandatory. The reference to the business response
// Returns either the decoded response or error if something goes wrong
func encodeReadUserByEmailResponse(
	ctx context.Context,
	response interface{}) (interface{}, error) {
	castedResponse := response.(*business.ReadUserByEmailResponse)

	if castedResponse.Err == nil {
		return &userGRPCContract.ReadUserByEmailResponse{
			Error: userGRPCContract.Error_NO_ERROR,
			User: &userGRPCContract.User{
				Email: castedResponse.User.Email,
			},
		}, nil
	}

	return &userGRPCContract.ReadUserByEmailResponse{
		Error:        mapError(castedResponse.Err),
		ErrorMessage: castedResponse.Err.Error(),
	}, nil
}

// decodeUpdateUserRequest decodes UpdateUser request message from GRPC object to business object
// context: Optional The reference to the context
// request: Mandatory. The reference to the GRPC request
// Returns either the decoded request or error if something goes wrong
func decodeUpdateUserRequest(
	ctx context.Context,
	request interface{}) (interface{}, error) {
	castedRequest := request.(*userGRPCContract.UpdateUserRequest)

	return &business.UpdateUserRequest{
		UserID: castedRequest.UserID,
		User: models.User{
			Email: castedRequest.User.Email,
		}}, nil
}

// encodeUpdateUserResponse encodes UpdateUser response from business object to GRPC object
// context: Optional The reference to the context
// request: Mandatory. The reference to the business response
// Returns either the decoded response or error if something goes wrong
func encodeUpdateUserResponse(
	ctx context.Context,
	response interface{}) (interface{}, error) {
	castedResponse := response.(*business.UpdateUserResponse)

	if castedResponse.Err == nil {
		return &userGRPCContract.UpdateUserResponse{
			Error: userGRPCContract.Error_NO_ERROR,
			User: &userGRPCContract.User{
				Email: castedResponse.User.Email,
			},
			Cursor: castedResponse.Cursor,
		}, nil
	}

	return &userGRPCContract.UpdateUserResponse{
		Error:        mapError(castedResponse.Err),
		ErrorMessage: castedResponse.Err.Error(),
	}, nil
}

// decodeDeleteUserRequest decodes DeleteUser request message from GRPC object to business object
// context: Optional The reference to the context
// request: Mandatory. The reference to the GRPC request
// Returns either the decoded request or error if something goes wrong
func decodeDeleteUserRequest(
	ctx context.Context,
	request interface{}) (interface{}, error) {
	castedRequest := request.(*userGRPCContract.DeleteUserRequest)

	return &business.DeleteUserRequest{
		UserID: castedRequest.UserID,
	}, nil
}

// encodeDeleteUserResponse encodes DeleteUser response from business object to GRPC object
// context: Optional The reference to the context
// request: Mandatory. The reference to the business response
// Returns either the decoded response or error if something goes wrong
func encodeDeleteUserResponse(
	ctx context.Context,
	response interface{}) (interface{}, error) {
	castedResponse := response.(*business.DeleteUserResponse)
	if castedResponse.Err == nil {
		return &userGRPCContract.DeleteUserResponse{
			Error: userGRPCContract.Error_NO_ERROR,
		}, nil
	}

	return &userGRPCContract.DeleteUserResponse{
		Error:        mapError(castedResponse.Err),
		ErrorMessage: castedResponse.Err.Error(),
	}, nil
}

// decodeSearchRequest decodes Search request message from GRPC object to business object
// context: Optional The reference to the context
// request: Mandatory. The reference to the GRPC request
// Returns either the decoded request or error if something goes wrongw
func decodeSearchRequest(
	ctx context.Context,
	request interface{}) (interface{}, error) {
	castedRequest := request.(*userGRPCContract.SearchRequest)
	sortingOptions := []common.SortingOptionPair{}

	if len(castedRequest.SortingOptions) > 0 {
		sortingOptions = funk.Map(
			castedRequest.SortingOptions,
			func(sortingOption *userGRPCContract.SortingOptionPair) common.SortingOptionPair {
				direction := common.Ascending

				if sortingOption.Direction == userGRPCContract.SortingDirection_DESCENDING {
					direction = common.Descending
				}

				return common.SortingOptionPair{
					Name:      sortingOption.Name,
					Direction: direction,
				}
			}).([]common.SortingOptionPair)
	}

	pagination := common.Pagination{}

	if castedRequest.Pagination.HasAfter {
		pagination.After = &castedRequest.Pagination.After
	}

	if castedRequest.Pagination.HasFirst {
		first := int(castedRequest.Pagination.First)
		pagination.First = &first
	}

	if castedRequest.Pagination.HasBefore {
		pagination.Before = &castedRequest.Pagination.Before
	}

	if castedRequest.Pagination.HasLast {
		last := int(castedRequest.Pagination.Last)
		pagination.Last = &last
	}

	return &business.SearchRequest{
		Pagination:     pagination,
		UserIDs:        castedRequest.UserIDs,
		SortingOptions: sortingOptions,
	}, nil
}

// encodeSearchResponse encodes Search response from business object to GRPC object
// context: Optional The reference to the context
// request: Mandatory. The reference to the business response
// Returns either the decoded response or error if something goes wrong
func encodeSearchResponse(
	ctx context.Context,
	response interface{}) (interface{}, error) {
	castedResponse := response.(*business.SearchResponse)
	if castedResponse.Err == nil {
		return &userGRPCContract.SearchResponse{
			Error:           userGRPCContract.Error_NO_ERROR,
			HasPreviousPage: castedResponse.HasPreviousPage,
			HasNextPage:     castedResponse.HasNextPage,
			TotalCount:      castedResponse.TotalCount,
			Users: funk.Map(castedResponse.Users, func(user models.UserWithCursor) *userGRPCContract.UserWithCursor {
				return &userGRPCContract.UserWithCursor{
					UserID: user.UserID,
					User: &userGRPCContract.User{
						Email: user.User.Email,
					},
					Cursor: user.Cursor,
				}
			}).([]*userGRPCContract.UserWithCursor),
		}, nil
	}

	return &userGRPCContract.SearchResponse{
		Error:        mapError(castedResponse.Err),
		ErrorMessage: castedResponse.Err.Error(),
	}, nil
}

func mapError(err error) userGRPCContract.Error {
	if business.IsUnknownError(err) {
		return userGRPCContract.Error_UNKNOWN
	}

	if business.IsUserAlreadyExistsError(err) {
		return userGRPCContract.Error_USER_ALREADY_EXISTS
	}

	if business.IsUserNotFoundError(err) {
		return userGRPCContract.Error_USER_NOT_FOUND
	}

	if commonErrors.IsArgumentNilError(err) || commonErrors.IsArgumentError(err) {
		return userGRPCContract.Error_BAD_REQUEST
	}

	panic("Error type undefined.")
}
