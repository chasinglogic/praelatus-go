package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/praelatus/praelatus/models"
)

// TODO: create a way to invalidate a session
// TODO: prevent session hijacking
// TODO: Make this actually secure
var signingKey = []byte("CHANGE ME")

func makeClaims(user models.User) jwt.MapClaims {
	now := time.Now()
	return jwt.MapClaims{
		"username": user.Username,
		"email":    user.Email,
		"is_admin": user.IsAdmin,
		"iat":      now,
		"exp":      now.Add(time.Hour * 24),
	}
}

func userFromClaims(claims jwt.MapClaims) models.User {
	maybeUsername, ok := claims["username"]
	if !ok {
		return models.User{}
	}

	maybeEmail, ok := claims["email"]
	if !ok {
		return models.User{}
	}

	maybeAdmin, ok := claims["is_admin"]
	if !ok {
		return models.User{}
	}

	username, ok := maybeUsername.(string)
	if !ok {
		return models.User{}
	}

	email, ok := maybeEmail.(string)
	if !ok {
		return models.User{}
	}

	isAdmin, ok := maybeAdmin.(bool)
	if !ok {
		return models.User{}
	}

	return models.User{
		Username: username,
		Email:    email,
		IsAdmin:  isAdmin,
	}
}

func getToken(r *http.Request) *jwt.Token {
	auth := r.Header.Get("Authorization")
	if len(auth) == 0 {
		return nil
	}

	tokenString := auth[len("Bearer "):]
	if tokenString == "" {
		return nil
	}

	token, err := jwt.Parse(tokenString,
		func(t *jwt.Token) (interface{}, error) {
			return signingKey, nil
		})
	if err != nil {
		fmt.Println("ERROR: [TOKEN_VALIDATION]", err.Error())
		return nil
	}

	return token
}

func getClaims(token *jwt.Token) jwt.MapClaims {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("ERROR: [TOKEN_VALIDATION] Invalid Claims", claims)
		return nil
	}

	return claims
}

// GetUserSession will check the given http.Request for a session token and if
// found it will return the corresponding user.
func GetUserSession(r *http.Request) *models.User {
	token := getToken(r)
	if token == nil {
		return nil
	}

	claims := getClaims(token)
	if claims == nil {
		return nil
	}

	user := userFromClaims(claims)
	return &user
}

// SetUserSession will generate a secure cookie for user u, will set the cookie
// on the response w and will add the user session to the session store
func SetUserSession(u models.User, w http.ResponseWriter) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, makeClaims(u))
	signed, err := token.SignedString(signingKey)
	if err != nil {
		return err
	}

	// TODO: Store the session with the client specific ID for session hijacking
	w.Header().Set("X-Praelatus-Token", signed)
	return nil
}

// RefreshSession will reset the expiry on the session for the given request
func RefreshSession(r *http.Request) error {
	token := getToken(r)
	if token == nil {
		return errors.New("no session on this request")
	}

	claims := getClaims(token)
	if claims == nil {
		return errors.New("no claims on token")
	}

	claims["exp"] = time.Now().Add(time.Hour * 24)
	return nil
}
