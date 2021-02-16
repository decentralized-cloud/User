// Package mongodb implements MongoDB repository services
package mongodb

import (
	"context"

	"github.com/decentralized-cloud/user/models"
	"github.com/decentralized-cloud/user/services/configuration"
	"github.com/decentralized-cloud/user/services/repository"
	commonErrors "github.com/micro-business/go-core/system/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type user struct {
	Email string `bson:"email" json:"email"`
}

type mongodbRepositoryService struct {
	connectionString       string
	databaseName           string
	databaseCollectionName string
}

// NewMongodbRepositoryService creates new instance of the mongodbRepositoryService, setting up all dependencies and returns the instance
// Returns the new service or error if something goes wrong
func NewMongodbRepositoryService(
	configurationService configuration.ConfigurationContract) (repository.RepositoryContract, error) {
	if configurationService == nil {
		return nil, commonErrors.NewArgumentNilError("configurationService", "configurationService is required")
	}

	connectionString, err := configurationService.GetDatabaseConnectionString()
	if err != nil {
		return nil, repository.NewUnknownErrorWithError("Failed to get connection string to mongodb", err)
	}

	databaseName, err := configurationService.GetDatabaseName()
	if err != nil {
		return nil, repository.NewUnknownErrorWithError("Failed to get the database name", err)
	}

	databaseCollectionName, err := configurationService.GetDatabaseCollectionName()
	if err != nil {
		return nil, repository.NewUnknownErrorWithError("Failed to get the database collection name", err)
	}

	return &mongodbRepositoryService{
		connectionString:       connectionString,
		databaseName:           databaseName,
		databaseCollectionName: databaseCollectionName,
	}, nil
}

// CreateUser creates a new user.
// context: Optional The reference to the context
// request: Mandatory. The request to create a new user
// Returns either the result of creating new user or error if something goes wrong.
func (service *mongodbRepositoryService) CreateUser(
	ctx context.Context,
	request *repository.CreateUserRequest) (*repository.CreateUserResponse, error) {
	client, collection, err := service.createClientAndCollection(ctx)
	if err != nil {
		return nil, err
	}

	defer disconnect(ctx, client)

	insertResult, err := collection.InsertOne(ctx, user{request.Email})
	if err != nil {
		return nil, repository.NewUnknownErrorWithError("User creation failed.", err)
	}

	userID := insertResult.InsertedID.(primitive.ObjectID).Hex()

	return &repository.CreateUserResponse{
		User:   request.User,
		Cursor: userID,
	}, nil
}

// ReadUser read an existing user
// ctx: Mandatory The reference to the context
// request: Mandatory. The request to read an existing user
// Returns either the result of reading an existing user or error if something goes wrong.
func (service *mongodbRepositoryService) ReadUser(
	ctx context.Context,
	request *repository.ReadUserRequest) (response *repository.ReadUserResponse, err error) {
	response, _, err = service.readUser(ctx, request)

	return
}

// UpdateUser update an existing user
// context: Optional The reference to the context
// request: Mandatory. The request to update an existing user
// Returns either the result of updateing an existing user or error if something goes wrong.
func (service *mongodbRepositoryService) UpdateUser(
	ctx context.Context,
	request *repository.UpdateUserRequest) (*repository.UpdateUserResponse, error) {
	client, collection, err := service.createClientAndCollection(ctx)
	if err != nil {
		return nil, err
	}

	defer disconnect(ctx, client)

	filter := bson.D{{Key: "email", Value: request.Email}}

	newUser := bson.M{"$set": bson.M{"email": request.Email}}
	response, err := collection.UpdateOne(ctx, filter, newUser)

	if err != nil {
		return nil, repository.NewUnknownErrorWithError("Update user failed.", err)
	}

	if response.MatchedCount == 0 {
		return nil, repository.NewUserNotFoundError(request.Email)
	}

	_, userID, err := service.readUser(ctx, &repository.ReadUserRequest{Email: request.Email})
	if err != nil {
		return nil, err
	}

	return &repository.UpdateUserResponse{
		User:   request.User,
		Cursor: userID,
	}, nil
}

// DeleteUser delete an existing user
// context: Optional The reference to the context
// request: Mandatory. The request to delete an existing user
// Returns either the result of deleting an existing user or error if something goes wrong.
func (service *mongodbRepositoryService) DeleteUser(
	ctx context.Context,
	request *repository.DeleteUserRequest) (*repository.DeleteUserResponse, error) {
	client, collection, err := service.createClientAndCollection(ctx)
	if err != nil {
		return nil, err
	}

	defer disconnect(ctx, client)

	filter := bson.D{{Key: "email", Value: request.Email}}
	response, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, repository.NewUnknownErrorWithError("Delete user failed.", err)
	}

	if response.DeletedCount == 0 {
		return nil, repository.NewUserNotFoundError(request.Email)
	}

	return &repository.DeleteUserResponse{}, nil
}

// ReadUser read an existing user
// ctx: Mandatory The reference to the context
// request: Mandatory. The request to read an existing user
// Returns either the result of reading an existing user or error if something goes wrong.
func (service *mongodbRepositoryService) readUser(
	ctx context.Context,
	request *repository.ReadUserRequest) (*repository.ReadUserResponse, string, error) {
	client, collection, err := service.createClientAndCollection(ctx)
	if err != nil {
		return nil, "", err
	}

	defer disconnect(ctx, client)

	filter := bson.D{{Key: "email", Value: request.Email}}
	var user user

	result := collection.FindOne(ctx, filter)
	err = result.Decode(&user)
	if err != nil {
		return nil, "", repository.NewUserNotFoundError(request.Email)
	}

	var userBson bson.M

	err = result.Decode(&userBson)
	if err != nil {
		return nil, "", repository.NewUnknownErrorWithError("Failed to load user bson data", err)
	}

	userID := userBson["_id"].(primitive.ObjectID).Hex()

	return &repository.ReadUserResponse{
		User: models.User{},
	}, userID, nil
}

func (service *mongodbRepositoryService) createClientAndCollection(ctx context.Context) (*mongo.Client, *mongo.Collection, error) {
	clientOptions := options.Client().ApplyURI(service.connectionString)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, nil, repository.NewUnknownErrorWithError("Could not connect to mongodb database.", err)
	}

	return client, client.Database(service.databaseName).Collection(service.databaseCollectionName), nil
}

func disconnect(ctx context.Context, client *mongo.Client) {
	_ = client.Disconnect(ctx)
}
