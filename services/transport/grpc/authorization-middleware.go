// Package grpc implements functions to expose user service endpoint using GRPC protocol.
package grpc

import (
	"context"

	userGRPCContract "github.com/decentralized-cloud/user/contract/grpc/go"
	"github.com/go-kit/kit/endpoint"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/micro-business/go-core/pkg/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authorizeFunc func(email string, request interface{}) error

var authorizedFuncs = map[string]authorizeFunc{
	"CreateUser": isAuthorizedToCallCreateUser,
	"ReadUser":   isAuthorizedToCallReadUser,
	"UpdateUser": isAuthorizedToCallUpdateUser,
	"DeleteUser": isAuthorizedToCallDeleteUser,
}

// CreateLoggingMiddleware creates the logging middleware.
// endpointName: Mandatory. The name of the endpoint
// Returns the new endpoint with logging middleware added
func (service *transportService) createAuthMiddleware(endpointName string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			token, err := util.ParseAndVerifyToken(ctx, service.jwksURL, true)
			if err != nil {
				return nil, err
			}

			if err = service.isAuthorized(token, endpointName, request); err != nil {
				return nil, err
			}

			return next(ctx, request)
		}
	}
}

func (service *transportService) isAuthorized(token jwt.Token, endpointName string, request interface{}) error {
	email := token.PrivateClaims()["email"].(string)

	if len(email) == 0 {
		return status.Errorf(codes.Unauthenticated, "Email address is not included in the claims")
	}

	return authorizedFuncs[endpointName](email, request)
}

func isAuthorizedToCallCreateUser(email string, request interface{}) error {
	return nil
}

func isAuthorizedToCallReadUser(email string, request interface{}) error {
	castedRequest := request.(*userGRPCContract.ReadUserRequest)

	if castedRequest.Email != email {
		return status.Errorf(codes.Unauthenticated, "Email address does not match the received one in the request")
	}

	return nil
}

func isAuthorizedToCallUpdateUser(email string, request interface{}) error {
	castedRequest := request.(*userGRPCContract.UpdateUserRequest)

	if castedRequest.Email != email {
		return status.Errorf(codes.Unauthenticated, "Email address does not match the received one in the request")
	}

	return nil
}

func isAuthorizedToCallDeleteUser(email string, request interface{}) error {
	castedRequest := request.(*userGRPCContract.DeleteUserRequest)

	if castedRequest.Email != email {
		return status.Errorf(codes.Unauthenticated, "Email address does not match the received one in the request")
	}

	return nil
}
