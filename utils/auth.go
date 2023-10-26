package utils

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// TokenMetadata represents JWT token metadata
type TokenMetadata struct {
	Expires int64
	UserId  string
}

// ExtractTokenMetadata extracts JWT token metadata
func ExtractTokenMetadata(r *http.Request) (*TokenMetadata, error) {
	// verify the JWT token
	token, err := verifyToken(r)

	// if verification is failed, return an error
	if err != nil {
		return nil, err
	}

	// get a JWT claim from the JWT token
	claims, ok := token.Claims.(jwt.MapClaims)

	// check if the token is valid
	var isValid bool = ok && token.Valid

	// if the JWT token is valid, return the JWT token metadata
	if isValid {
		// set token expiration
		expires := int64(claims["exp"].(float64))
		// set user ID for the token
		userId := claims["userId"].(string)

		// return the JWT token metadata
		return &TokenMetadata{
			Expires: expires,
			UserId:  userId,
		}, nil
	}

	// return an error
	return nil, err
}

// GenerateNewAccessToken generates a new JWT token
func GenerateNewAccessToken(userId string) (string, error) {
	// get the JWT secret key from .env file
	secret := GetValue("JWT_SECRET_KEY")

	// get the JWT secret key expiration in minutes from .env file
	minutesCount, _ := strconv.Atoi(GetValue("JWT_SECRET_KEY_EXPIRE_MINUTES_COUNT"))

	// create a JWT claim
	claims := jwt.MapClaims{}

	// assign an expiration time for the token
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(minutesCount)).Unix()
	// assign a data for user ID
	claims["userId"] = userId

	// create a JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// convert the JWT token into the string
	t, err := token.SignedString([]byte(secret))

	// if conversion is failed, return an error
	if err != nil {
		return "", err
	}

	// return the generated JWT token
	return t, nil
}

// CheckToken checks JWT token
func CheckToken(r *http.Request) (*TokenMetadata, error) {
	// get the current time
	var now int64 = time.Now().Unix()

	// extract the JWT token metadata
	claims, err := ExtractTokenMetadata(r)
	// if extraction is failed, return an error
	if err != nil {
		return nil, err
	}

	// get the expiration time
	var expires int64 = claims.Expires

	// if the token is expired, return an error
	if now > expires {
		return nil, err
	}

	// return JWT claims from the JWT token
	return claims, nil
}

// verifyToken verifies JWT token
func verifyToken(r *http.Request) (*jwt.Token, error) {
	// get the token
	var tokenString string = extractToken(r)

	// parse the JWT token
	token, err := jwt.Parse(tokenString, jwtKeyFunc)

	// if parsing is failed, return an error
	if err != nil {
		return nil, err
	}

	// return JWT token
	return token, nil
}

// extractToken extracts JWT token from the Authorization header
func extractToken(r *http.Request) string {
	// get the Authorization header
	var header string = r.Header.Get("Authorization")
	// split the content inside the header to get the JWT token
	token := strings.Split(header, " ")

	// check if the JWT token is empty
	var isEmpty bool = header == "" || len(token) < 2

	// if the JWT token is empty return an empty string
	if isEmpty {
		return ""
	}

	// return JWT token from the header
	return token[1]
}

// jwtKeyFunc return JWT secret key
func jwtKeyFunc(token *jwt.Token) (interface{}, error) {
	return []byte(GetValue("JWT_SECRET_KEY")), nil
}
