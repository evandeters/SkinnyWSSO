package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"SkinnyWSSO/token"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/go-ldap/ldap/v3"
)

func authRequired(c *gin.Context) {
	session := sessions.Default(c)
	id := session.Get("id")
	if id == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	c.Next()
}

func initCookies(router *gin.Engine) {
	router.Use(sessions.Sessions("kamino", cookie.NewStore([]byte("kamino")))) // change to secret
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
	redirectUrl := c.Request.URL.Query().Encode()

	fmt.Println(redirectUrl)

	// Validate form input
	if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username or password can't be empty."})
		return
	}

	l, err := ldap.DialURL("ldap://ldap:389")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer l.Close()

	err = l.Bind("uid="+username+",ou=users,dc=skinny,dc=wsso", password)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Incorrect username or password."})
		return
	}

	session.Set("id", username)

	prvKey, err := ioutil.ReadFile(os.Getenv("JWT_PRIVATE_KEY"))
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	pubKey, err := ioutil.ReadFile(os.Getenv("JWT_PUBLIC_KEY"))
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	jwtToken := token.NewJWT(prvKey, pubKey)
	tok, err := jwtToken.Create(time.Hour, "auth")
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.SetCookie("token", tok, 86400, "/", "*", false, true)

	if err := session.Save(); err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	c.Redirect(http.StatusFound, redirectUrl)
}

func logout(c *gin.Context) {
	session := sessions.Default(c)
	id := session.Get("id")

	cookie, err := c.Request.Cookie("token")

	if cookie != nil && err == nil {
		c.SetCookie("token", "", -1, "/", "*", false, true)
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

	if err != 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": message})

}

func authFromToken(c *gin.Context) {
	tok := c.Param("token")

	prvKey, err := os.ReadFile(os.Getenv("JWT_PRIVATE_KEY"))
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	pubKey, err := os.ReadFile(os.Getenv("JWT_PUBLIC_KEY"))
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	jwtToken := token.NewJWT(prvKey, pubKey)
	fmt.Println(tok)
	auth, _ := jwtToken.Validate(tok)
	fmt.Println(auth)
	if auth != "auth" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged in!."})
}

func adminAuthRequired(c *gin.Context) int {
	user, password, hasAuth := c.Request.BasicAuth()
	if !hasAuth || (user != os.Getenv("WSSO_ADMIN_USERNAME") && password != os.Getenv("WSSO_ADMIN_PASSWORD")) {
		return 1
	}
	return 0
}
