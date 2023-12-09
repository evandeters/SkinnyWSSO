package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	router := gin.Default()
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

func TestLogin(t *testing.T) {
	router := gin.Default()
	router.POST("/api/users/login", login)
	w := httptest.NewRecorder()

	// Create a request to send to the above route
	jsonParam := `{"username": "testuser", "password": "testpassword"}`
	req, err := http.NewRequest("POST", "/api/users/login", bytes.NewBufferString(jsonParam))

	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"message":"Successfully logged in!"}`
	assert.Equal(t, expected, w.Body.String())
}

func TestLogoutWithoutAuth(t *testing.T) {
	router := gin.Default()
	router.GET("/logout", logout)
	w := httptest.NewRecorder()

	// Create a request to send to the above route
	req, err := http.NewRequest("GET", "/logout", nil)

	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	// Check the status code is what we expect.
	assert.Equal(t, 200, w.Code)

	// Check the response body is what we expect.
	expected := `{"error":"Unauthorized"}`
	assert.Equal(t, expected, w.Body.String())

}
