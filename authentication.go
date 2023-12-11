// authentication.go

package main

import (
	"fmt"
	"net/http"
	"strings"

	"SkinnyWSSO/token"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/go-ldap/ldap/v3"
)

func getId(c *gin.Context) interface{} {
	session := sessions.Default(c)
	id := session.Get("id")
	return id
}

type jwtData struct {
	Username string   `json:"username"`
	Groups   []string `json:"groups"`
	Admin    bool     `json:"admin"`
}

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

	jwtContent := jwtData{
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

func logout(c *gin.Context) {
	session := sessions.Default(c)
	id := session.Get("id")

	cookie, err := c.Request.Cookie("token")

	if cookie != nil && err == nil {
		c.SetCookie("auth_token", "", -1, "/", "*", false, true)
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

func authFromToken(c *gin.Context) {
	tok := c.Param("token")

	claims, _ := token.GetClaimsFromToken(tok)

	if claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged in!."})
}

func adminAuthRequired(c *gin.Context) {

	session := sessions.Default(c)
	id := session.Get("id")
	if id == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	tokenString, err := c.Cookie("auth_token")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	fmt.Println(tokenString)

	claims, err := token.GetClaimsFromToken(tokenString)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if val, ok := claims["UserInfo"]; ok {
		userInfo := val.(map[string]interface{})
		if userInfo["admin"] != true {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.Next()

}

func AddClaimsToContext(c *gin.Context) {
	session := sessions.Default(c)
	id := session.Get("id")
	if id == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	auth_header := c.GetHeader("Authorization")
	if !strings.HasPrefix(auth_header, "Bearer ") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	tokenString := strings.TrimPrefix(auth_header, "Bearer ")

	claims, err := token.GetClaimsFromToken(tokenString)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	token.SetJWTClaimsContext(c, claims)

	c.Next()
}
