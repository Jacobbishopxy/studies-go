package routes

import (
	"blog/api/controller"
	"blog/infrastructure"
)

// api 最终暴露的 Route
type PostRoute struct {
	Controller *controller.PostController
	Handler    *infrastructure.GinRouter
}

func NewPostRoute(
	controller *controller.PostController,
	handler *infrastructure.GinRouter,
) PostRoute {
	return PostRoute{controller, handler}
}

func (p *PostRoute) Setup() {
	// Router 组
	post := p.Handler.Gin.Group("/posts")
	{
		post.GET("/", p.Controller.GetPosts)
		post.POST("/", p.Controller.AddPost)
		post.GET("/:id", p.Controller.GetPost)
		post.DELETE("/:id", p.Controller.DeletePost)
		post.PUT("/:id", p.Controller.UpdatePost)
	}

}
