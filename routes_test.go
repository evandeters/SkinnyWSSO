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
	assert.Equal(t, "{\"message\":\"Account created successfully!\"}", w.Body.String())
}
