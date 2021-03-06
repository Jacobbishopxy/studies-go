package main

import (
	"blog/api/controller"
	"blog/api/repository"
	"blog/api/routes"
	"blog/api/service"
	"blog/infrastructure"
	"blog/models"
)

func main() {
	// 加载 env
	infrastructure.LoadEnv()

	router := infrastructure.NewGinRouter()
	db := infrastructure.NewDatabase()
	postRepository := repository.NewPostRepository(&db)
	postService := service.NewPostService(&postRepository)
	postController := controller.NewPostController(&postService)
	postRoute := routes.NewPostRoute(&postController, &router)

	postRoute.Setup()

	db.DB.AutoMigrate(&models.Post{})

	router.Gin.Run(":8080")
}
