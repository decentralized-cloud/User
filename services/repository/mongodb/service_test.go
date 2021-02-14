package mongodb_test

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/decentralized-cloud/user/models"
	configurationMock "github.com/decentralized-cloud/user/services/configuration/mock"
	"github.com/decentralized-cloud/user/services/repository"
	"github.com/decentralized-cloud/user/services/repository/mongodb"
	"github.com/golang/mock/gomock"
	"github.com/lucsky/cuid"
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
			Email: cuid.New() + "@test.com",
			User:  models.User{}}
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
				Ω(response.Cursor).ShouldNot(BeNil())
				assertUser(response.User, createRequest.User)
			})
		})
	})

	Context("user already exists", func() {
		var (
			email string
		)

		BeforeEach(func() {
			_, _ = sut.CreateUser(ctx, &createRequest)
			email = createRequest.Email
		})

		When("user reads a user by Id", func() {
			It("should return a user", func() {
				response, err := sut.ReadUser(ctx, &repository.ReadUserRequest{Email: email})
				Ω(err).Should(BeNil())
				assertUser(response.User, createRequest.User)
			})
		})

		When("user updates the existing user", func() {
			It("should update the user information", func() {
				updateRequest := repository.UpdateUserRequest{
					Email: email,
					User:  models.User{}}

				updateResponse, err := sut.UpdateUser(ctx, &updateRequest)
				Ω(err).Should(BeNil())
				Ω(updateResponse.Cursor).ShouldNot(BeNil())
				assertUser(updateResponse.User, updateRequest.User)

				readResponse, err := sut.ReadUser(ctx, &repository.ReadUserRequest{Email: email})
				Ω(err).Should(BeNil())
				assertUser(readResponse.User, updateRequest.User)
			})
		})

		When("user deletes the user", func() {
			It("should delete the user", func() {
				_, err := sut.DeleteUser(ctx, &repository.DeleteUserRequest{Email: email})
				Ω(err).Should(BeNil())

				response, err := sut.ReadUser(ctx, &repository.ReadUserRequest{Email: email})
				Ω(err).Should(HaveOccurred())
				Ω(response).Should(BeNil())

				Ω(repository.IsUserNotFoundError(err)).Should(BeTrue())

				var notFoundErr repository.UserNotFoundError
				_ = errors.As(err, &notFoundErr)

				Ω(notFoundErr.Email).Should(Equal(email))
			})
		})
	})

	Context("user does not exist", func() {
		var (
			email string
		)

		BeforeEach(func() {
			email = cuid.New() + "@test.com"
		})

		When("user reads the user", func() {
			It("should return NotFoundError", func() {
				response, err := sut.ReadUser(ctx, &repository.ReadUserRequest{Email: email})
				Ω(err).Should(HaveOccurred())
				Ω(response).Should(BeNil())

				Ω(repository.IsUserNotFoundError(err)).Should(BeTrue())

				var notFoundErr repository.UserNotFoundError
				_ = errors.As(err, &notFoundErr)

				Ω(notFoundErr.Email).Should(Equal(email))
			})
		})

		When("user tries to update the user", func() {
			It("should return NotFoundError", func() {
				updateRequest := repository.UpdateUserRequest{
					Email: email,
					User:  models.User{}}

				response, err := sut.UpdateUser(ctx, &updateRequest)
				Ω(err).Should(HaveOccurred())
				Ω(response).Should(BeNil())

				Ω(repository.IsUserNotFoundError(err)).Should(BeTrue())

				var notFoundErr repository.UserNotFoundError
				_ = errors.As(err, &notFoundErr)

				Ω(notFoundErr.Email).Should(Equal(email))
			})
		})

		When("user tries to delete the user", func() {
			It("should return NotFoundError", func() {
				response, err := sut.DeleteUser(ctx, &repository.DeleteUserRequest{Email: email})
				Ω(err).Should(HaveOccurred())
				Ω(response).Should(BeNil())

				Ω(repository.IsUserNotFoundError(err)).Should(BeTrue())

				var notFoundErr repository.UserNotFoundError
				_ = errors.As(err, &notFoundErr)

				Ω(notFoundErr.Email).Should(Equal(email))
			})
		})
	})

})

func assertUser(user, expectedUser models.User) {
	Ω(user).ShouldNot(BeNil())
}
