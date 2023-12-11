package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	router.Use(CORSMiddleware())
	router.LoadHTMLGlob("templates/*")
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	public := router.Group("/")
	addPublicRoutes(public)

	private := router.Group("/")
	private.Use(validateAgainstSSO)
	addPrivateRoutes(private)

	if os.Getenv("USE_HTTPS") == "true" {
		log.Fatalln(router.RunTLS(":8080", "/opt/skinnywsso/tls/cert.pem", "/opt/skinnywsso/tls/key.pem"))
	} else {
		log.Fatalln(router.Run(":8080"))
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "https://wsso.gfed.dev")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Origin, X-Requested-With")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		}

		c.Next()
	}
}
