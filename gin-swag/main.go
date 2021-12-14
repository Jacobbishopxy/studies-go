package main

import (
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "gin-swag/docs"
)

// 一个 Todo List 的任务
type todo struct {
	ID   string `json:"id"`
	Task string `json:"task"`
}

// 响应信息
type message struct {
	Message string `json:"message"`
}

// mock 初始数据
var todoList = []todo{
	{"1", "Learn Go"},
	{"2", "Learn Gin"},
	{"3", "Learn Swagger"},
	{"4", "Try to finish this tutorial"},
	{"5", "Have fun with it"},
	{"6", "Do not forget to star it"},
	{"7", "This is a test"},
	{"8", "What a great day"},
	{"9", "I am learning Go"},
	{"10", "Useful for learning Go"},
}

// @Summary get all items in the todo list
// @ID get-all-todos
// @Produce json
// @Success 200 {array} todo
// @Router /todo [get]
func getAllTodos(c *gin.Context) {
	c.JSON(http.StatusOK, todoList)
}

// @Summary get a todo item by ID
// @ID get-todo-by-id
// @Produce json
// @Param id path string true "todo ID"
// @Success 200 {object} todo
// @Failure 404 {object} message
// @Router /todo/{id} [get]
func getTodoByID(c *gin.Context) {
	id := c.Param("id")

	for _, todo := range todoList {
		if todo.ID == id {
			c.JSON(http.StatusOK, todo)
			return
		}
	}

	r := message{"todo not found"}

	c.JSON(http.StatusNotFound, r)
}

// @Summary get a todo item by Pagination
// @ID get-todo-by-pagination
// @Produce json
// @Param offset query string true "todo list offset"
// @Param limit query string true "todo list limit"
// @Success 200 {array} todo
// @Failure 404 {object} message
// @Router /todo_pagination [get]
func getTodoByPagination(c *gin.Context) {
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		c.JSON(http.StatusBadRequest, message{"invalid offset"})
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "5"))
	if err != nil {
		c.JSON(http.StatusBadRequest, message{"invalid limit"})
		return
	}

	if offset <= len(todoList) {
		end := offset + limit
		if end > len(todoList) {
			c.JSON(http.StatusOK, todoList[offset:])
		} else {
			c.JSON(http.StatusOK, todoList[offset:offset+limit])
		}
	} else {
		c.JSON(http.StatusOK, []todo{})
	}
}

// @Summary add a new item to the todo list
// @ID create-todo
// @Produce json
// @Param data body todo true "todo data"
// @Success 200 {object} todo
// @Failure 400 {object} message
// @Router /todo [post]
func createTodo(c *gin.Context) {
	var newTodo todo

	if err := c.BindJSON(&newTodo); err != nil {
		r := message{"invalid todo"}
		c.JSON(http.StatusBadRequest, r)
		return
	}

	todoList = append(todoList, newTodo)
	c.JSON(http.StatusCreated, newTodo)
}

// @Summary delete a todo item by ID
// @ID delete-todo-by-id
// @Produce json
// @Param id path string true "todo ID"
// @Success 200 {object} todo
// @Failure 404 {object} message
// @Router /todo/{id} [delete]
func deleteTodo(c *gin.Context) {
	id := c.Param("id")

	for index, todo := range todoList {
		if todo.ID == id {
			todoList = append(todoList[:index], todoList[index+1:]...)
			c.JSON(http.StatusOK, message{"todo deleted"})
			return
		}
	}

	r := message{"todo not found"}
	c.JSON(http.StatusNotFound, r)
}

// @title Go + Gin Todo API
// @version 1.0
// @description This is a sample server todo server.

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @query.collection.format multi
func main() {

	// Gin logging
	gin.DisableConsoleColor()
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	// 配置 Gin 服务
	router := gin.Default()
	router.GET("/todo", getAllTodos)
	router.GET("/todo/:id", getTodoByID)
	router.GET("/todo_pagination", getTodoByPagination)
	router.POST("/todo", createTodo)
	router.DELETE("/todo", deleteTodo)

	// swagger 文档
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 运行 Gin 服务
	router.Run()
}
