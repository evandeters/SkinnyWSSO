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
	g.POST("/api/users/login", login)
	g.POST("/api/users/register", register)
	g.GET("/api/status", status)
}

func addPrivateRoutes(g *gin.RouterGroup) {
	g.GET("/logout", logout)
	g.GET("/admin", gin.BasicAuth(gin.Accounts{os.Getenv("WSSO_ADMIN_USERNAME"): os.Getenv("WSSO_ADMIN_PASSWORD")}), func(c *gin.Context) {
		c.HTML(200, "admin.html", gin.H{})
	})
	g.GET("/api/users/list", listUsers)
	g.DELETE("/api/users/delete/:username", deleteUser)
}

func status(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}
