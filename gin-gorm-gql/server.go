package main

import (
	"fmt"
	"gin-gorm-gql/graph/model"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"gin-gorm-gql/graph/handlers"
)

const defaultPort = "8080"
const dbname = "dev"

var db *gorm.DB

func initDB() *gorm.DB {
	var err error
	dsn := fmt.Sprintf("host=localhost user=root password=secret dbname=%s port=5432 sslmode=disable TimeZone=Asia/Shanghai", dbname)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}

	db.Debug()
	// db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbname))
	// db.Exec(fmt.Sprintf("USE %s", dbname))

	db.AutoMigrate(&model.Order{}, &model.Item{})
	return db
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	initDB()

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.POST("/api/v1/", handlers.GraphqlHandler(db))
	r.GET("/", handlers.PlaygroundHandler())

	r.GET("/js/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "js.html", nil)
	})

	r.Run(":" + defaultPort)
}
