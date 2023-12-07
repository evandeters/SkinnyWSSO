package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/gin-contrib/sessions"
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

func login(c *gin.Context) {
	session := sessions.Default(c)
	username := c.PostForm("username")
	password := c.PostForm("password")

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
	session.Save()
}

func logout(c *gin.Context) {
	session := sessions.Default(c)
	id := session.Get("id")
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

	matched, _ := regexp.MatchString(`^\w{1,16}$`, username)

	if !matched {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Username must not exceed 16 characters and may only contain letters, numbers, or an underscore (_)!"})
		return
	}

	message, err := registerUser(username, password, os.Getenv("LDAP_PASSWORD"))

	if err != 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": message})

}
