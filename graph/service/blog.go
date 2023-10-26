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
	"go.mongodb.org/mongo-driver/mongo/options"
)

// BlogService represents service component
type BlogService struct{}

func (b *BlogService) GetAllBlogs() []*model.Blog {
	var (
		query       primitive.D          = bson.D{{}}
		findOptions *options.FindOptions = options.Find()
	)

	findOptions.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := database.GetCollection(utils.BLOG_COLLECTION).Find(context.TODO(), query, findOptions)
	if err != nil {
		return []*model.Blog{}
	}

	blogs := make([]*model.Blog, 0)

	if err := cursor.All(context.TODO(), &blogs); err != nil {
		return []*model.Blog{}
	}

	return blogs
}

func (b *BlogService) GetBlogByID(id string) (*model.Blog, error) {
	blogID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &model.Blog{}, errors.New("id is invalid")
	}

	var (
		query      primitive.D       = bson.D{{Key: "_id", Value: blogID}}
		collection *mongo.Collection = database.GetCollection(utils.BLOG_COLLECTION)
	)

	blogData := collection.FindOne(context.TODO(), query)

	if blogData.Err() != nil {
		return &model.Blog{}, errors.New("blog not found")
	}

	blog := &model.Blog{}
	blogData.Decode(blog)

	return blog, nil
}

func (b *BlogService) CreateBlog(input model.NewBlog, user model.User) (*model.Blog, error) {
	var (
		blog model.Blog = model.Blog{
			Title:     input.Title,
			Content:   input.Content,
			Author:    &user,
			CreatedAt: time.Now(),
		}
		collection *mongo.Collection = database.GetCollection(utils.BLOG_COLLECTION)
	)

	result, err := collection.InsertOne(context.TODO(), blog)

	if err != nil {
		return &model.Blog{}, errors.New("create blog failed")
	}

	filter := bson.D{{Key: "_id", Value: result.InsertedID}}
	createdRecord := collection.FindOne(context.TODO(), filter)

	createdBlog := &model.Blog{}

	createdRecord.Decode(createdBlog)

	return createdBlog, nil
}

func (b *BlogService) EditBlog(input model.EditBlog, user model.User) (*model.Blog, error) {
	blogID, err := primitive.ObjectIDFromHex(input.BlogID)
	if err != nil {
		return &model.Blog{}, errors.New("id is invalid")
	}

	var (
		query primitive.D = bson.D{
			{Key: "_id", Value: blogID},
			{Key: "author._id", Value: user.ID},
		}

		update primitive.D = bson.D{{
			Key: "$set",
			Value: bson.D{
				{Key: "title", Value: input.Title},
				{Key: "content", Value: input.Content},
				{Key: "updatedAt", Value: time.Now()},
			},
		}}

		collection *mongo.Collection = database.GetCollection(utils.BLOG_COLLECTION)
	)

	updateResult := collection.FindOneAndUpdate(
		context.TODO(),
		query,
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	if updateResult.Err() != nil {
		if err == mongo.ErrNoDocuments {
			return &model.Blog{}, errors.New("blog not found")
		}
		return &model.Blog{}, errors.New("update blog failed")
	}

	var editedBlog *model.Blog = &model.Blog{}

	updateResult.Decode(editedBlog)

	return editedBlog, nil
}

func (b *BlogService) DeleteBlog(input model.DeleteBlog, user model.User) bool {
	blogID, err := primitive.ObjectIDFromHex(input.BlogID)
	if err != nil {
		return false
	}

	var (
		query primitive.D = bson.D{
			{Key: "_id", Value: blogID},
			{Key: "author._id", Value: user.ID},
		}
		collection *mongo.Collection = database.GetCollection(utils.BLOG_COLLECTION)
	)

	result, err := collection.DeleteOne(context.TODO(), query)
	isFailed := err != nil || result.DeletedCount < 1

	if isFailed {
		return !isFailed
	}

	return true
}
