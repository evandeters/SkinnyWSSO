package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	public := router.Group("/")
	addPublicRoutes(public)

	private := router.Group("/")
	private.Use(validateAgainstSSO)
	addPrivateRoutes(private)

	if os.Getenv("USE_HTTPS") == "true" {
		log.Fatalln(router.RunTLS(":4433", os.Getenv("CERT_PATH"), os.Getenv("KEY_PATH")))
	} else {
		log.Fatalln(router.Run(":8080"))
	}
}
