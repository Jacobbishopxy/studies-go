package infrastructure

import (
	"log"

	"github.com/joho/godotenv"
)

// 从 .env 文件中加载 env 变量
func LoadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("unable to load .env file")
	}
}
