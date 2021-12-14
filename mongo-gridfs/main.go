package main

import (
	"github.com/gin-gonic/gin"

	"mongo-gridfs/service"
)

func main() {

	router := gin.Default()
	router.POST("/upload", service.FileUpload)
	router.GET("/download", service.FileDownload)

	router.Run()
}
