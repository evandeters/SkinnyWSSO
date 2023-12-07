package main

import (
	"github.com/gin-gonic/gin"
)

func listUsers(c *gin.Context) {
	if adminAuthRequired(c) != 0 {
		c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	users, err := getLdapUsers()
	if err != 0 {
		c.AbortWithStatusJSON(500, gin.H{"error": "Failed to get users."})
		return
	}
	c.JSON(200, gin.H{"users": users})
}

func deleteUser(c *gin.Context) {
	if adminAuthRequired(c) != 0 {
		c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	username := c.Param("username")
	message, err := deleteLdapUser(username)
	if err != 0 {
		c.JSON(500, gin.H{"error": message})
		return
	}
	c.JSON(200, gin.H{"message": message})
}
