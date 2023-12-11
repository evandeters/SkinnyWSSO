package main

import (
	"net/http"

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
	g.GET("/dashboard", viewDashboard)
	g.GET("/logout", logout)
}

func addAdminRoutes(g *gin.RouterGroup) {
	g.GET("/admin", viewAdmin)

	g.GET("/api/users/list", listUsers)
	g.DELETE("/api/users/delete/:username", deleteUser)
}

// Misc functions

func pageData(c *gin.Context, specialData gin.H) gin.H {
	var data gin.H
	data["isAdmin"] = false

	// iterate over all keys in specialData and add them to data
	for key, value := range specialData {
		data[key] = value
	}
	return data
}

// View handlers

func viewIndex(c *gin.Context) {
	isLoggedIn, err := isLoggedIn(c)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", pageData(c, gin.H{"error": err}))
	}

	if isLoggedIn == false {
		c.HTML(http.StatusOK, "index.html", pageData(c, gin.H{}))
	} else {
		c.Redirect(http.StatusFound, "/dashboard")
	}
}

func viewDashboard(c *gin.Context) {
	c.HTML(http.StatusOK, "dashboard.html", pageData(c, gin.H{}))
}

func viewRedirect(c *gin.Context) {
	c.HTML(http.StatusOK, "login_redirect.html", pageData(c, gin.H{}))
}

func viewLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", pageData(c, gin.H{}))
}

func viewRegister(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", pageData(c, gin.H{}))
}

func viewAdmin(c *gin.Context) {
	c.HTML(http.StatusOK, "admin.html", pageData(c, gin.H{}))
}

// API handlers

func status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
