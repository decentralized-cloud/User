package endpoint_test

import (
	"context"
	"errors"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/decentralized-cloud/user/models"
	"github.com/decentralized-cloud/user/services/business"
	businessMock "github.com/decentralized-cloud/user/services/business/mock"
	"github.com/decentralized-cloud/user/services/endpoint"
	gokitendpoint "github.com/go-kit/kit/endpoint"
	"github.com/golang/mock/gomock"
	"github.com/lucsky/cuid"
	commonErrors "github.com/micro-business/go-core/system/errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEndpointCreatorService(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())

	RegisterFailHandler(Fail)
	RunSpecs(t, "Endpoint Creator Service Tests")
}

var _ = Describe("Endpoint Creator Service Tests", func() {
	var (
		mockCtrl            *gomock.Controller
		sut                 endpoint.EndpointCreatorContract
		mockBusinessService *businessMock.MockBusinessContract
		ctx                 context.Context
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())

		mockBusinessService = businessMock.NewMockBusinessContract(mockCtrl)
		sut, _ = endpoint.NewEndpointCreatorService(mockBusinessService)
		ctx = context.Background()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("user tries to instantiate EndpointCreatorService", func() {
		When("user business service is not provided and NewEndpointCreatorService is called", func() {
			It("should return ArgumentNilError", func() {
				service, err := endpoint.NewEndpointCreatorService(nil)
				Ω(service).Should(BeNil())
				assertArgumentNilError("businessService", "", err)
			})
		})

		When("all dependencies are resolved and NewEndpointCreatorService is called", func() {
			It("should instantiate the new EndpointCreatorService", func() {
				service, err := endpoint.NewEndpointCreatorService(mockBusinessService)
				Ω(err).Should(BeNil())
				Ω(service).ShouldNot(BeNil())
			})
		})
	})

	Context("EndpointCreatorService is instantiated", func() {
		When("CreateUserEndpoint is called", func() {
			It("should return valid function", func() {
				endpoint := sut.CreateUserEndpoint()
				Ω(endpoint).ShouldNot(BeNil())
			})

			var (
				endpoint gokitendpoint.Endpoint
				request  business.CreateUserRequest
				response business.CreateUserResponse
			)

			BeforeEach(func() {
				endpoint = sut.CreateUserEndpoint()
				request = business.CreateUserRequest{
					Email: cuid.New() + "@test.com",
					User:  models.User{},
				}

				response = business.CreateUserResponse{
					User:   models.User{},
					Cursor: cuid.New(),
				}
			})

			Context("CreateUserEndpoint function is returned", func() {
				When("endpoint is called with nil context", func() {
					It("should return ArgumentNilError", func() {
						returnedResponse, err := endpoint(nil, &request)

						Ω(err).Should(BeNil())
						Ω(returnedResponse).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.CreateUserResponse)
						assertArgumentNilError("ctx", "", castedResponse.Err)
					})
				})

				When("endpoint is called with nil request", func() {
					It("should return ArgumentNilError", func() {
						returnedResponse, err := endpoint(ctx, nil)

						Ω(err).Should(BeNil())
						Ω(returnedResponse).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.CreateUserResponse)
						assertArgumentNilError("request", "", castedResponse.Err)
					})
				})

				When("endpoint is called with invalid request", func() {
					It("should return ArgumentNilError", func() {
						invalidRequest := business.CreateUserRequest{
							Email: "",
							User:  models.User{}}
						returnedResponse, err := endpoint(ctx, &invalidRequest)

						Ω(err).Should(BeNil())
						Ω(response).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.CreateUserResponse)
						validationErr := invalidRequest.Validate()
						assertArgumentError("request", validationErr.Error(), castedResponse.Err, validationErr)
					})
				})

				When("endpoint is called with valid request", func() {
					It("should call business service CreateUser method", func() {
						mockBusinessService.
							EXPECT().
							CreateUser(ctx, gomock.Any()).
							Do(func(_ context.Context, mappedRequest *business.CreateUserRequest) {
								Ω(mappedRequest.User).Should(Equal(request.User))
							}).
							Return(&response, nil)

						returnedResponse, err := endpoint(ctx, &request)

						Ω(err).Should(BeNil())
						Ω(response).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.CreateUserResponse)
						Ω(castedResponse.Err).Should(BeNil())
					})
				})

				When("business service CreateUser returns error", func() {
					It("should return the same error", func() {
						expectedErr := errors.New(cuid.New())
						mockBusinessService.
							EXPECT().
							CreateUser(gomock.Any(), gomock.Any()).
							Return(nil, expectedErr)

						_, err := endpoint(ctx, &request)

						Ω(err).Should(Equal(expectedErr))
					})
				})

				When("business service CreateUser returns response", func() {
					It("should return the same response", func() {
						mockBusinessService.
							EXPECT().
							CreateUser(gomock.Any(), gomock.Any()).
							Return(&response, nil)

						returnedResponse, err := endpoint(ctx, &request)

						Ω(err).Should(BeNil())
						Ω(returnedResponse).Should(Equal(&response))
					})
				})
			})
		})
	})

	Context("EndpointCreatorService is instantiated", func() {
		When("ReadUserEndpoint is called", func() {
			It("should return valid function", func() {
				endpoint := sut.ReadUserEndpoint()
				Ω(endpoint).ShouldNot(BeNil())
			})

			var (
				endpoint gokitendpoint.Endpoint
				request  business.ReadUserRequest
				response business.ReadUserResponse
			)

			BeforeEach(func() {
				endpoint = sut.ReadUserEndpoint()
				request = business.ReadUserRequest{
					Email: cuid.New() + "@test.com",
				}

				response = business.ReadUserResponse{
					User: models.User{},
				}
			})

			Context("ReadUserEndpoint function is returned", func() {
				When("endpoint is called with nil context", func() {
					It("should return ArgumentNilError", func() {
						returnedResponse, err := endpoint(nil, &request)

						Ω(err).Should(BeNil())
						Ω(returnedResponse).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.ReadUserResponse)
						assertArgumentNilError("ctx", "", castedResponse.Err)
					})
				})

				When("endpoint is called with nil request", func() {
					It("should return ArgumentNilError", func() {
						returnedResponse, err := endpoint(ctx, nil)

						Ω(err).Should(BeNil())
						Ω(returnedResponse).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.ReadUserResponse)
						assertArgumentNilError("request", "", castedResponse.Err)
					})
				})

				When("endpoint is called with invalid request", func() {
					It("should return ArgumentNilError", func() {
						invalidRequest := business.ReadUserRequest{
							Email: "",
						}
						returnedResponse, err := endpoint(ctx, &invalidRequest)

						Ω(err).Should(BeNil())
						Ω(response).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.ReadUserResponse)
						validationErr := invalidRequest.Validate()
						assertArgumentError("request", validationErr.Error(), castedResponse.Err, validationErr)
					})
				})

				When("endpoint is called with valid request", func() {
					It("should call business service ReadUser method", func() {
						mockBusinessService.
							EXPECT().
							ReadUser(ctx, gomock.Any()).
							Do(func(_ context.Context, mappedRequest *business.ReadUserRequest) {
								Ω(mappedRequest.Email).Should(Equal(request.Email))
							}).
							Return(&response, nil)

						returnedResponse, err := endpoint(ctx, &request)

						Ω(err).Should(BeNil())
						Ω(response).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.ReadUserResponse)
						Ω(castedResponse.Err).Should(BeNil())
					})
				})

				When("business service ReadUser returns error", func() {
					It("should return the same error", func() {
						expectedErr := errors.New(cuid.New())
						mockBusinessService.
							EXPECT().
							ReadUser(gomock.Any(), gomock.Any()).
							Return(nil, expectedErr)

						_, err := endpoint(ctx, &request)

						Ω(err).Should(Equal(expectedErr))
					})
				})

				When("business service ReadUser returns response", func() {
					It("should return the same response", func() {
						mockBusinessService.
							EXPECT().
							ReadUser(gomock.Any(), gomock.Any()).
							Return(&response, nil)

						returnedResponse, err := endpoint(ctx, &request)

						Ω(err).Should(BeNil())
						Ω(returnedResponse).Should(Equal(&response))
					})
				})
			})
		})
	})

	Context("EndpointCreatorService is instantiated", func() {
		When("UpdateUserEndpoint is called", func() {
			It("should return valid function", func() {
				endpoint := sut.UpdateUserEndpoint()
				Ω(endpoint).ShouldNot(BeNil())
			})

			var (
				endpoint gokitendpoint.Endpoint
				request  business.UpdateUserRequest
				response business.UpdateUserResponse
			)

			BeforeEach(func() {
				endpoint = sut.UpdateUserEndpoint()
				request = business.UpdateUserRequest{
					Email: cuid.New() + "@test.com",
					User:  models.User{}}

				response = business.UpdateUserResponse{
					User:   models.User{},
					Cursor: cuid.New(),
				}
			})

			Context("UpdateUserEndpoint function is returned", func() {
				When("endpoint is called with nil context", func() {
					It("should return ArgumentNilError", func() {
						returnedResponse, err := endpoint(nil, &request)

						Ω(err).Should(BeNil())
						Ω(returnedResponse).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.UpdateUserResponse)
						assertArgumentNilError("ctx", "", castedResponse.Err)
					})
				})

				When("endpoint is called with nil request", func() {
					It("should return ArgumentNilError", func() {
						returnedResponse, err := endpoint(ctx, nil)

						Ω(err).Should(BeNil())
						Ω(returnedResponse).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.UpdateUserResponse)
						assertArgumentNilError("request", "", castedResponse.Err)
					})
				})

				When("endpoint is called with invalid request", func() {
					It("should return ArgumentNilError", func() {
						invalidRequest := business.UpdateUserRequest{
							Email: "",
							User:  models.User{}}
						returnedResponse, err := endpoint(ctx, &invalidRequest)

						Ω(err).Should(BeNil())
						Ω(response).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.UpdateUserResponse)
						validationErr := invalidRequest.Validate()
						assertArgumentError("request", validationErr.Error(), castedResponse.Err, validationErr)
					})
				})

				When("endpoint is called with valid request", func() {
					It("should call business service UpdateUser method", func() {
						mockBusinessService.
							EXPECT().
							UpdateUser(ctx, gomock.Any()).
							Do(func(_ context.Context, mappedRequest *business.UpdateUserRequest) {
								Ω(mappedRequest.Email).Should(Equal(request.Email))
							}).
							Return(&response, nil)

						returnedResponse, err := endpoint(ctx, &request)

						Ω(err).Should(BeNil())
						Ω(response).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.UpdateUserResponse)
						Ω(castedResponse.Err).Should(BeNil())
					})
				})

				When("business service UpdateUser returns error", func() {
					It("should return the same error", func() {
						expectedErr := errors.New(cuid.New())
						mockBusinessService.
							EXPECT().
							UpdateUser(gomock.Any(), gomock.Any()).
							Return(nil, expectedErr)

						_, err := endpoint(ctx, &request)

						Ω(err).Should(Equal(expectedErr))
					})
				})

				When("business service UpdateUser returns response", func() {
					It("should return the same response", func() {
						mockBusinessService.
							EXPECT().
							UpdateUser(gomock.Any(), gomock.Any()).
							Return(&response, nil)

						returnedResponse, err := endpoint(ctx, &request)

						Ω(err).Should(BeNil())
						Ω(returnedResponse).Should(Equal(&response))
					})
				})
			})
		})
	})

	Context("EndpointCreatorService is instantiated", func() {
		When("DeleteUserEndpoint is called", func() {
			It("should return valid function", func() {
				endpoint := sut.DeleteUserEndpoint()
				Ω(endpoint).ShouldNot(BeNil())
			})

			var (
				endpoint gokitendpoint.Endpoint
				request  business.DeleteUserRequest
				response business.DeleteUserResponse
			)

			BeforeEach(func() {
				endpoint = sut.DeleteUserEndpoint()
				request = business.DeleteUserRequest{
					Email: cuid.New() + "@test.com",
				}

				response = business.DeleteUserResponse{}
			})

			Context("DeleteUserEndpoint function is returned", func() {
				When("endpoint is called with nil context", func() {
					It("should return ArgumentNilError", func() {
						returnedResponse, err := endpoint(nil, &request)

						Ω(err).Should(BeNil())
						Ω(returnedResponse).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.DeleteUserResponse)
						assertArgumentNilError("ctx", "", castedResponse.Err)
					})
				})

				When("endpoint is called with nil request", func() {
					It("should return ArgumentNilError", func() {
						returnedResponse, err := endpoint(ctx, nil)

						Ω(err).Should(BeNil())
						Ω(returnedResponse).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.DeleteUserResponse)
						assertArgumentNilError("request", "", castedResponse.Err)
					})
				})

				When("endpoint is called with invalid request", func() {
					It("should return ArgumentNilError", func() {
						invalidRequest := business.DeleteUserRequest{
							Email: "",
						}
						returnedResponse, err := endpoint(ctx, &invalidRequest)

						Ω(err).Should(BeNil())
						Ω(response).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.DeleteUserResponse)
						validationErr := invalidRequest.Validate()
						assertArgumentError("request", validationErr.Error(), castedResponse.Err, validationErr)
					})
				})

				When("endpoint is called with valid request", func() {
					It("should call business service DeleteUser method", func() {
						mockBusinessService.
							EXPECT().
							DeleteUser(ctx, gomock.Any()).
							Do(func(_ context.Context, mappedRequest *business.DeleteUserRequest) {
								Ω(mappedRequest.Email).Should(Equal(request.Email))
							}).
							Return(&response, nil)

						returnedResponse, err := endpoint(ctx, &request)

						Ω(err).Should(BeNil())
						Ω(response).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.DeleteUserResponse)
						Ω(castedResponse.Err).Should(BeNil())
					})
				})

				When("business service DeleteUser returns error", func() {
					It("should return the same error", func() {
						expectedErr := errors.New(cuid.New())
						mockBusinessService.
							EXPECT().
							DeleteUser(gomock.Any(), gomock.Any()).
							Return(nil, expectedErr)

						_, err := endpoint(ctx, &request)

						Ω(err).Should(Equal(expectedErr))
					})
				})

				When("business service DeleteUser returns response", func() {
					It("should return the same response", func() {
						mockBusinessService.
							EXPECT().
							DeleteUser(gomock.Any(), gomock.Any()).
							Return(&response, nil)

						returnedResponse, err := endpoint(ctx, &request)

						Ω(err).Should(BeNil())
						Ω(returnedResponse).Should(Equal(&response))
					})
				})
			})
		})
	})

})

func assertArgumentNilError(expectedArgumentName, expectedMessage string, err error) {
	Ω(commonErrors.IsArgumentNilError(err)).Should(BeTrue())

	var argumentNilErr commonErrors.ArgumentNilError
	_ = errors.As(err, &argumentNilErr)

	if expectedArgumentName != "" {
		Ω(argumentNilErr.ArgumentName).Should(Equal(expectedArgumentName))
	}

	if expectedMessage != "" {
		Ω(strings.Contains(argumentNilErr.Error(), expectedMessage)).Should(BeTrue())
	}
}

func assertArgumentError(expectedArgumentName, expectedMessage string, err error, nestedErr error) {
	Ω(commonErrors.IsArgumentError(err)).Should(BeTrue())

	var argumentErr commonErrors.ArgumentError
	_ = errors.As(err, &argumentErr)

	Ω(argumentErr.ArgumentName).Should(Equal(expectedArgumentName))
	Ω(strings.Contains(argumentErr.Error(), expectedMessage)).Should(BeTrue())
	Ω(errors.Unwrap(err)).Should(Equal(nestedErr))
}
