package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	db "fintech-banking-app/db/sqlc"
	"fintech-banking-app/token"
	"fintech-banking-app/util"
)

// 定义 server 结构体
// `util.Config` 加载 env 文件中的变量
// `db.Store` 即之前所编写的实现。在处理客户端的 API 请求时，它允许我们与进行数据库交互。
// `token.Maker` 加密的实现
// `gin.Enine` 将帮助我们发送每个 API 请求至正确的函数进行处理加工。
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// 入参 store 为数据库 interface，即真实 DB 与 Mock DB 所实现的接口
func NewServer(config util.Config, store db.Store) (*Server, error) {

	// 加密
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmeticKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{store: store, tokenMaker: tokenMaker}

	// 将自定义 validator 引入 Gin 的 validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册 currency 的 validator
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	// 传入一个或多个 handler 函数。
	// 只有最后的函数为 handler，前面部分的函数皆为 middlewares
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)
	router.POST("/transfers", server.createTransfer)

	server.router = router
}

func (server *Server) Start(address string) error {
	// `server.router` 为私有字段，不可被 `api` package 外所访问
	return server.router.Run(address)
}

// 接受 error 作为入参，返回 `gin.H` 对象（`map[string]interface{}`的简写）。
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
