package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
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
	assert.Equal(t, http.StatusOK, w.Code)

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
	assert.Equal(t, http.StatusOK, w.Code)
	expected := `{"message":"Successfully logged in!"}`
	assert.Equal(t, expected, w.Body.String())

	cookies := w.Result().Cookies()

	// Create and send logout request
	logoutReq, _ := http.NewRequest("GET", "/logout", nil)

	// Add cookies to request
	for _, cookie := range cookies {
		logoutReq.AddCookie(cookie)
	}

	// Important: Use the same recorder to maintain the session state
	router.ServeHTTP(w, logoutReq)

	// Check the status code is what we expect.
	assert.Equal(t, http.StatusOK, w.Code)

	// Check the response body is what we expect.
	expected = expected + `{"message":"Successfully logged out!"}`
	assert.Equal(t, expected, w.Body.String())

}

func TestAdminAuthorization(t *testing.T) {

	router := gin.Default()
	initCookies(router) // Make sure this correctly initializes any required middleware
	router.POST("/api/users/login", login)

	// Create and send login request
	loginBody := strings.NewReader(fmt.Sprintf(`{"username": "%s", "password": "%s"}`, os.Getenv("WSSO_ADMIN_USR"), os.Getenv("WSSO_ADMIN_PSW")))
	loginReq, _ := http.NewRequest("POST", "/api/users/login", loginBody)
	loginReq.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, loginReq)

	// Verify login was successful
	assert.Equal(t, http.StatusOK, w.Code)
	expected := `{"message":"Successfully logged in!"}`
	assert.Equal(t, expected, w.Body.String())

	cookies := w.Result().Cookies()

	// Create and send User List Request (Requires Admin)
	adminReq, _ := http.NewRequest("GET", "/api/users/list", nil)

	// Add cookies to request
	for _, cookie := range cookies {
		adminReq.AddCookie(cookie)
	}

	router.ServeHTTP(w, adminReq)

	assert.Equal(t, http.StatusOK, w.Code)

	// Now try to access the admin page with a non-admin user

}

func TestFailedAdminAuthorization(t *testing.T) {

	router := gin.Default()
	initCookies(router) // Make sure this correctly initializes any required middleware
	router.POST("/api/users/login", login)
	router.GET("/api/users/list", listUsers)

	loginBody := strings.NewReader(`{"username": "testuser", "password": "testpassword"}`)
	loginReq, _ := http.NewRequest("POST", "/api/users/login", loginBody)
	loginReq.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, loginReq)

	// Verify login was successful
	assert.Equal(t, http.StatusOK, w.Code)
	expected := `{"message":"Successfully logged in!"}`
	assert.Equal(t, expected, w.Body.String())

	cookies := w.Result().Cookies()

	adminReq, _ := http.NewRequest("GET", "/api/users/list", nil)

	for _, cookie := range cookies {
		fmt.Println(cookie)
		adminReq.AddCookie(cookie)
	}

	router.ServeHTTP(w, adminReq)

	expected = expected + `{"error":"Unauthorized."}`
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, expected, w.Body.String())
}

func TestDeleteUser(t *testing.T) {

	router := gin.Default()
	initCookies(router) // Make sure this correctly initializes any required middleware
	router.POST("/api/users/login", login)
	router.DELETE("/api/users/delete/:username", deleteUser)

	// Create and send login request
	loginBody := strings.NewReader(fmt.Sprintf(`{"username": "%s", "password": "%s"}`, os.Getenv("WSSO_ADMIN_USR"), os.Getenv("WSSO_ADMIN_PSW")))
	loginReq, _ := http.NewRequest("POST", "/api/users/login", loginBody)
	loginReq.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, loginReq)

	// Verify login was successful
	assert.Equal(t, http.StatusOK, w.Code)
	expected := `{"message":"Successfully logged in!"}`
	assert.Equal(t, expected, w.Body.String())

	cookies := w.Result().Cookies()

	// Create Delete Request
	deleteReq, _ := http.NewRequest("DELETE", "/api/users/delete/testuser", nil)

	// Add cookies to request
	for _, cookie := range cookies {
		deleteReq.AddCookie(cookie)
	}

	router.ServeHTTP(w, deleteReq)

	assert.Equal(t, http.StatusOK, w.Code)

	expected = expected + `{"message":"Account deleted successfully!"}`
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
	assert.Equal(t, http.StatusOK, w.Code)

	// Check the response body is what we expect.
	expected := `{"message":"No session."}`
	assert.Equal(t, expected, w.Body.String())
}
