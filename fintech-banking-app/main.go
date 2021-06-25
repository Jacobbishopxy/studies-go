package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" // 必须引入 driver 才能连接数据库

	"fintech-banking-app/api"
	db "fintech-banking-app/db/sqlc"
	"fintech-banking-app/util"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("connot start server:", err)
	}
}
