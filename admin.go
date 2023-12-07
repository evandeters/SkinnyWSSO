package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

func listUsers(c *gin.Context) {
	users, err := getLdapUsers(os.Getenv("LDAP_ADMIN_PASSWORD"))
	if err != 0 {
		c.AbortWithStatusJSON(500, gin.H{"error": "Failed to get users."})
		return
	}
	c.JSON(200, gin.H{"users": users})
}

func deleteUser(c *gin.Context) {
	username := c.Param("username")
	message, err := deleteLdapUser(username, os.Getenv("LDAP_ADMIN_PASSWORD"))
	if err != 0 {
		c.JSON(500, gin.H{"error": message})
		return
	}
	c.JSON(200, gin.H{"message": message})
}
