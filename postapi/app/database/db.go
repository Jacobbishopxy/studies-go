package database

import (
	"log"
	"postapi/app/models"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// 数据库的业务接口，需要实现数据库连接与断开，以及业务逻辑
type PostDB interface {
	Open() error
	Close() error
	CreatePost(p *models.Post) error
	GetPosts() ([]*models.Post, error)
}

// 数据库结构体，包装三方库的结构体
type DB struct {
	db *sqlx.DB
}

// 数据库连接，初始化表
func (d *DB) Open() error {
	pg, err := sqlx.Open("postgres", pgConnStr)
	if err != nil {
		return err
	}
	log.Println("Connected to Database!")

	pg.MustExec(createSchema)

	// 为 DB 结构体赋值
	d.db = pg

	return nil
}

// 数据库断开
func (d *DB) Close() error {
	return d.db.Close()
}
