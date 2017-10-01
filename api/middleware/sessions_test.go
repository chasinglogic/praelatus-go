package middleware

import (
	"net/http/httptest"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/praelatus/praelatus/models"
)

func TestGetToken(t *testing.T) {
	testToken := jwt.New(jwt.SigningMethodHS256)
	signed, err := testToken.SignedString(signingKey)
	if err != nil {
		t.Error(err)
		return
	}

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Bearer "+signed)

	token := getToken(r)
	if token == nil {
		t.Error("Token was nil")
		return
	}

	if token.Raw != signed {
		t.Errorf("Expected:\n%s\nGot:\n%s\n", signed, token.Raw)
	}
}

func TestMakeClaims(t *testing.T) {
	u := models.User{
		Username: "test",
		FullName: "Test Testerson",
		Email:    "test@example.com",
		IsAdmin:  false,
	}

	claims := makeClaims(u)

	if claims["username"] != u.Username {
		t.Errorf("Expected %s Got %s\n", u.Username, claims["username"])
	}

	if claims["email"] != u.Email {
		t.Errorf("Expected %s Got %s\n", u.Email, claims["email"])
	}

	if claims["is_admin"] != u.IsAdmin {
		t.Errorf("Expected %v Got %v\n", u.IsAdmin, claims["is_admin"])
	}

	if _, ok := claims["iat"]; !ok {
		t.Error("Expected iat to be set but is not.")
	}

	if _, ok := claims["exp"]; !ok {
		t.Error("Expected exp to be set but is not.")
	}
}

func TestMakeAdminClaims(t *testing.T) {
	u := models.User{
		Username: "test",
		FullName: "Test Testerson",
		Email:    "test@example.com",
		IsAdmin:  true,
	}

	claims := makeClaims(u)

	if claims["username"] != u.Username {
		t.Errorf("Expected %s Got %s\n", u.Username, claims["username"])
	}

	if claims["email"] != u.Email {
		t.Errorf("Expected %s Got %s\n", u.Email, claims["email"])
	}

	if claims["is_admin"] != u.IsAdmin {
		t.Errorf("Expected %v Got %v\n", u.IsAdmin, claims["is_admin"])
	}

	if _, ok := claims["iat"]; !ok {
		t.Error("Expected iat to be set but is not.")
	}

	if _, ok := claims["exp"]; !ok {
		t.Error("Expected exp to be set but is not.")
	}
}

func TestUserFromClaims(t *testing.T) {
	u := models.User{
		Username: "test",
		FullName: "Test Testerson",
		Email:    "test@example.com",
		IsAdmin:  true,
	}

	claims := makeClaims(u)
	newUser := userFromClaims(claims)

	if u.Username != newUser.Username {
		t.Errorf("Expected %s Got %s", u.Username, newUser.Username)
	}

	if u.Email != newUser.Email {
		t.Errorf("Expected %s Got %s", u.Email, newUser.Email)
	}

	if u.IsAdmin != newUser.IsAdmin {
		t.Errorf("Expected %v Got %v", u.IsAdmin, newUser.IsAdmin)
	}
}

func TestTokenIntegration(t *testing.T) {
	u := models.User{
		Username: "test",
		FullName: "Test Testerson",
		Email:    "test@example.com",
		IsAdmin:  true,
	}

	w := httptest.NewRecorder()

	err := SetUserSession(u, w)
	if err != nil {
		t.Error(err)
		return
	}

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Bearer "+w.Header().Get("X-Praelatus-Token"))

	returnedUser := GetUserSession(r)

	if u.Username != returnedUser.Username {
		t.Errorf("Got: %s Expected: %s\n", returnedUser.Username, u.Username)
	}

	if u.Email != returnedUser.Email {
		t.Errorf("Got: %s Expected: %s\n", returnedUser.Email, u.Email)
	}

	if u.IsAdmin != returnedUser.IsAdmin {
		t.Errorf("Got: %v Expected: %v\n", returnedUser.IsAdmin, u.IsAdmin)
	}
}
