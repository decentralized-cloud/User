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
	commonErrors "github.com/micro-business/go-core/system/errors"
	"github.com/micro-business/gokit-core/middleware"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type transportService struct {
	logger                    *zap.Logger
	configurationService      configuration.ConfigurationContract
	endpointCreatorService    endpoint.EndpointCreatorContract
	middlewareProviderService middleware.MiddlewareProviderContract
	createUserHandler         gokitgrpc.Handler
	readUserHandler           gokitgrpc.Handler
	readUserByEmailHandler    gokitgrpc.Handler
	updateUserHandler         gokitgrpc.Handler
	deleteUserHandler         gokitgrpc.Handler
	searchHandler             gokitgrpc.Handler
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

	return &transportService{
		logger:                    logger,
		configurationService:      configurationService,
		endpointCreatorService:    endpointCreatorService,
		middlewareProviderService: middlewareProviderService,
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
	userGRPCContract.RegisterUserServiceServer(gRPCServer, service)
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
	service.createUserHandler = gokitgrpc.NewServer(
		endpoint,
		decodeCreateUserRequest,
		encodeCreateUserResponse,
	)

	endpoint = service.endpointCreatorService.ReadUserEndpoint()
	endpoint = service.middlewareProviderService.CreateLoggingMiddleware("ReadUser")(endpoint)
	service.readUserHandler = gokitgrpc.NewServer(
		endpoint,
		decodeReadUserRequest,
		encodeReadUserResponse,
	)

	endpoint = service.endpointCreatorService.ReadUserByEmailEndpoint()
	endpoint = service.middlewareProviderService.CreateLoggingMiddleware("ReadUserByEmail")(endpoint)
	service.readUserByEmailHandler = gokitgrpc.NewServer(
		endpoint,
		decodeReadUserByEmailRequest,
		encodeReadUserByEmailResponse,
	)

	endpoint = service.endpointCreatorService.UpdateUserEndpoint()
	endpoint = service.middlewareProviderService.CreateLoggingMiddleware("UpdateUser")(endpoint)
	service.updateUserHandler = gokitgrpc.NewServer(
		endpoint,
		decodeUpdateUserRequest,
		encodeUpdateUserResponse,
	)

	endpoint = service.endpointCreatorService.DeleteUserEndpoint()
	endpoint = service.middlewareProviderService.CreateLoggingMiddleware("DeleteUser")(endpoint)
	service.deleteUserHandler = gokitgrpc.NewServer(
		endpoint,
		decodeDeleteUserRequest,
		encodeDeleteUserResponse,
	)

	endpoint = service.endpointCreatorService.SearchEndpoint()
	endpoint = service.middlewareProviderService.CreateLoggingMiddleware("Search")(endpoint)
	service.searchHandler = gokitgrpc.NewServer(
		endpoint,
		decodeSearchRequest,
		encodeSearchResponse,
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

// ReadUserByEmail read an existing user by email address
// context: Mandatory. The reference to the context
// request: Mandatory. The request to read an existing user by email address
// Returns the result of reading an existing user by email address
func (service *transportService) ReadUserByEmail(
	ctx context.Context,
	request *userGRPCContract.ReadUserByEmailRequest) (*userGRPCContract.ReadUserByEmailResponse, error) {
	_, response, err := service.readUserByEmailHandler.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.(*userGRPCContract.ReadUserByEmailResponse), nil

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

// Search returns the list  of user that matched the provided criteria
// context: Mandatory. The reference to the context
// request: Mandatory. The request contains the filter criteria to look for existing user
// Returns the list of user that matched the provided criteria
func (service *transportService) Search(
	ctx context.Context,
	request *userGRPCContract.SearchRequest) (*userGRPCContract.SearchResponse, error) {
	_, response, err := service.searchHandler.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.(*userGRPCContract.SearchResponse), nil
}
