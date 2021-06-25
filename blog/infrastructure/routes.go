package infrastructure

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 主 Router
type GinRouter struct {
	Gin *gin.Engine
}

func NewGinRouter() GinRouter {

	httpRouter := gin.Default()

	httpRouter.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"data": "Up and Running..."})
	})
	return GinRouter{Gin: httpRouter}
}
