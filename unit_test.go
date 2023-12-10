package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	router := gin.Default()
	initCookies(router)
	router.POST("/api/users/register", register)
	w := httptest.NewRecorder()

	// Create a request to send to the above route
	jsonParam := `{"username": "testuser", "password": "testpassword", "email": "test@test.test"}`
	req, err := http.NewRequest("POST", "/api/users/register", bytes.NewBufferString(jsonParam))

	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	// Check the status code is what we expect.
	assert.Equal(t, 200, w.Code)

	// Check the response body is what we expect.
	expected := `{"message":"Account created successfully!"}`
	assert.Equal(t, expected, w.Body.String())
}

func TestLoginAndLogout(t *testing.T) {

	router := gin.Default()
	initCookies(router) // Make sure this correctly initializes any required middleware
	router.POST("/api/users/login", login)
	router.GET("/logout", logout)

	// Create and send login request
	loginBody := strings.NewReader(`{"username": "testuser", "password": "testpassword"}`)
	loginReq, _ := http.NewRequest("POST", "/api/users/login", loginBody)
	loginReq.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, loginReq)

	// Verify login was successful
	assert.Equal(t, 200, w.Code)
	expected := `{"message":"Successfully logged in!"}`
	assert.Equal(t, expected, w.Body.String())

	cookies := w.Result().Cookies()

	// Create and send logout request
	logoutReq, _ := http.NewRequest("GET", "/logout", nil)

	// Add cookies to request
	for _, cookie := range cookies {
		logoutReq.AddCookie(cookie)
	}

	w = httptest.NewRecorder()
	router.ServeHTTP(w, logoutReq)

	// Important: Use the same recorder to maintain the session state
	router.ServeHTTP(w, logoutReq)

	// Check the status code is what we expect.
	assert.Equal(t, 200, w.Code)

	// Check the response body is what we expect.
	expected = expected + `{"message":"Successfully logged out!"}`
	assert.Equal(t, expected, w.Body.String())

}

func TestLogoutWithoutAuth(t *testing.T) {
	router := gin.Default()
	initCookies(router)
	router.GET("/logout", logout)
	w := httptest.NewRecorder()

	// Create a request to send to the above route
	req, err := http.NewRequest("GET", "/logout", nil)

	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	// Check the status code is what we expect.
	assert.Equal(t, 200, w.Code)

	// Check the response body is what we expect.
	expected := `{"message":"No session."}`
	assert.Equal(t, expected, w.Body.String())
}
