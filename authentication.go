// authentication.go

package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"SkinnyWSSO/token"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/go-ldap/ldap/v3"
)

func initCookies(router *gin.Engine) {
	router.Use(sessions.Sessions("kamino", cookie.NewStore([]byte("kamino")))) // change to secret
}

func logout(c *gin.Context) {
	session := sessions.Default(c)
	id := session.Get("id")

	cookie, err := c.Request.Cookie("auth_token")

	if cookie != nil && err == nil {
		c.SetCookie("auth_token", "", -1, "/", "*", false, true)
	}

	err = session.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	if id == nil {
		c.JSON(http.StatusOK, gin.H{"message": "No session."})
		return
	}
	session.Delete("id")
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out!"})
	c.Redirect(http.StatusSeeOther, "/")
}

func register(c *gin.Context) {
	var jsonData map[string]interface{}
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		fmt.Print(&jsonData)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing fields"})
		return
	}

	username := jsonData["username"].(string)
	password := jsonData["password"].(string)
	email := jsonData["email"].(string)

	message, err := registerUser(username, password, email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account created successfully!"})
}

func login(c *gin.Context) {
	session := sessions.Default(c)
	var jsonData map[string]interface{}
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing fields"})
		return
	}

	username := jsonData["username"].(string)
	password := jsonData["password"].(string)

	userDN := "uid=" + username + ",ou=users,dc=skinny,dc=wsso"

	// Validate form input
	if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username or password can't be empty."})
		return
	}

	// Authenticate user
	l, err := ldap.DialURL("ldap://localhost:389")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer l.Close()

	err = l.Bind(userDN, password)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Incorrect username or password."})
		return
	}

	session.Set("id", username)

	groups, err := GetGroupMembership(username)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT"})
		return
	}
	isAdmin, err := IsMemberOf(username, "admins")
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT"})
		return
	}

	jwtContent := token.UserJWTData{
		Username: username,
		Groups:   groups,
		Admin:    isAdmin,
	}

	tok, err := token.Create(userDN, jwtContent)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT"})
		return
	}

	c.SetCookie("auth_token", tok, 3600, "/", "dev.gfed", false, true)

	if err := session.Save(); err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged in!"})
}

func isLoggedIn(c *gin.Context) (bool, error) {
	session := sessions.Default(c)
	id := session.Get("id")
	if id == nil {
		return false, errors.New("No ID")
	}
	return true, nil
}

func authRequired(c *gin.Context) {
	_, err := isLoggedIn(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	c.Next()
}

func authFromToken(c *gin.Context) {
	tok := c.Param("token")

	claims, err := token.GetClaimsFromToken(tok)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token."})
		return
	}

	if claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged in!"})
}

func isAdmin(c *gin.Context) (bool, error) {

	isLoggedIn, err := isLoggedIn(c)
	if err != nil {
		return false, err
	}

	if isLoggedIn == false {
		return false, errors.New("Not logged in")
	}

	tokenString, err := c.Cookie("auth_token")
	fmt.Println(tokenString)
	if err != nil {
		return false, err
	}

	claims, err := token.GetClaimsFromToken(tokenString)
	if err != nil {
		return false, err
	}

	if val, ok := claims["UserInfo"]; ok {
		userInfo := val.(map[string]interface{})
		if userInfo["admin"] != true {
			return false, errors.New("Not admin")
		}
	} else {
		return false, errors.New("No user info")
	}
	return true, nil
}

func adminAuthRequired(c *gin.Context) {
	isAdmin, err := isAdmin(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if isAdmin == false {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.Next()
}
