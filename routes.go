package main

import (
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
}

func health(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}
