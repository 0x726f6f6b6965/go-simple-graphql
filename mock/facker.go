package mock

import (
	"context"
	"errors"
	"math/rand"
	"reflect"
	"time"

	"github.com/0x726f6f6b6965/go-simple-graphql/database"
	"github.com/0x726f6f6b6965/go-simple-graphql/graph/model"
	"github.com/0x726f6f6b6965/go-simple-graphql/utils"

	"github.com/go-faker/faker/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserFaker struct {
	ID        string     `json:"id" bson:"_id,omitempty"`
	Username  string     `json:"username" bson:"username" faker:"username"`
	Email     string     `json:"email" bson:"email" faker:"email"`
	Password  string     `json:"password" bson:"password" faker:"password"`
	CreatedAt time.Time  `json:"createdAt" bson:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt" bson:"updatedAt"`
}

type BlogFaker struct {
	ID        string     `json:"id" bson:"_id,omitempty"`
	Title     string     `json:"title" bson:"title" faker:"title"`
	Content   string     `json:"content" bson:"content" faker:"content"`
	Author    *UserFaker `json:"author" bson:"author"`
	CreatedAt time.Time  `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt" bson:"updatedAt"`
}

func CreateFaker[T any]() (T, error) {
	var fakerData *T = new(T)
	_ = faker.AddProvider("content", func(v reflect.Value) (interface{}, error) {
		b := make([]rune, 20)
		for i := 0; i < 20; i++ {
			randRune := rune(rand.Intn(122-97) + 97)
			b[i] = randRune
		}
		return string(b), nil
	})

	_ = faker.AddProvider("title", func(v reflect.Value) (interface{}, error) {
		b := make([]rune, 10)
		for i := 0; i < 10; i++ {
			randRune := rune(rand.Intn(122-97) + 97)
			b[i] = randRune
		}
		return string(b), nil
	})
	err := faker.FakeData(fakerData)
	if err != nil {
		return *fakerData, err
	}

	return *fakerData, nil
}

func SeedUser() (model.User, error) {
	// create a faker for user data
	// this faker data will be stored in the database
	userFaker, err := CreateFaker[UserFaker]()

	// if faker creation failed, return an error
	if err != nil {
		return model.User{}, err
	}

	// create a password
	bs, err := bcrypt.GenerateFromPassword([]byte(userFaker.Password), bcrypt.DefaultCost)

	// if password creation failed, return an error
	if err != nil {
		return model.User{}, err
	}

	// get the password into the string format
	var password string = string(bs)

	// create a new user in the "user" variable
	var user model.User = model.User{
		Username:  userFaker.Username,
		Email:     userFaker.Email,
		Password:  password,
		CreatedAt: time.Now(),
	}

	// get the "users" collection
	var collection *mongo.Collection = database.GetCollection(utils.USER_COLLECTION)

	// insert the user into the "users" collection
	result, err := collection.InsertOne(context.TODO(), user)

	// if user insertion is failed, return an error
	if err != nil {
		return model.User{}, err
	}

	// assign the user password in original format (unencrypted)
	user.Password = userFaker.Password

	// assign the user ID in string format
	user.ID = result.InsertedID.(primitive.ObjectID).Hex()

	// return the recently created user
	return user, nil
}

func SeedBlog() (model.Blog, error) {
	// create a faker for blog data
	// this faker data will be stored in the database
	blogFaker, err := CreateFaker[BlogFaker]()

	// if faker creation failed, return an error
	if err != nil {
		return model.Blog{}, err
	}

	// create data for the author
	author, err := SeedUser()

	// if author creation failed, return an error
	if err != nil {
		return model.Blog{}, err
	}

	// create a new blog in the "blog" variable
	var blog model.Blog = model.Blog{
		Title:     blogFaker.Title,
		Content:   blogFaker.Content,
		Author:    &author,
		CreatedAt: time.Now(),
	}

	// get the "blogs" collection
	var collection *mongo.Collection = database.GetCollection(utils.BLOG_COLLECTION)

	// insert the blog into the "blogs" collection
	result, err := collection.InsertOne(context.TODO(), blog)

	// if blog insertion is failed, return an error
	if err != nil {
		return model.Blog{}, errors.New("create blog failed")
	}

	// create a filter to get the blog data by ID
	var filter primitive.D = bson.D{{Key: "_id", Value: result.InsertedID}}

	// get the recently created blog from the collection
	var createdRecord *mongo.SingleResult = collection.FindOne(context.TODO(), filter)

	// create a variable called "createdBlog" to store the recently created blog
	var createdBlog *model.Blog = &model.Blog{}

	// decode the recently created blog to the "createdBlog" variable
	createdRecord.Decode(createdBlog)

	// return the recently created blog
	return *createdBlog, nil
}

func CleanSeeders() {
	// get the "users" collection
	var userCollection *mongo.Collection = database.GetCollection(utils.USER_COLLECTION)

	// get the "blogs" collection
	var blogCollection *mongo.Collection = database.GetCollection(utils.BLOG_COLLECTION)

	// delete all data inside the "users" collection
	userErr := userCollection.Drop(context.TODO())

	// delete all data inside the "blogs" collection
	blogErr := blogCollection.Drop(context.TODO())

	// check if both operations are failed
	var isFailed bool = userErr != nil || blogErr != nil

	// if both operations are failed, return an error
	if isFailed {
		panic("error when deleting all data inside collection")
	}
}
