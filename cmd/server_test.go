package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/0x726f6f6b6965/go-simple-graphql/database"
	"github.com/0x726f6f6b6965/go-simple-graphql/graph/model"
	"github.com/0x726f6f6b6965/go-simple-graphql/mock"
	"github.com/0x726f6f6b6965/go-simple-graphql/utils"
	"github.com/joho/godotenv"

	"github.com/steinfletcher/apitest"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	var err error = godotenv.Load("./../.env")
	if err != nil {
		fmt.Println("Error while loading .env file")
	}
	err = database.Connect(os.Getenv("TEST_DATABASE_NAME"))
	if err != nil {
		fmt.Printf("Cannot connect to the database: %v, %s\n", err, os.Getenv("TEST_DATABASE_NAME"))
		return
	}
	fmt.Printf("\033[1;33m%s\033[0m", "> Setup completed\n")
}

func teardown() {
	database.Mongo.Client.Disconnect(context.Background())
	fmt.Printf("\033[1;33m%s\033[0m", "> Teardown completed")
	fmt.Printf("\n")
}

func TestSignup_Success(t *testing.T) {

	// create a test
	apitest.New().
		// run the cleanup() function after the test is finished
		Observe(cleanup).
		// add an application to be tested
		Handler(NewGraphQLHandler()).
		// send a POST request for sign-up
		Post("/query").
		// define the query for sign-up
		GraphQLQuery(`mutation {
            register(input:{
                email:"test@test.com",
                username:"test",
                password:"123123"
            })
        }`).
		// expect the status code is equals to 200
		Expect(t).
		Status(http.StatusOK).
		End()
}

func TestLogin_Success(t *testing.T) {

	// create a new user data
	user := getUser()

	// create a query to login
	var query string = `mutation {
        login(input:{
            email:"` + user.Email + `",
            password:"` + user.Password + `"
        })
    }`

	// create a test
	apitest.New().
		// run the cleanup() function after the test is finished
		Observe(cleanup).
		// add an application to be tested
		Handler(NewGraphQLHandler()).
		// send a POST request for login
		Post("/query").
		// define the query for the login
		GraphQLQuery(query).
		// expect the status code is equals to 200
		Expect(t).
		Status(http.StatusOK).
		End()
}

func TestLogin_Failed(t *testing.T) {
	var result string = `{
		"errors": [
			{
				"message": "login failed, invalid email or password",
				"path": [
					"login"
				]
			}
		],
		"data": null
	}`

	var query string = `mutation {
		login(input:{
			email:"wrong@mail.com",
			password:"123456"
		})
	}`

	apitest.New().
		Observe(cleanup).
		Handler(NewGraphQLHandler()).
		Post("/query").
		GraphQLQuery(query).
		Expect(t).
		Status(http.StatusOK).
		Body(result).
		End()
}

func TestGetBlogs_Success(t *testing.T) {
	// create a test
	apitest.New().
		// add an application to be tested
		Handler(NewGraphQLHandler()).
		// send a POST request for getting all blogs
		Post("/query").
		// define the query for getting all blogs
		GraphQLQuery(`query { blogs { title } }`).
		// expect the status code is equals to 200
		Expect(t).
		Status(http.StatusOK).
		End()
}

func TestGetBlog_Success(t *testing.T) {
	// create a blog data
	var blog model.Blog = getBlog()

	// create a query for getting a blog by ID
	var query string = `query {
        blog(id:"` + blog.ID + `") {
            title
            content
        }
    }`

	// create an expected result body
	var result string = `{
        "data": {
            "blog": {
                "title": "` + blog.Title + `",
                "content": "` + blog.Content + `"
            }
        }
    }`

	// create a test
	apitest.New().
		// run the cleanup() function after the test is finished
		Observe(cleanup).
		// add an application to be tested
		Handler(NewGraphQLHandler()).
		// send a POST request for getting a blog by ID
		Post("/query").
		// define the query for getting a blog by ID
		GraphQLQuery(query).
		// expect the status code is equals to 200
		Expect(t).
		Status(http.StatusOK).
		// expect the response body is equal to the expected result
		Body(result).
		End()
}

func TestGetBlog_Failed(t *testing.T) {
	// create a query to get the blog by ID
	var query string = `query {
        blog(id:"62c24f3b4896bb25c21e49b9") {
            title
            content
        }
    }`

	// create an expected result body
	var result string = `{
        "errors": [
            {
                "message": "blog not found",
                "path": [
                    "blog"
                ]
            }
        ],
        "data": null
    }`

	// create a test
	apitest.New().
		// run the cleanup() function after the test is finished
		Observe(cleanup).
		// add an application to be tested
		Handler(NewGraphQLHandler()).
		// send a POST request for getting a blog by ID
		Post("/query").
		// define the query for getting a blog by ID
		GraphQLQuery(query).
		// expect the status code is equals to 200
		Expect(t).
		Status(http.StatusOK).
		// expect the response body is equal to the expected result
		Body(result).
		End()
}

func TestCreateBlog_Success(t *testing.T) {
	// create a new user data
	var author model.User = getUser()

	// generate a JWT token for authentication
	var token string = getJWTToken(author)

	// create a query for creating a new blog
	var query string = `mutation {
        newBlog(input:{
            title:"my blog",
            content:"this is the content"
        }) {
            title
            content
        }
    }`

	// create an expected result body
	var result string = `{
        "data": {
            "newBlog": {
                "title": "my blog",
                "content": "this is the content"
            }
        }
    }`

	// create a test
	apitest.New().
		// run the cleanup() function after the test is finished
		Observe(cleanup).
		// add an application to be tested
		Handler(NewGraphQLHandler()).
		// send a POST request for creating a new blog
		Post("/query").
		// attach the JWT token to the Authorization header
		Header("Authorization", token).
		// define the query for creating a new blog
		GraphQLQuery(query).
		// expect the status code is equals to 200
		Expect(t).
		Status(http.StatusOK).
		// expect the response body is equal to the expected result
		Body(result).
		End()
}

func TestCreateBlog_Failed(t *testing.T) {
	// create a query for creating a new blog
	var query string = `mutation {
        newBlog(input:{
            title:"my blog",
            content:"this is the content"
        }) {
            title
            content
        }
    }`

	// create an expected result body
	var result string = `{
        "errors": [
            {
                "message": "access denied",
                "path": [
                    "newBlog"
                ]
            }
        ],
        "data": null
    }`

	// create a test
	apitest.New().
		// run the cleanup() function after the test is finished
		Observe(cleanup).
		// add an application to be tested
		Handler(NewGraphQLHandler()).
		// send a POST request for creating a new blog
		Post("/query").
		// define the query for creating a new blog
		GraphQLQuery(query).
		// expect the status code is equals to 200
		Expect(t).
		Status(http.StatusOK).
		// expect the response body is equal to the expected result
		Body(result).
		End()
}

func TestEditBlog_Success(t *testing.T) {
	// create a new blog data
	var blog model.Blog = getBlog()
	// generate a JWT token for authentication
	var token string = getJWTToken(*blog.Author)

	// create a query for updating a blog
	var query string = `mutation {
        editBlog(input:{
            blogId:"` + blog.ID + `"
            title:"my blog",
            content:"this is the content"
        }) {
            title
            content
        }
    }`

	// create an expected result body
	var result string = `{
        "data": {
            "editBlog": {
                "title": "my blog",
                "content": "this is the content"
            }
        }
    }`

	// create a test
	apitest.New().
		// run the cleanup() function after the test is finished
		Observe(cleanup).
		// add an application to be tested
		Handler(NewGraphQLHandler()).
		// send a POST request for updating a blog
		Post("/query").
		// attach the JWT token to the Authorization header
		Header("Authorization", token).
		// define the query for updating a blog
		GraphQLQuery(query).
		// expect the status code is equals to 200
		Expect(t).
		Status(http.StatusOK).
		// expect the response body is equal to the expected result
		Body(result).
		End()
}

func TestEditBlog_Failed(t *testing.T) {
	// create a new blog data
	var blog model.Blog = getBlog()

	// create a query for updating a blog
	var query string = `mutation {
        editBlog(input:{
            blogId:"` + blog.ID + `"
            title:"my blog",
            content:"this is the content"
        }) {
            title
            content
        }
    }`

	// create an expected result body
	var result string = `{
        "errors": [
            {
                "message": "access denied",
                "path": [
                    "editBlog"
                ]
            }
        ],
        "data": null
    }`

	// create a test
	apitest.New().
		// run the cleanup() function after the test is finished
		Observe(cleanup).
		// add an application to be tested
		Handler(NewGraphQLHandler()).
		// send a POST request for updating a blog
		Post("/query").
		// define the query for updating a blog
		GraphQLQuery(query).
		// expect the status code is equals to 200
		Expect(t).
		Status(http.StatusOK).
		// expect the response body is equal to the expected result
		Body(result).
		End()
}

func TestDeleteBlog_Success(t *testing.T) {
	// create a new blog data
	var blog model.Blog = getBlog()

	// generate a JWT token for authentication
	var token string = getJWTToken(*blog.Author)

	// create a query for deleting a blog
	var query string = `mutation {
        deleteBlog(input:{
            blogId:"` + blog.ID + `"
        })
    }`

	// create an expected result body
	var result string = `{
        "data": {
            "deleteBlog": true
        }
    }`

	// create a test
	apitest.New().
		// run the cleanup() function after the test is finished
		Observe(cleanup).
		// add an application to be tested
		Handler(NewGraphQLHandler()).
		// send a POST request for deleting a blog
		Post("/query").
		// attach the JWT token to the Authorization header
		Header("Authorization", token).
		// define the query for deleting a blog
		GraphQLQuery(query).
		// expect the status code is equals to 200
		Expect(t).
		Status(http.StatusOK).
		// expect the response body is equal to the expected result
		Body(result).
		End()
}

func TestDeleteBlog_Failed(t *testing.T) {
	// create a new blog data
	var blog model.Blog = getBlog()

	// create a query for deleting a blog
	var query string = `mutation {
        deleteBlog(input:{
            blogId:"` + blog.ID + `"
        })
    }`

	// create an expected result body
	var result string = `{
        "errors": [
            {
                "message": "access denied",
                "path": [
                    "deleteBlog"
                ]
            }
        ],
        "data": null
    }`

	// create a test
	apitest.New().
		// run the cleanup() function after the test is finished
		Observe(cleanup).
		// add an application to be tested
		Handler(NewGraphQLHandler()).
		// send a POST request for deleting a blog
		Post("/query").
		// define the query for deleting a blog
		GraphQLQuery(query).
		// expect the status code is equals to 200
		Expect(t).
		Status(http.StatusOK).
		// expect the response body is equal to the expected result
		Body(result).
		End()
}

func cleanup(res *http.Response, req *http.Request, apiTest *apitest.APITest) {
	if http.StatusOK == res.StatusCode {
		mock.CleanSeeders()
	}
}

func getJWTToken(user model.User) string {
	// generate JWT token
	token, err := utils.GenerateNewAccessToken(user.ID)

	// if token generation failed, return an error
	if err != nil {
		panic(err)
	}

	// return the JWT token for the authorization header
	return "Bearer " + token
}

func getBlog() model.Blog {

	// create a new blog data for testing
	blog, err := mock.SeedBlog()

	// if blog creation failed, return an error
	if err != nil {
		panic(err)
	}

	// return the blog
	return blog
}

func getUser() model.User {

	// create new user data for testing
	user, err := mock.SeedUser()

	// if user creation failed, return an error
	if err != nil {
		panic(err)
	}

	// return the user
	return user
}
