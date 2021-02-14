// Package grpc implements functions to expose user service endpoint using GRPC protocol.
package grpc

import (
	"context"

	userGRPCContract "github.com/decentralized-cloud/user/contract/grpc/go"
	"github.com/decentralized-cloud/user/models"
	"github.com/decentralized-cloud/user/services/business"
	"github.com/micro-business/go-core/pkg/util"
	commonErrors "github.com/micro-business/go-core/system/errors"
)

// decodeCreateUserRequest decodes CreateUser request message from GRPC object to business object
// context: Mandatory The reference to the context
// request: Mandatory. The reference to the GRPC request
// Returns either the decoded request or error if something goes wrong
func decodeCreateUserRequest(
	ctx context.Context,
	request interface{}) (interface{}, error) {
	token, err := util.ParseAndVerifyToken(ctx, "", false)
	if err != nil {
		return nil, err
	}

	return &business.CreateUserRequest{
		Email: token.PrivateClaims()["email"].(string),
		User:  models.User{}}, nil
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
			User:   &userGRPCContract.User{},
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
		Email: castedRequest.Email,
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
			User:  &userGRPCContract.User{},
		}, nil
	}

	return &userGRPCContract.ReadUserResponse{
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
		Email: castedRequest.Email,
		User:  models.User{}}, nil
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
			Error:  userGRPCContract.Error_NO_ERROR,
			User:   &userGRPCContract.User{},
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
		Email: castedRequest.Email,
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
