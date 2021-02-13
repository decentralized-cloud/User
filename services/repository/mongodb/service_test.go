package mongodb_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/decentralized-cloud/user/models"
	configurationMock "github.com/decentralized-cloud/user/services/configuration/mock"
	"github.com/decentralized-cloud/user/services/repository"
	"github.com/decentralized-cloud/user/services/repository/mongodb"
	"github.com/golang/mock/gomock"
	"github.com/lucsky/cuid"
	"github.com/micro-business/go-core/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMongodbRepositoryService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mongodb Repository Service Tests")
}

var _ = Describe("Mongodb Repository Service Tests", func() {
	var (
		mockCtrl      *gomock.Controller
		sut           repository.RepositoryContract
		ctx           context.Context
		createRequest repository.CreateUserRequest
	)

	BeforeEach(func() {
		connectionString := os.Getenv("DATABASE_CONNECTION_STRING")
		if strings.Trim(connectionString, " ") == "" {
			connectionString = "mongodb://mongodb:27017"
		}

		mockCtrl = gomock.NewController(GinkgoT())
		mockConfigurationService := configurationMock.NewMockConfigurationContract(mockCtrl)
		mockConfigurationService.
			EXPECT().
			GetDatabaseConnectionString().
			Return(connectionString, nil)

		mockConfigurationService.
			EXPECT().
			GetDatabaseName().
			Return("user", nil)

		mockConfigurationService.
			EXPECT().
			GetDatabaseCollectionName().
			Return("user", nil)

		sut, _ = mongodb.NewMongodbRepositoryService(mockConfigurationService)
		ctx = context.Background()
		createRequest = repository.CreateUserRequest{
			User: models.User{
				Email: cuid.New() + "@test.com",
			}}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("user tries to instantiate RepositoryService", func() {
		When("all dependencies are resolved and NewRepositoryService is called", func() {
			It("should instantiate the new RepositoryService", func() {
				mockConfigurationService := configurationMock.NewMockConfigurationContract(mockCtrl)
				mockConfigurationService.
					EXPECT().
					GetDatabaseConnectionString().
					Return(cuid.New(), nil)

				mockConfigurationService.
					EXPECT().
					GetDatabaseName().
					Return(cuid.New(), nil)

				mockConfigurationService.
					EXPECT().
					GetDatabaseCollectionName().
					Return(cuid.New(), nil)

				service, err := mongodb.NewMongodbRepositoryService(mockConfigurationService)
				Ω(err).Should(BeNil())
				Ω(service).ShouldNot(BeNil())
			})
		})
	})

	Context("user going to create a new user", func() {
		When("create user is called", func() {
			It("should create the new user", func() {
				response, err := sut.CreateUser(ctx, &createRequest)
				Ω(err).Should(BeNil())
				Ω(response.UserID).ShouldNot(BeNil())
				Ω(response.Cursor).Should(Equal(response.UserID))
				assertUser(response.User, createRequest.User)
			})
		})
	})

	Context("user already exists", func() {
		var (
			userID string
		)

		BeforeEach(func() {
			response, _ := sut.CreateUser(ctx, &createRequest)
			userID = response.UserID
		})

		When("user reads a user by Id", func() {
			It("should return a user", func() {
				response, err := sut.ReadUser(ctx, &repository.ReadUserRequest{UserID: userID})
				Ω(err).Should(BeNil())
				assertUser(response.User, createRequest.User)
			})
		})

		When("user updates the existing user", func() {
			It("should update the user information", func() {
				updateRequest := repository.UpdateUserRequest{
					UserID: userID,
					User: models.User{
						Email: cuid.New() + "@test.com",
					}}

				updateResponse, err := sut.UpdateUser(ctx, &updateRequest)
				Ω(err).Should(BeNil())
				Ω(updateResponse.Cursor).Should(Equal(userID))
				assertUser(updateResponse.User, updateRequest.User)

				readResponse, err := sut.ReadUser(ctx, &repository.ReadUserRequest{UserID: userID})
				Ω(err).Should(BeNil())
				assertUser(readResponse.User, updateRequest.User)
			})
		})

		When("user deletes the user", func() {
			It("should delete the user", func() {
				_, err := sut.DeleteUser(ctx, &repository.DeleteUserRequest{UserID: userID})
				Ω(err).Should(BeNil())

				response, err := sut.ReadUser(ctx, &repository.ReadUserRequest{UserID: userID})
				Ω(err).Should(HaveOccurred())
				Ω(response).Should(BeNil())

				Ω(repository.IsUserNotFoundError(err)).Should(BeTrue())

				var notFoundErr repository.UserNotFoundError
				_ = errors.As(err, &notFoundErr)

				Ω(notFoundErr.UserID).Should(Equal(userID))
			})
		})
	})

	Context("user does not exist", func() {
		var (
			userID string
		)

		BeforeEach(func() {
			userID = cuid.New()
		})

		When("user reads the user", func() {
			It("should return NotFoundError", func() {
				response, err := sut.ReadUser(ctx, &repository.ReadUserRequest{UserID: userID})
				Ω(err).Should(HaveOccurred())
				Ω(response).Should(BeNil())

				Ω(repository.IsUserNotFoundError(err)).Should(BeTrue())

				var notFoundErr repository.UserNotFoundError
				_ = errors.As(err, &notFoundErr)

				Ω(notFoundErr.UserID).Should(Equal(userID))
			})
		})

		When("user tries to update the user", func() {
			It("should return NotFoundError", func() {
				updateRequest := repository.UpdateUserRequest{
					UserID: userID,
					User: models.User{
						Email: cuid.New() + "@test.com",
					}}

				response, err := sut.UpdateUser(ctx, &updateRequest)
				Ω(err).Should(HaveOccurred())
				Ω(response).Should(BeNil())

				Ω(repository.IsUserNotFoundError(err)).Should(BeTrue())

				var notFoundErr repository.UserNotFoundError
				_ = errors.As(err, &notFoundErr)

				Ω(notFoundErr.UserID).Should(Equal(userID))
			})
		})

		When("user tries to delete the user", func() {
			It("should return NotFoundError", func() {
				response, err := sut.DeleteUser(ctx, &repository.DeleteUserRequest{UserID: userID})
				Ω(err).Should(HaveOccurred())
				Ω(response).Should(BeNil())

				Ω(repository.IsUserNotFoundError(err)).Should(BeTrue())

				var notFoundErr repository.UserNotFoundError
				_ = errors.As(err, &notFoundErr)

				Ω(notFoundErr.UserID).Should(Equal(userID))
			})
		})
	})

	Context("user already exists", func() {
		var (
			userIDs []string
		)

		BeforeEach(func() {
			userIDs = []string{}

			for i := 0; i < 10; i++ {
				email := fmt.Sprintf("Name%d@test.com", i)
				createRequest.User.Email = email
				response, _ := sut.CreateUser(ctx, &createRequest)
				userIDs = append(userIDs, response.UserID)
			}
		})

		When("user searches for users with selected user Ids and first 10 users", func() {
			It("should return first 10 users", func() {
				first := 10
				searchRequest := repository.SearchRequest{
					UserIDs: userIDs,
					Pagination: common.Pagination{
						After: nil,
						First: &first,
					},
					SortingOptions: []common.SortingOptionPair{},
				}

				response, err := sut.Search(ctx, &searchRequest)
				Ω(err).Should(BeNil())
				Ω(response.Users).ShouldNot(BeNil())
				Ω(len(response.Users)).Should(Equal(10))
				for i := 0; i < 10; i++ {
					Ω(response.Users[i].UserID).Should(Equal(userIDs[i]))
					email := fmt.Sprintf("Name%d@test.com", i)
					Ω(response.Users[i].User.Email).Should(Equal(email))
				}
			})
		})

		When("user searches for users with selected user Ids and first 5 users", func() {
			It("should return first 5 users", func() {
				first := 5
				searchRequest := repository.SearchRequest{
					UserIDs: userIDs,
					Pagination: common.Pagination{
						After: nil,
						First: &first,
					},
					SortingOptions: []common.SortingOptionPair{},
				}

				response, err := sut.Search(ctx, &searchRequest)
				Ω(err).Should(BeNil())
				Ω(response.Users).ShouldNot(BeNil())
				Ω(len(response.Users)).Should(Equal(5))
				for i := 0; i < 5; i++ {
					Ω(response.Users[i].UserID).Should(Equal(userIDs[i]))
					email := fmt.Sprintf("Name%d@test.com", i)
					Ω(response.Users[i].User.Email).Should(Equal(email))
				}
			})
		})

		When("user searches for users with selected user Ids with After parameter provided.", func() {
			It("should return first 9 users after provided user id", func() {
				first := 9
				searchRequest := repository.SearchRequest{
					UserIDs: userIDs,
					Pagination: common.Pagination{
						After: &userIDs[0],
						First: &first,
					},
					SortingOptions: []common.SortingOptionPair{},
				}

				response, err := sut.Search(ctx, &searchRequest)
				Ω(err).Should(BeNil())
				Ω(response.Users).ShouldNot(BeNil())
				Ω(len(response.Users)).Should(Equal(9))
				for i := 1; i < 10; i++ {
					Ω(response.Users[i-1].UserID).Should(Equal(userIDs[i]))
					email := fmt.Sprintf("Name%d@test.com", i)
					Ω(response.Users[i-1].User.Email).Should(Equal(email))
				}
			})
		})

		When("user searches for users with selected user Ids with After parameter provided.", func() {
			It("should return first 5 users after provided user id", func() {
				first := 5
				searchRequest := repository.SearchRequest{
					UserIDs: userIDs,
					Pagination: common.Pagination{
						After: &userIDs[0],
						First: &first,
					},
					SortingOptions: []common.SortingOptionPair{},
				}

				response, err := sut.Search(ctx, &searchRequest)
				Ω(err).Should(BeNil())
				Ω(response.Users).ShouldNot(BeNil())
				Ω(len(response.Users)).Should(Equal(5))
				for i := 1; i < 5; i++ {
					Ω(response.Users[i-1].UserID).Should(Equal(userIDs[i]))
					email := fmt.Sprintf("Name%d@test.com", i)
					Ω(response.Users[i-1].User.Email).Should(Equal(email))
				}
			})
		})

		//TODO : this test does not make sense
		When("user searches for users with selected user Ids and last 10 users", func() {
			It("should return first 10 users", func() {
				last := 10
				searchRequest := repository.SearchRequest{
					UserIDs: userIDs,
					Pagination: common.Pagination{
						Before: nil,
						Last:   &last,
					},
					SortingOptions: []common.SortingOptionPair{},
				}

				response, err := sut.Search(ctx, &searchRequest)
				Ω(err).Should(BeNil())
				Ω(response.Users).ShouldNot(BeNil())
				Ω(len(response.Users)).Should(Equal(10))
				for i := 0; i < 10; i++ {
					Ω(response.Users[i].UserID).Should(Equal(userIDs[i]))
					userName := fmt.Sprintf("Name%d@test.com", i)
					Ω(response.Users[i].User.Email).Should(Equal(userName))
				}
			})
		})

		When("user searches for users with selected user Ids with Before parameter provided.", func() {
			It("should return first 9 users before provided user id", func() {
				last := 9
				searchRequest := repository.SearchRequest{
					UserIDs: userIDs,
					Pagination: common.Pagination{
						Before: &userIDs[9],
						Last:   &last,
					},
					SortingOptions: []common.SortingOptionPair{},
				}

				response, err := sut.Search(ctx, &searchRequest)
				Ω(err).Should(BeNil())
				Ω(response.Users).ShouldNot(BeNil())
				Ω(len(response.Users)).Should(Equal(9))
				for i := 0; i < 9; i++ {
					Ω(response.Users[i].UserID).Should(Equal(userIDs[i]))
					email := fmt.Sprintf("Name%d@test.com", i)
					Ω(response.Users[i].User.Email).Should(Equal(email))
				}
			})
		})

		When("user searches for users with selected user Ids and first 10 users with ascending order on name property", func() {
			It("should return first 10 users in adcending order on name field", func() {
				first := 10
				searchRequest := repository.SearchRequest{
					UserIDs: userIDs,
					Pagination: common.Pagination{
						After: nil,
						First: &first,
					},
					SortingOptions: []common.SortingOptionPair{
						{Name: "email", Direction: common.Ascending},
					},
				}

				response, err := sut.Search(ctx, &searchRequest)
				Ω(err).Should(BeNil())
				Ω(response.Users).ShouldNot(BeNil())
				Ω(len(response.Users)).Should(Equal(10))
				for i := 0; i < 10; i++ {
					Ω(response.Users[i].UserID).Should(Equal(userIDs[i]))
					email := fmt.Sprintf("Name%d@test.com", i)
					Ω(response.Users[i].User.Email).Should(Equal(email))
				}
			})
		})

		When("user searches for users with selected user Ids and first 10 users with descending order on name property", func() {
			It("should return first 10 users in descending order on name field", func() {
				first := 10
				searchRequest := repository.SearchRequest{
					UserIDs: userIDs,
					Pagination: common.Pagination{
						After: nil,
						First: &first,
					},
					SortingOptions: []common.SortingOptionPair{
						{Name: "email", Direction: common.Descending},
					},
				}

				response, err := sut.Search(ctx, &searchRequest)
				Ω(err).Should(BeNil())
				Ω(response.Users).ShouldNot(BeNil())
				Ω(len(response.Users)).Should(Equal(10))
				for i := 0; i < 10; i++ {
					Ω(response.Users[i].UserID).Should(Equal(userIDs[9-i]))
					email := fmt.Sprintf("Name%d@test.com", 9-i)
					Ω(response.Users[i].User.Email).Should(Equal(email))
				}
			})
		})

	})

})

func assertUser(user, expectedUser models.User) {
	Ω(user).ShouldNot(BeNil())
	Ω(user.Email).Should(Equal(expectedUser.Email))
}
