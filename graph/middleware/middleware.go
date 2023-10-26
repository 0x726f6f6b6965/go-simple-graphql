package middleware

import (
	"context"
	"net/http"

	"github.com/0x726f6f6b6965/go-simple-graphql/graph/model"
	"github.com/0x726f6f6b6965/go-simple-graphql/graph/service"
	"github.com/0x726f6f6b6965/go-simple-graphql/utils"
)

// create a context key
type contextKey struct {
	name string
}

// create a context key for user data
var userCtxKey = &contextKey{"user"}

// NewMiddleware returns a middleware for authentication
func NewMiddleware() func(http.Handler) http.Handler {
	// return handler that acts as a middleware
	return func(next http.Handler) http.Handler {
		// return handler function
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get header data from Authorization header
			var header string = r.Header.Get("Authorization")

			// if header data is empty
			// continue to serve HTTP
			if header == "" {
				next.ServeHTTP(w, r)
				return
			}

			// get the JWT token from the header
			tokenData, err := utils.CheckToken(r)

			// if the JWT token is invalid, return an error
			// the next request cannot be proceed
			if err != nil {
				http.Error(w, "invalid token", http.StatusForbidden)
				return
			}

			// create a userService component
			var userService service.UserService = service.UserService{}

			// get the user data by ID from the JWT token
			userData, err := userService.GetUser(tokenData.UserId)

			// if a user is not found, return an error
			// the next request cannot be proceed
			if err != nil {
				http.Error(w, "user not found", http.StatusForbidden)
				return
			}

			// store the user data from the database
			// into "user" variable
			var user model.User = *userData

			// create a context with value
			// the context value is user data
			ctx := context.WithValue(r.Context(), userCtxKey, &user)

			// add context to the request object
			r = r.WithContext(ctx)
			// continue to serve HTTP
			next.ServeHTTP(w, r)
		})
	}
}

// ForContext returns value from the context
func ForContext(ctx context.Context) *model.User {
	// get context value for user data
	raw, _ := ctx.Value(userCtxKey).(*model.User)
	// return context value
	return raw
}
