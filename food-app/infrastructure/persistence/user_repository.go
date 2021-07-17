package persistence

import (
	"errors"
	"strings"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"

	"food-app/domain/entity"
	"food-app/domain/repository"
	"food-app/infrastructure/security"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepo {
	return &UserRepo{db}
}

// UserRepo 实现 repository.UserRepository 接口
var _ repository.UserRepository = &UserRepo{}

func (r *UserRepo) SaveUser(user *entity.User) (*entity.User, map[string]string) {
	dbErr := map[string]string{}
	err := r.db.Debug().Create(&user).Error
	if err != nil {
		// 如果 email 已被占用
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicated") {
			dbErr["email_take"] = "email already taken"
			return nil, dbErr
		}
		// 其它 db 错误
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return user, nil
}

func (r *UserRepo) GetUser(id uint64) (*entity.User, error) {
	var user entity.User
	err := r.db.Debug().Where("id = ?", id).Take(&user).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (r *UserRepo) GetUsers() ([]entity.User, error) {
	var users []entity.User
	err := r.db.Debug().Find(&users).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("user not found")
	}
	return users, nil
}

func (r *UserRepo) GetUserByEmailAndPassword(u *entity.User) (*entity.User, map[string]string) {
	var user entity.User
	dbErr := make(map[string]string)
	err := r.db.Debug().Where("email = ?", u.Email).Take(&user).Error
	if gorm.IsRecordNotFoundError(err) {
		dbErr["no_user"] = "user not found"
		return nil, dbErr
	}
	if err != nil {
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	// 验证密码
	err = security.VerifyPassword(user.Password, u.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		dbErr["incorrect_password"] = "incorrect password"
		return nil, dbErr
	}

	return &user, nil
}
