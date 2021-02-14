package business_test

import (
	"context"
	"errors"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/decentralized-cloud/user/models"
	"github.com/decentralized-cloud/user/services/business"
	repository "github.com/decentralized-cloud/user/services/repository"
	repsoitoryMock "github.com/decentralized-cloud/user/services/repository/mock"
	"github.com/golang/mock/gomock"
	"github.com/lucsky/cuid"
	commonErrors "github.com/micro-business/go-core/system/errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBusinessService(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())

	RegisterFailHandler(Fail)
	RunSpecs(t, "Business Service Tests")
}

var _ = Describe("Business Service Tests", func() {
	var (
		mockCtrl              *gomock.Controller
		sut                   business.BusinessContract
		mockRepositoryService *repsoitoryMock.MockRepositoryContract
		ctx                   context.Context
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())

		mockRepositoryService = repsoitoryMock.NewMockRepositoryContract(mockCtrl)
		sut, _ = business.NewBusinessService(mockRepositoryService)
		ctx = context.Background()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("user tries to instantiate BusinessService", func() {
		When("user repository service is not provided and NewBusinessService is called", func() {
			It("should return ArgumentNilError", func() {
				service, err := business.NewBusinessService(nil)
				Ω(service).Should(BeNil())
				assertArgumentNilError("repositoryService", "", err)
			})
		})

		When("all dependencies are resolved and NewBusinessService is called", func() {
			It("should instantiate the new BusinessService", func() {
				service, err := business.NewBusinessService(mockRepositoryService)
				Ω(err).Should(BeNil())
				Ω(service).ShouldNot(BeNil())
			})
		})
	})

	Describe("CreateUser", func() {
		var (
			request business.CreateUserRequest
		)

		BeforeEach(func() {
			request = business.CreateUserRequest{
				Email: cuid.New() + "@test.com",
				User:  models.User{}}
		})

		Context("user service is instantiated", func() {
			When("CreateUser is called", func() {
				It("should call user repository CreateUser method", func() {
					mockRepositoryService.
						EXPECT().
						CreateUser(ctx, gomock.Any()).
						Do(func(_ context.Context, mappedRequest *repository.CreateUserRequest) {
							Ω(mappedRequest.User).Should(Equal(request.User))
						}).
						Return(&repository.CreateUserResponse{}, nil)

					response, err := sut.CreateUser(ctx, &request)
					Ω(err).Should(BeNil())
					Ω(response.Err).Should(BeNil())
				})

				When("And user repository CreateUser return UserAlreadyExistError", func() {
					It("should return UserAlreadyExistsError", func() {
						expectedError := repository.NewUserAlreadyExistsError()
						mockRepositoryService.
							EXPECT().
							CreateUser(gomock.Any(), gomock.Any()).
							Return(nil, expectedError)

						response, err := sut.CreateUser(ctx, &request)
						Ω(err).Should(BeNil())
						assertUserAlreadyExistsError(response.Err, expectedError)
					})
				})

				When("And user repository CreateUser return any other error", func() {
					It("should return UnknownError", func() {
						expectedError := errors.New(cuid.New())
						mockRepositoryService.
							EXPECT().
							CreateUser(gomock.Any(), gomock.Any()).
							Return(nil, expectedError)

						response, err := sut.CreateUser(ctx, &request)
						Ω(err).Should(BeNil())
						assertUnknowError(expectedError.Error(), response.Err, expectedError)
					})
				})

				When("And user repository CreateUser return no error", func() {
					It("should return expected details", func() {
						expectedResponse := repository.CreateUserResponse{
							User:   models.User{},
							Cursor: cuid.New(),
						}

						mockRepositoryService.
							EXPECT().
							CreateUser(gomock.Any(), gomock.Any()).
							Return(&expectedResponse, nil)

						response, err := sut.CreateUser(ctx, &request)
						Ω(err).Should(BeNil())
						Ω(response.Err).Should(BeNil())
						assertUser(response.User, expectedResponse.User)
					})
				})
			})
		})
	})

	Describe("ReadUser", func() {
		var (
			request business.ReadUserRequest
		)

		BeforeEach(func() {
			request = business.ReadUserRequest{
				Email: cuid.New() + "@test.com",
			}
		})

		Context("user service is instantiated", func() {
			When("ReadUser is called", func() {
				It("should call user repository ReadUser method", func() {
					mockRepositoryService.
						EXPECT().
						ReadUser(ctx, gomock.Any()).
						Do(func(_ context.Context, mappedRequest *repository.ReadUserRequest) {
							Ω(mappedRequest.Email).Should(Equal(request.Email))
						}).
						Return(&repository.ReadUserResponse{}, nil)

					response, err := sut.ReadUser(ctx, &request)
					Ω(err).Should(BeNil())
					Ω(response.Err).Should(BeNil())
				})
			})

			When("And user repository ReadUser cannot find provided user", func() {
				It("should return UserNotFoundError", func() {
					expectedError := repository.NewUserNotFoundError(request.Email)
					mockRepositoryService.
						EXPECT().
						ReadUser(gomock.Any(), gomock.Any()).
						Return(nil, expectedError)

					response, err := sut.ReadUser(ctx, &request)
					Ω(err).Should(BeNil())
					assertUserNotFoundError(request.Email, response.Err, expectedError)
				})
			})

			When("And user repository ReadUser return any other error", func() {
				It("should return UnknownError", func() {
					expectedError := errors.New(cuid.New())
					mockRepositoryService.
						EXPECT().
						ReadUser(gomock.Any(), gomock.Any()).
						Return(nil, expectedError)

					response, err := sut.ReadUser(ctx, &request)
					Ω(err).Should(BeNil())
					assertUnknowError(expectedError.Error(), response.Err, expectedError)
				})
			})

			When("And user repository ReadUser return no error", func() {
				It("should return the user details", func() {
					expectedResponse := repository.ReadUserResponse{
						User: models.User{},
					}

					mockRepositoryService.
						EXPECT().
						ReadUser(gomock.Any(), gomock.Any()).
						Return(&expectedResponse, nil)

					response, err := sut.ReadUser(ctx, &request)
					Ω(err).Should(BeNil())
					Ω(response.Err).Should(BeNil())
					assertUser(response.User, expectedResponse.User)
				})
			})
		})
	})

	Describe("UpdateUser", func() {
		var (
			request business.UpdateUserRequest
		)

		BeforeEach(func() {
			request = business.UpdateUserRequest{
				Email: cuid.New() + "@test.com",
				User:  models.User{},
			}
		})

		Context("user service is instantiated", func() {
			When("UpdateUser is called", func() {
				It("should call user repository UpdateUser method", func() {
					mockRepositoryService.
						EXPECT().
						UpdateUser(ctx, gomock.Any()).
						Do(func(_ context.Context, mappedRequest *repository.UpdateUserRequest) {
							Ω(mappedRequest.Email).Should(Equal(request.Email))
						}).
						Return(&repository.UpdateUserResponse{}, nil)

					response, err := sut.UpdateUser(ctx, &request)
					Ω(err).Should(BeNil())
					Ω(response.Err).Should(BeNil())
				})
			})

			When("And user repository UpdateUser cannot find provided user", func() {
				It("should return UserNotFoundError", func() {
					expectedError := repository.NewUserNotFoundError(request.Email)
					mockRepositoryService.
						EXPECT().
						UpdateUser(gomock.Any(), gomock.Any()).
						Return(nil, expectedError)

					response, err := sut.UpdateUser(ctx, &request)
					Ω(err).Should(BeNil())
					assertUserNotFoundError(request.Email, response.Err, expectedError)
				})
			})

			When("And user repository UpdateUser return any other error", func() {
				It("should return UnknownError", func() {
					expectedError := errors.New(cuid.New())
					mockRepositoryService.
						EXPECT().
						UpdateUser(gomock.Any(), gomock.Any()).
						Return(nil, expectedError)

					response, err := sut.UpdateUser(ctx, &request)
					Ω(err).Should(BeNil())
					assertUnknowError(expectedError.Error(), response.Err, expectedError)
				})
			})

			When("And user repository UpdateUser return no error", func() {
				It("should return expected details", func() {
					expectedResponse := repository.UpdateUserResponse{
						User:   models.User{},
						Cursor: cuid.New(),
					}
					mockRepositoryService.
						EXPECT().
						UpdateUser(gomock.Any(), gomock.Any()).
						Return(&expectedResponse, nil)

					response, err := sut.UpdateUser(ctx, &request)
					Ω(err).Should(BeNil())
					Ω(response.Err).Should(BeNil())
					assertUser(response.User, expectedResponse.User)
				})
			})
		})
	})

	Describe("DeleteUser is called", func() {
		var (
			request business.DeleteUserRequest
		)

		BeforeEach(func() {
			request = business.DeleteUserRequest{
				Email: cuid.New() + "@test.com",
			}
		})

		Context("user service is instantiated", func() {
			When("DeleteUser is called", func() {
				It("should call user repository DeleteUser method", func() {
					mockRepositoryService.
						EXPECT().
						DeleteUser(ctx, gomock.Any()).
						Do(func(_ context.Context, mappedRequest *repository.DeleteUserRequest) {
							Ω(mappedRequest.Email).Should(Equal(request.Email))
						}).
						Return(&repository.DeleteUserResponse{}, nil)

					response, err := sut.DeleteUser(ctx, &request)
					Ω(err).Should(BeNil())
					Ω(response.Err).Should(BeNil())
				})
			})

			When("user repository DeleteUser cannot find provided user", func() {
				It("should return UserNotFoundError", func() {
					expectedError := repository.NewUserNotFoundError(request.Email)
					mockRepositoryService.
						EXPECT().
						DeleteUser(gomock.Any(), gomock.Any()).
						Return(nil, expectedError)

					response, err := sut.DeleteUser(ctx, &request)
					Ω(err).Should(BeNil())
					assertUserNotFoundError(request.Email, response.Err, expectedError)
				})
			})

			When("user repository DeleteUser is faced with any other error", func() {
				It("should return UnknownError", func() {
					expectedError := errors.New(cuid.New())
					mockRepositoryService.
						EXPECT().
						DeleteUser(gomock.Any(), gomock.Any()).
						Return(nil, expectedError)

					response, err := sut.DeleteUser(ctx, &request)
					Ω(err).Should(BeNil())
					assertUnknowError(expectedError.Error(), response.Err, expectedError)
				})
			})

			When("user repository DeleteUser completes successfully", func() {
				It("should return no error", func() {
					mockRepositoryService.
						EXPECT().
						DeleteUser(gomock.Any(), gomock.Any()).
						Return(&repository.DeleteUserResponse{}, nil)

					response, err := sut.DeleteUser(ctx, &request)
					Ω(err).Should(BeNil())
					Ω(response.Err).Should(BeNil())
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

func assertUnknowError(expectedMessage string, err error, nestedErr error) {
	Ω(business.IsUnknownError(err)).Should(BeTrue())

	var unknownErr business.UnknownError
	_ = errors.As(err, &unknownErr)

	Ω(strings.Contains(unknownErr.Error(), expectedMessage)).Should(BeTrue())
	Ω(errors.Unwrap(err)).Should(Equal(nestedErr))
}

func assertUserAlreadyExistsError(err error, nestedErr error) {
	Ω(business.IsUserAlreadyExistsError(err)).Should(BeTrue())
	Ω(errors.Unwrap(err)).Should(Equal(nestedErr))
}

func assertUserNotFoundError(expectedEmail string, err error, nestedErr error) {
	Ω(business.IsUserNotFoundError(err)).Should(BeTrue())

	var userNotFoundErr business.UserNotFoundError
	_ = errors.As(err, &userNotFoundErr)

	Ω(userNotFoundErr.Email).Should(Equal(expectedEmail))
	Ω(errors.Unwrap(err)).Should(Equal(nestedErr))
}

func assertUser(user, expectedUser models.User) {
	Ω(user).ShouldNot(BeNil())
}
