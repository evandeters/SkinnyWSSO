package main

import (
	"github.com/gin-gonic/gin"
)

func listUsers(c *gin.Context) {

	users, err := getLdapUsers()
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "Failed to get users."})
		return
	}
	c.JSON(200, gin.H{"users": users})
}

func deleteUser(c *gin.Context) {

	username := c.Param("username")
	message, err := deleteLdapUser(username)
	if err != nil {
		c.JSON(500, gin.H{"error": message})
		return
	}
	c.JSON(200, gin.H{"message": message})
}
