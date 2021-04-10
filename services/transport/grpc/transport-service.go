// Package grpc implements functions to expose user service endpoint using GRPC protocol.
package grpc

import (
	"context"
	"fmt"
	"net"

	userGRPCContract "github.com/decentralized-cloud/user/contract/grpc/go"
	"github.com/decentralized-cloud/user/services/configuration"
	"github.com/decentralized-cloud/user/services/endpoint"
	"github.com/decentralized-cloud/user/services/transport"
	gokitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/micro-business/go-core/gokit/middleware"
	commonErrors "github.com/micro-business/go-core/system/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type transportService struct {
	logger                    *zap.Logger
	configurationService      configuration.ConfigurationContract
	endpointCreatorService    endpoint.EndpointCreatorContract
	middlewareProviderService middleware.MiddlewareProviderContract
	jwksURL                   string
	createUserHandler         gokitgrpc.Handler
	readUserHandler           gokitgrpc.Handler
	updateUserHandler         gokitgrpc.Handler
	deleteUserHandler         gokitgrpc.Handler
}

var Live bool
var Ready bool

func init() {
	Live = false
	Ready = false
}

// NewTransportService creates new instance of the transportService, setting up all dependencies and returns the instance
// logger: Mandatory. Reference to the logger service
// configurationService: Mandatory. Reference to the service that provides required configurations
// endpointCreatorService: Mandatory. Reference to the service that creates go-kit compatible endpoints
// middlewareProviderService: Mandatory. Reference to the service that provides different go-kit middlewares
// Returns the new service or error if something goes wrong
func NewTransportService(
	logger *zap.Logger,
	configurationService configuration.ConfigurationContract,
	endpointCreatorService endpoint.EndpointCreatorContract,
	middlewareProviderService middleware.MiddlewareProviderContract) (transport.TransportContract, error) {
	if logger == nil {
		return nil, commonErrors.NewArgumentNilError("logger", "logger is required")
	}

	if configurationService == nil {
		return nil, commonErrors.NewArgumentNilError("configurationService", "configurationService is required")
	}

	if endpointCreatorService == nil {
		return nil, commonErrors.NewArgumentNilError("endpointCreatorService", "endpointCreatorService is required")
	}

	if middlewareProviderService == nil {
		return nil, commonErrors.NewArgumentNilError("middlewareProviderService", "middlewareProviderService is required")
	}

	jwksURL, err := configurationService.GetJwksURL()
	if err != nil {
		return nil, err
	}

	return &transportService{
		logger:                    logger,
		configurationService:      configurationService,
		endpointCreatorService:    endpointCreatorService,
		middlewareProviderService: middlewareProviderService,
		jwksURL:                   jwksURL,
	}, nil
}

// Start starts the GRPC transport service
// Returns error if something goes wrong
func (service *transportService) Start() error {
	service.setupHandlers()

	host, err := service.configurationService.GetGrpcHost()
	if err != nil {
		return err
	}

	port, err := service.configurationService.GetGrpcPort()
	if err != nil {
		return err
	}

	address := fmt.Sprintf("%s:%d", host, port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	gRPCServer := grpc.NewServer()
	userGRPCContract.RegisterServiceServer(gRPCServer, service)
	service.logger.Info("gRPC service started", zap.String("address", address))

	Live = true
	Ready = true

	err = gRPCServer.Serve(listener)

	Live = false
	Ready = false

	return err
}

// Stop stops the GRPC transport service
// Returns error if something goes wrong
func (service *transportService) Stop() error {
	return nil
}

func (service *transportService) setupHandlers() {
	endpoint := service.endpointCreatorService.CreateUserEndpoint()
	endpoint = service.middlewareProviderService.CreateLoggingMiddleware("CreateUser")(endpoint)
	endpoint = service.createAuthMiddleware("CreateUser")(endpoint)
	service.createUserHandler = gokitgrpc.NewServer(
		endpoint,
		decodeCreateUserRequest,
		encodeCreateUserResponse,
	)

	endpoint = service.endpointCreatorService.ReadUserEndpoint()
	endpoint = service.middlewareProviderService.CreateLoggingMiddleware("ReadUser")(endpoint)
	endpoint = service.createAuthMiddleware("ReadUser")(endpoint)
	service.readUserHandler = gokitgrpc.NewServer(
		endpoint,
		decodeReadUserRequest,
		encodeReadUserResponse,
	)

	endpoint = service.endpointCreatorService.UpdateUserEndpoint()
	endpoint = service.middlewareProviderService.CreateLoggingMiddleware("UpdateUser")(endpoint)
	endpoint = service.createAuthMiddleware("UpdateUser")(endpoint)
	service.updateUserHandler = gokitgrpc.NewServer(
		endpoint,
		decodeUpdateUserRequest,
		encodeUpdateUserResponse,
	)

	endpoint = service.endpointCreatorService.DeleteUserEndpoint()
	endpoint = service.middlewareProviderService.CreateLoggingMiddleware("DeleteUser")(endpoint)
	endpoint = service.createAuthMiddleware("DeleteUser")(endpoint)
	service.deleteUserHandler = gokitgrpc.NewServer(
		endpoint,
		decodeDeleteUserRequest,
		encodeDeleteUserResponse,
	)
}

// CreateUser creates a new user
// context: Mandatory. The reference to the context
// request: mandatory. The request to create a new user
// Returns the result of creating new user
func (service *transportService) CreateUser(
	ctx context.Context,
	request *userGRPCContract.CreateUserRequest) (*userGRPCContract.CreateUserResponse, error) {
	_, response, err := service.createUserHandler.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.(*userGRPCContract.CreateUserResponse), nil
}

// ReadUser read an existing user
// context: Mandatory. The reference to the context
// request: Mandatory. The request to read an existing user
// Returns the result of reading an existing user
func (service *transportService) ReadUser(
	ctx context.Context,
	request *userGRPCContract.ReadUserRequest) (*userGRPCContract.ReadUserResponse, error) {
	_, response, err := service.readUserHandler.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.(*userGRPCContract.ReadUserResponse), nil

}

// UpdateUser update an existing user
// context: Mandatory. The reference to the context
// request: Mandatory. The request to update an existing user
// Returns the result of updateing an existing user
func (service *transportService) UpdateUser(
	ctx context.Context,
	request *userGRPCContract.UpdateUserRequest) (*userGRPCContract.UpdateUserResponse, error) {
	_, response, err := service.updateUserHandler.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.(*userGRPCContract.UpdateUserResponse), nil

}

// DeleteUser delete an existing user
// context: Mandatory. The reference to the context
// request: Mandatory. The request to delete an existing user
// Returns the result of deleting an existing user
func (service *transportService) DeleteUser(
	ctx context.Context,
	request *userGRPCContract.DeleteUserRequest) (*userGRPCContract.DeleteUserResponse, error) {
	_, response, err := service.deleteUserHandler.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.(*userGRPCContract.DeleteUserResponse), nil

}
