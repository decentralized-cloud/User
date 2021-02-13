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
	"github.com/micro-business/go-core/common"
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
				User: models.User{
					Email: cuid.New() + "@test.com",
				}}
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
							UserID: cuid.New(),
							User: models.User{
								Email: cuid.New() + "@test.com",
							},
							Cursor: cuid.New(),
						}

						mockRepositoryService.
							EXPECT().
							CreateUser(gomock.Any(), gomock.Any()).
							Return(&expectedResponse, nil)

						response, err := sut.CreateUser(ctx, &request)
						Ω(err).Should(BeNil())
						Ω(response.Err).Should(BeNil())
						Ω(response.UserID).ShouldNot(BeNil())
						Ω(response.UserID).Should(Equal(expectedResponse.UserID))
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
				UserID: cuid.New(),
			}
		})

		Context("user service is instantiated", func() {
			When("ReadUser is called", func() {
				It("should call user repository ReadUser method", func() {
					mockRepositoryService.
						EXPECT().
						ReadUser(ctx, gomock.Any()).
						Do(func(_ context.Context, mappedRequest *repository.ReadUserRequest) {
							Ω(mappedRequest.UserID).Should(Equal(request.UserID))
						}).
						Return(&repository.ReadUserResponse{}, nil)

					response, err := sut.ReadUser(ctx, &request)
					Ω(err).Should(BeNil())
					Ω(response.Err).Should(BeNil())
				})
			})

			When("And user repository ReadUser cannot find provided user", func() {
				It("should return UserNotFoundError", func() {
					expectedError := repository.NewUserNotFoundError(request.UserID)
					mockRepositoryService.
						EXPECT().
						ReadUser(gomock.Any(), gomock.Any()).
						Return(nil, expectedError)

					response, err := sut.ReadUser(ctx, &request)
					Ω(err).Should(BeNil())
					assertUserNotFoundError(request.UserID, response.Err, expectedError)
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
						User: models.User{Email: cuid.New() + "@test.com"},
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
				UserID: cuid.New(),
				User:   models.User{Email: cuid.New() + "@test.com"},
			}
		})

		Context("user service is instantiated", func() {
			When("UpdateUser is called", func() {
				It("should call user repository UpdateUser method", func() {
					mockRepositoryService.
						EXPECT().
						UpdateUser(ctx, gomock.Any()).
						Do(func(_ context.Context, mappedRequest *repository.UpdateUserRequest) {
							Ω(mappedRequest.UserID).Should(Equal(request.UserID))
							Ω(mappedRequest.User.Email).Should(Equal(request.User.Email))
						}).
						Return(&repository.UpdateUserResponse{}, nil)

					response, err := sut.UpdateUser(ctx, &request)
					Ω(err).Should(BeNil())
					Ω(response.Err).Should(BeNil())
				})
			})

			When("And user repository UpdateUser cannot find provided user", func() {
				It("should return UserNotFoundError", func() {
					expectedError := repository.NewUserNotFoundError(request.UserID)
					mockRepositoryService.
						EXPECT().
						UpdateUser(gomock.Any(), gomock.Any()).
						Return(nil, expectedError)

					response, err := sut.UpdateUser(ctx, &request)
					Ω(err).Should(BeNil())
					assertUserNotFoundError(request.UserID, response.Err, expectedError)
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
						User: models.User{
							Email: cuid.New() + "@test.com",
						},
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
				UserID: cuid.New(),
			}
		})

		Context("user service is instantiated", func() {
			When("DeleteUser is called", func() {
				It("should call user repository DeleteUser method", func() {
					mockRepositoryService.
						EXPECT().
						DeleteUser(ctx, gomock.Any()).
						Do(func(_ context.Context, mappedRequest *repository.DeleteUserRequest) {
							Ω(mappedRequest.UserID).Should(Equal(request.UserID))
						}).
						Return(&repository.DeleteUserResponse{}, nil)

					response, err := sut.DeleteUser(ctx, &request)
					Ω(err).Should(BeNil())
					Ω(response.Err).Should(BeNil())
				})
			})

			When("user repository DeleteUser cannot find provided user", func() {
				It("should return UserNotFoundError", func() {
					expectedError := repository.NewUserNotFoundError(request.UserID)
					mockRepositoryService.
						EXPECT().
						DeleteUser(gomock.Any(), gomock.Any()).
						Return(nil, expectedError)

					response, err := sut.DeleteUser(ctx, &request)
					Ω(err).Should(BeNil())
					assertUserNotFoundError(request.UserID, response.Err, expectedError)
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

	Describe("Search is called", func() {
		var (
			request business.SearchRequest
			userIDs []string
		)

		BeforeEach(func() {
			userIDs = []string{}
			for idx := 0; idx < rand.Intn(20)+1; idx++ {
				userIDs = append(userIDs, cuid.New())
			}

			request = business.SearchRequest{
				Pagination: common.Pagination{
					After:  convertStringToPointer(cuid.New()),
					First:  convertIntToPointer(rand.Intn(1000)),
					Before: convertStringToPointer(cuid.New()),
					Last:   convertIntToPointer(rand.Intn(1000)),
				},
				SortingOptions: []common.SortingOptionPair{
					common.SortingOptionPair{
						Name:      cuid.New(),
						Direction: common.Ascending,
					},
					common.SortingOptionPair{
						Name:      cuid.New(),
						Direction: common.Descending,
					},
				},
				UserIDs: userIDs,
			}
		})

		Context("user service is instantiated", func() {
			When("Search is called", func() {
				It("should call user repository Search method", func() {
					mockRepositoryService.
						EXPECT().
						Search(ctx, gomock.Any()).
						Do(func(_ context.Context, mappedRequest *repository.SearchRequest) {
							Ω(mappedRequest.Pagination).Should(Equal(request.Pagination))
							Ω(mappedRequest.SortingOptions).Should(Equal(request.SortingOptions))
							Ω(mappedRequest.UserIDs).Should(Equal(request.UserIDs))
						}).
						Return(&repository.SearchResponse{}, nil)

					response, err := sut.Search(ctx, &request)
					Ω(err).Should(BeNil())
					Ω(response.Err).Should(BeNil())
				})
			})

			When("user repository Search is faced with any other error", func() {
				It("should return UnknownError", func() {
					expectedError := errors.New(cuid.New())
					mockRepositoryService.
						EXPECT().
						Search(gomock.Any(), gomock.Any()).
						Return(nil, expectedError)

					response, err := sut.Search(ctx, &request)
					Ω(err).Should(BeNil())
					assertUnknowError(expectedError.Error(), response.Err, expectedError)
				})
			})

			When("user repository Search completes successfully", func() {
				It("should return the list of matched userIDs", func() {
					users := []models.UserWithCursor{}

					for idx := 0; idx < rand.Intn(20)+1; idx++ {
						users = append(users, models.UserWithCursor{
							UserID: cuid.New(),
							User: models.User{
								Email: cuid.New() + "@test.com",
							},
							Cursor: cuid.New(),
						})
					}

					expectedResponse := repository.SearchResponse{
						HasPreviousPage: (rand.Intn(10) % 2) == 0,
						HasNextPage:     (rand.Intn(10) % 2) == 0,
						TotalCount:      rand.Int63n(1000),
						Users:           users,
					}

					mockRepositoryService.
						EXPECT().
						Search(gomock.Any(), gomock.Any()).
						Return(&expectedResponse, nil)

					response, err := sut.Search(ctx, &request)
					Ω(err).Should(BeNil())
					Ω(response.Err).Should(BeNil())
					Ω(response.HasPreviousPage).Should(Equal(expectedResponse.HasPreviousPage))
					Ω(response.HasNextPage).Should(Equal(expectedResponse.HasNextPage))
					Ω(response.TotalCount).Should(Equal(expectedResponse.TotalCount))
					Ω(response.Users).Should(Equal(expectedResponse.Users))
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

func assertUserNotFoundError(expectedUserID string, err error, nestedErr error) {
	Ω(business.IsUserNotFoundError(err)).Should(BeTrue())

	var userNotFoundErr business.UserNotFoundError
	_ = errors.As(err, &userNotFoundErr)

	Ω(userNotFoundErr.UserID).Should(Equal(expectedUserID))
	Ω(errors.Unwrap(err)).Should(Equal(nestedErr))
}

func assertUser(user, expectedUser models.User) {
	Ω(user).ShouldNot(BeNil())
	Ω(user.Email).Should(Equal(expectedUser.Email))
}

func convertStringToPointer(str string) *string {
	return &str
}

func convertIntToPointer(i int) *int {
	return &i
}
