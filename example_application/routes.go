package main

import (
	"github.com/gin-gonic/gin"
)

func addPublicRoutes(g *gin.RouterGroup) {
	g.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{})
	})
	g.POST("/auth", tokenAuth)
}

func addPrivateRoutes(g *gin.RouterGroup) {
	g.GET("/logout", logout)
	g.GET("/status", status)
}

func status(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}
