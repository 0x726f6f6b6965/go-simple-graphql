package main

import (
	"log"
	"net/http"
	"os"

	"github.com/0x726f6f6b6965/go-simple-graphql/database"
	"github.com/0x726f6f6b6965/go-simple-graphql/graph"
	"github.com/0x726f6f6b6965/go-simple-graphql/graph/middleware"
	"github.com/0x726f6f6b6965/go-simple-graphql/utils"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	var handler *chi.Mux = NewGraphQLHandler()

	// connect to the database
	err := database.Connect(utils.GetValue("DATABASE_NAME"))
	if err != nil {
		log.Fatalf("Cannot connect to the database: %v\n", err)
	}

	log.Println("Connected to the database")

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

// NewGraphQLHandler returns handler for GraphQL application
func NewGraphQLHandler() *chi.Mux {
	// create a new router
	var router *chi.Mux = chi.NewRouter()

	// use the middleware component
	router.Use(middleware.NewMiddleware())

	// create a GraphQL server
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	// assign some handlers for the GraphQL server
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	// return the handler
	return router
}
