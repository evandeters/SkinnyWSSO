package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

func addPublicRoutes(g *gin.RouterGroup) {
	g.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{})
	})
	g.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", gin.H{})
	})
	g.GET("/register", func(c *gin.Context) {
		c.HTML(200, "register.html", gin.H{})
	})
	g.POST("/login", login)
	g.POST("/register", register)
	g.GET("/health", health)
}

func addPrivateRoutes(g *gin.RouterGroup) {
	g.GET("/logout", logout)
	g.GET("/admin", gin.BasicAuth(gin.Accounts{"admin": os.Getenv("WSSO_ADMIN_PASSWORD")}), func(c *gin.Context) {
		c.HTML(200, "admin.html", gin.H{})
	})
	g.GET("/api/users/list", listUsers)
	g.DELETE("/api/users/delete/:username", deleteUser)
}

func health(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}
