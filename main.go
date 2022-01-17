package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/smarest/smarest-paste/application"
)

func main() {
	bean, err := application.InitBean()
	if err != nil {
		log.Fatalln("can not create bean", err)
	}
	router := gin.Default()
	router.GET("/cookie/:uid", bean.FileUploadService.SetCookie)

	router.Static("/paste/view", bean.FileUploadService.Directory)
	router.POST("/paste", bean.FileUploadService.Post)
	router.Run(":8080")
}
