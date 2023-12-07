package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	public := router.Group("/")
	addPublicRoutes(public)

	private := router.Group("/")
	private.Use(authRequired)
	addPrivateRoutes(private)

	if os.Getenv("USE_HTTPS") == "true" {
		log.Fatalln(router.RunTLS(":443", os.Getenv("CERT_PATH"), os.Getenv("KEY_PATH")))
	} else {
		log.Fatalln(router.Run(":80"))
	}
}
