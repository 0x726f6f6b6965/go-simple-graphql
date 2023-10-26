//go:generate go get -d github.com/99designs/gqlgen
//go:generate go run ./../cmd/gen/generate.go
package graph

import "github.com/0x726f6f6b6965/go-simple-graphql/graph/service"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	blogService service.BlogService
	userService service.UserService
}
