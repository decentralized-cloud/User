// Package util implements different utilities required by the user service
package util

import (
	"log"
	"os"
	"os/signal"

	"github.com/decentralized-cloud/user/services/business"
	"github.com/decentralized-cloud/user/services/configuration"
	"github.com/decentralized-cloud/user/services/endpoint"
	"github.com/decentralized-cloud/user/services/repository/mongodb"
	"github.com/decentralized-cloud/user/services/transport/grpc"
	"github.com/decentralized-cloud/user/services/transport/https"
	"github.com/micro-business/go-core/gokit/middleware"
	"go.uber.org/zap"
)

var configurationService configuration.ConfigurationContract
var endpointCreatorService endpoint.EndpointCreatorContract
var middlewareProviderService middleware.MiddlewareProviderContract

// StartService setups all dependecies required to start the user service and
// start the service
func StartService() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		_ = logger.Sync()
	}()

	if err = setupDependencies(logger); err != nil {
		logger.Fatal("failed to setup dependecies", zap.Error(err))
	}

	grpcTransportService, err := grpc.NewTransportService(
		logger,
		configurationService,
		endpointCreatorService,
		middlewareProviderService)
	if err != nil {
		logger.Fatal("failed to create gRPC transport service", zap.Error(err))
	}

	httpsTansportService, err := https.NewTransportService(
		logger,
		configurationService)
	if err != nil {
		logger.Fatal("failed to create HTTPS transport service", zap.Error(err))
	}

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan struct{})
	signal.Notify(signalChan, os.Interrupt)

	go func() {
		if serviceErr := grpcTransportService.Start(); serviceErr != nil {
			logger.Fatal("failed to start gRPC transport service", zap.Error(serviceErr))
		}
	}()

	go func() {
		if serviceErr := httpsTansportService.Start(); serviceErr != nil {
			logger.Fatal("failed to start HTTPS transport service", zap.Error(serviceErr))
		}
	}()

	go func() {
		<-signalChan
		logger.Info("Received an interrupt, stopping services...")

		if err := grpcTransportService.Stop(); err != nil {
			logger.Error("failed to stop gRPC transport service", zap.Error(err))
		}

		if err := httpsTansportService.Stop(); err != nil {
			logger.Error("failed to stop HTTPS transport service", zap.Error(err))
		}

		close(cleanupDone)
	}()
	<-cleanupDone
}

func setupDependencies(logger *zap.Logger) (err error) {
	if configurationService, err = configuration.NewEnvConfigurationService(); err != nil {
		return
	}

	if middlewareProviderService, err = middleware.NewMiddlewareProviderService(logger, true, ""); err != nil {
		return
	}

	repositoryService, err := mongodb.NewMongodbRepositoryService(configurationService)
	if err != nil {
		return
	}

	businessService, err := business.NewBusinessService(repositoryService)
	if err != nil {
		return err
	}

	if endpointCreatorService, err = endpoint.NewEndpointCreatorService(businessService); err != nil {
		return
	}

	return
}
