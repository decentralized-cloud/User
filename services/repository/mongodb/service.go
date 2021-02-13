// Package mongodb implements MongoDB repository services
package mongodb

import (
	"context"
	"fmt"

	"github.com/decentralized-cloud/user/models"
	"github.com/decentralized-cloud/user/services/configuration"
	"github.com/decentralized-cloud/user/services/repository"
	"github.com/micro-business/go-core/common"
	commonErrors "github.com/micro-business/go-core/system/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

	insertResult, err := collection.InsertOne(ctx, request.User)
	if err != nil {
		return nil, repository.NewUnknownErrorWithError("Insert user failed.", err)
	}

	userID := insertResult.InsertedID.(primitive.ObjectID).Hex()

	return &repository.CreateUserResponse{
		UserID: userID,
		User:   request.User,
		Cursor: userID,
	}, nil
}

// ReadUser read an existing user
// context: Optional The reference to the context
// request: Mandatory. The request to read an existing user
// Returns either the result of reading an existing user or error if something goes wrong.
func (service *mongodbRepositoryService) ReadUser(
	ctx context.Context,
	request *repository.ReadUserRequest) (*repository.ReadUserResponse, error) {
	client, collection, err := service.createClientAndCollection(ctx)
	if err != nil {
		return nil, err
	}

	defer disconnect(ctx, client)

	ObjectID, _ := primitive.ObjectIDFromHex(request.UserID)
	filter := bson.D{{Key: "_id", Value: ObjectID}}
	var user models.User

	err = collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, repository.NewUserNotFoundError(request.UserID)
	}

	return &repository.ReadUserResponse{
		User: user,
	}, nil
}

// ReadUserByEmail read an existing user by email address
// ctx: Mandatory The reference to the context
// request: Mandatory. The request to read an existing user address
// Returns either the result of reading an existing user by email address or error if something goes wrong.
func (service *mongodbRepositoryService) ReadUserByEmail(
	ctx context.Context,
	request *repository.ReadUserByEmailRequest) (*repository.ReadUserByEmailResponse, error) {
	client, collection, err := service.createClientAndCollection(ctx)
	if err != nil {
		return nil, err
	}

	defer disconnect(ctx, client)

	filter := bson.D{{Key: "email", Value: request.Email}}
	var user models.User

	err = collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, repository.NewUserByEmailNotFoundError(request.Email)
	}

	return &repository.ReadUserByEmailResponse{
		User: user,
	}, nil
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

	ObjectID, _ := primitive.ObjectIDFromHex(request.UserID)
	filter := bson.D{{Key: "_id", Value: ObjectID}}

	newUser := bson.M{"$set": bson.M{"email": request.User.Email}}
	response, err := collection.UpdateOne(ctx, filter, newUser)

	if err != nil {
		return nil, repository.NewUnknownErrorWithError("Update user failed.", err)
	}

	if response.MatchedCount == 0 {
		return nil, repository.NewUserNotFoundError(request.UserID)
	}

	return &repository.UpdateUserResponse{
		User:   request.User,
		Cursor: request.UserID,
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

	ObjectID, _ := primitive.ObjectIDFromHex(request.UserID)
	filter := bson.D{{Key: "_id", Value: ObjectID}}
	response, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, repository.NewUnknownErrorWithError("Delete user failed.", err)
	}

	if response.DeletedCount == 0 {
		return nil, repository.NewUserNotFoundError(request.UserID)
	}

	return &repository.DeleteUserResponse{}, nil
}

// Search returns the list of users that matched the criteria
// ctx: Mandatory The reference to the context
// request: Mandatory. The request contains the search criteria
// Returns the list of users that matched the criteria
func (service *mongodbRepositoryService) Search(
	ctx context.Context,
	request *repository.SearchRequest) (*repository.SearchResponse, error) {
	response := &repository.SearchResponse{
		HasPreviousPage: false,
		HasNextPage:     false,
	}

	ids := []primitive.ObjectID{}
	for _, userID := range request.UserIDs {
		objectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			return nil, repository.NewUnknownErrorWithError(fmt.Sprintf("Failed to decode the userID: %s.", userID), err)
		}

		ids = append(ids, objectID)
	}

	filter := bson.M{}
	if len(request.UserIDs) > 0 {
		filter["_id"] = bson.M{"$in": ids}
	}

	client, collection, err := service.createClientAndCollection(ctx)
	if err != nil {
		return nil, err
	}

	defer disconnect(ctx, client)

	response.TotalCount, err = collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, repository.NewUnknownErrorWithError("Failed to retrieve the number of users that match the filter criteria", err)
	}

	if response.TotalCount == 0 {
		// No tennat matched the filter criteria
		return response, nil
	}

	if request.Pagination.After != nil {
		after := *request.Pagination.After
		objectID, err := primitive.ObjectIDFromHex(after)
		if err != nil {
			return nil, repository.NewUnknownErrorWithError(fmt.Sprintf("Failed to decode the After: %s.", after), err)
		}

		if len(filter) > 0 {
			filter["$and"] = []interface{}{
				bson.M{"_id": bson.M{"$gt": objectID}},
			}
		} else {
			filter["_id"] = bson.M{"$gt": objectID}
		}
	}

	if request.Pagination.Before != nil {
		before := *request.Pagination.Before
		objectID, err := primitive.ObjectIDFromHex(before)
		if err != nil {
			return nil, repository.NewUnknownErrorWithError(fmt.Sprintf("Failed to decode the Before: %s.", before), err)
		}

		if len(filter) > 0 {
			filter["$and"] = []interface{}{
				bson.M{"_id": bson.M{"$lt": objectID}},
			}
		} else {
			filter["_id"] = bson.M{"$lt": objectID}
		}
	}

	findOptions := options.Find()

	if request.Pagination.First != nil {
		findOptions.SetLimit(int64(*request.Pagination.First))
	}

	if request.Pagination.Last != nil {
		findOptions.SetLimit(int64(*request.Pagination.Last))
	}

	if len(request.SortingOptions) > 0 {
		var sortOptionPairs bson.D

		for _, sortingOption := range request.SortingOptions {
			direction := 1
			if sortingOption.Direction == common.Descending {
				direction = -1
			}

			sortOptionPairs = append(
				sortOptionPairs,
				bson.E{
					Key:   sortingOption.Name,
					Value: direction,
				})
		}

		findOptions.SetSort(sortOptionPairs)
	}

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, repository.NewUnknownErrorWithError("Failed to call the Find function on the collection.", err)
	}

	users := []models.UserWithCursor{}
	for cursor.Next(ctx) {
		var user models.User
		//TODO : below line need to be removed, if we pass 'ShowRecordID' in findOption, ObjectID will be available
		var userBson bson.M

		err := cursor.Decode(&user)
		if err != nil {
			return nil, repository.NewUnknownErrorWithError("Failed to decode the user", err)
		}

		err = cursor.Decode(&userBson)
		if err != nil {
			return nil, repository.NewUnknownErrorWithError("Could not load the data.", err)
		}

		userID := userBson["_id"].(primitive.ObjectID).Hex()
		userWithCursor := models.UserWithCursor{
			UserID: userID,
			User:   user,
			Cursor: userID,
		}

		users = append(users, userWithCursor)
	}

	response.Users = users
	if (request.Pagination.After != nil && request.Pagination.First != nil && int64(*request.Pagination.First) < response.TotalCount) ||
		(request.Pagination.Before != nil && request.Pagination.Last != nil && int64(*request.Pagination.Last) < response.TotalCount) {
		response.HasNextPage = true
		response.HasPreviousPage = true
	} else if request.Pagination.After == nil && request.Pagination.First != nil && int64(*request.Pagination.First) < response.TotalCount {
		response.HasNextPage = true
		response.HasPreviousPage = false
	} else if request.Pagination.Before == nil && request.Pagination.Last != nil && int64(*request.Pagination.Last) < response.TotalCount {
		response.HasNextPage = false
		response.HasPreviousPage = true
	}

	return response, nil
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
