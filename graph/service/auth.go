package service

import (
	"context"
	"errors"
	"time"

	"github.com/0x726f6f6b6965/go-simple-graphql/database"
	"github.com/0x726f6f6b6965/go-simple-graphql/graph/model"
	"github.com/0x726f6f6b6965/go-simple-graphql/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// create a new service
type UserService struct{}

// Register returns JWT token for authentication
func (u *UserService) Register(input model.NewUser) string {
	// create a password with bcrypt encryption
	bs, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}

	// create a variable to store the encrypted password
	var password string = string(bs)

	// create a new user
	var user model.User = model.User{
		Username:  input.Username,
		Email:     input.Email,
		Password:  password,
		CreatedAt: time.Now(),
	}

	// get the "users" collection
	var collection *mongo.Collection = database.GetCollection(utils.USER_COLLECTION)

	// add a new user to the "users" collection
	res, err := collection.InsertOne(context.TODO(), user)

	// if a user failed to add, return an empty string
	if err != nil {
		return ""
	}

	// convert ObjectID into the string
	var userId string = res.InsertedID.(primitive.ObjectID).Hex()

	// generate a new JWT token
	token, err := utils.GenerateNewAccessToken(userId)

	// if token generation failed, return an empty string
	if err != nil {
		return ""
	}

	// return the JWT token
	return token
}

// Login returns JWT token for authentication
func (u *UserService) Login(input model.LoginInput) string {
	// get the "users" collection from the database
	var collection *mongo.Collection = database.GetCollection(utils.USER_COLLECTION)

	// create a variable to store a user from the database
	var user *model.User = &model.User{}
	// create a filter query to filter data by user email
	filter := bson.M{"email": input.Email}

	// find the user data by email
	var res *mongo.SingleResult = collection.FindOne(context.TODO(), filter)

	// if a user is not found, return the empty string
	if err := res.Decode(user); err != nil {
		return ""
	}

	// compare the user password with the password from the input
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))

	// If the password does not match, return the empty string
	if err != nil {
		return ""
	}

	// generate a JWT token
	token, err := utils.GenerateNewAccessToken(user.ID)

	// if token generation failed, return the empty string
	if err != nil {
		return ""
	}

	// return the JWT token
	return token
}

func (u *UserService) GetUser(id string) (*model.User, error) {
	// create an ObjectID from id
	userID, err := primitive.ObjectIDFromHex(id)

	// if ObjectID is failed to create, return an error
	if err != nil {
		return &model.User{}, errors.New("id is invalid")
	}

	// create a query to filter data by id (_id)
	var query primitive.D = bson.D{{Key: "_id", Value: userID}}

	// get the "users" collection from the database
	var collection *mongo.Collection = database.GetCollection(utils.USER_COLLECTION)

	// get the user data by ID
	var userData *mongo.SingleResult = collection.FindOne(context.TODO(), query)

	// if user data is not found return an error
	if userData.Err() != nil {
		return &model.User{}, errors.New("user not found")
	}

	// create a variable to store the user from the database
	var user *model.User = &model.User{}

	// decode the user data from the database into the "user" variable
	userData.Decode(user)

	// return the user from the database
	return user, nil
}
