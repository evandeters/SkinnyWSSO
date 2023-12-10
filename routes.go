package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

func addPublicRoutes(g *gin.RouterGroup) {
	g.GET("/", viewIndex)
	g.GET("/login/redirect", viewRedirect)
	g.GET("/login", viewLogin)
	g.GET("/register", viewRegister)

	g.GET("/api/users/auth/:token", authFromToken)
	g.POST("/api/users/login", login)
	g.POST("/api/users/register", register)
	g.POST("/api/users/verify", verify)
	g.GET("/api/status", status)
}

func addPrivateRoutes(g *gin.RouterGroup) {
	g.GET("/logout", logout)
	g.GET("/admin", gin.BasicAuth(gin.Accounts{os.Getenv("WSSO_ADMIN_USERNAME"): os.Getenv("WSSO_ADMIN_PASSWORD")}), func(c *gin.Context) {
		c.HTML(200, "admin.html", gin.H{})
	})
	g.GET("/api/users/list", listUsers)
	g.GET("/dashboard", func(c *gin.Context) {
		c.HTML(200, "dashboard.html", gin.H{})
	})
	g.DELETE("/api/users/delete/:username", deleteUser)
}

// View handlers

func viewIndex(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{})
}

func viewRedirect(c *gin.Context) {
	c.HTML(200, "login_redirect.html", gin.H{})
}

func viewLogin(c *gin.Context) {
	c.HTML(200, "login.html", gin.H{})
}

func viewRegister(c *gin.Context) {
	c.HTML(200, "register.html", gin.H{})
}

// API handlers

func status(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}
