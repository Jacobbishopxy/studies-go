package application

import (
	"food-app/domain/entity"
	"food-app/domain/repository"
)

type userApp struct {
	us repository.UserRepository
}

var _ UserAppInterface = &userApp{}

type UserAppInterface interface {
	SaveUser(*entity.User) (*entity.User, map[string]string)
	GetUsers() ([]entity.User, error)
	GetUser(uint64) (*entity.User, error)
	GetUserByEmailAndPassword(*entity.User) (*entity.User, map[string]string)
}

func (u *userApp) SaveUser(user *entity.User) (*entity.User, map[string]string) {
	return u.us.SaveUser(user)
}

func (u *userApp) GetUsers() ([]entity.User, error) {
	return u.us.GetUsers()
}

func (u *userApp) GetUser(userId uint64) (*entity.User, error) {
	return u.us.GetUser(userId)
}

func (u *userApp) GetUserByEmailAndPassword(user *entity.User) (*entity.User, map[string]string) {
	return u.us.GetUserByEmailAndPassword(user)
}
