package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("templates/*.html")
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	initCookies(router)

	public := router.Group("/")
	addPublicRoutes(public)

	private := router.Group("/")
	private.Use(authRequired)
	addPrivateRoutes(private)

	admin := router.Group("/")
	admin.Use(adminAuthRequired)
	addAdminRoutes(admin)

	if os.Getenv("USE_HTTPS") == "true" {
		log.Fatalln(router.RunTLS(":443", os.Getenv("CERT_PATH"), os.Getenv("KEY_PATH")))
	} else {
		log.Fatalln(router.Run(":80"))
	}
}
