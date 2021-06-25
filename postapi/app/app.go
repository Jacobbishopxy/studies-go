package app

import (
	"postapi/app/database"

	"github.com/gorilla/mux"
)

// App 结构体，存储路由结构体与自定义数据库接口类型
type App struct {
	Router *mux.Router
	DB     database.PostDB
}

// App 结构体的构造器
func New() *App {
	a := &App{Router: mux.NewRouter()}

	a.initRoutes()
	return a
}

// App 路由，注册各个路由的具体业务
func (a *App) initRoutes() {
	a.Router.HandleFunc("/", a.IndexHandler()).Methods("GET")
	a.Router.HandleFunc("/api/posts", a.CreatePostHandler()).Methods("POST")
	a.Router.HandleFunc("/api/posts", a.GetPostsHandler()).Methods("GET")
}
