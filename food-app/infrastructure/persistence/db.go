package persistence

import (
	"fmt"
	"food-app/domain/entity"
	"food-app/domain/repository"

	"github.com/jinzhu/gorm"
)

type Repositories struct {
	User repository.UserRepository
	Food repository.FoodRepository
	db   *gorm.DB
}

func NewRepositories(DbDriver, DbUser, DbPassword, DbHost, DbPort, DbName string) (*Repositories, error) {
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", DbHost, DbPort, DbUser, DbName, DbPassword)
	db, err := gorm.Open(DbDriver, DBURL)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)

	return &Repositories{
		User: NewUserRepository(db),
		Food: NewFoodRepository(db),
		db:   db,
	}, nil
}

// 关闭数据库连接
func (s *Repositories) Close() error {
	return s.db.Close()
}

// 迁移所有表
func (s *Repositories) Automigrate() error {
	return s.db.AutoMigrate(&entity.User{}, &entity.Food{}).Error
}
