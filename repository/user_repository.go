package repository

import (
	"go-rest-api/model"

	"gorm.io/gorm"
)

// goでは、interfaceは、メソッドの一覧となっている
type IUserRepository interface {
	// koko
	GetUserByEmail(user *model.User, email string) error
	CreateUser(user *model.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	// ポインタをreturnで返している
	return &userRepository{db}
}

// kokoは一緒にならないといけない
func (ur *userRepository) GetUserByEmail(user *model.User, email string) error {
	if err := ur.db.Where("email=?", email).First(user).Error; err != nil {
		return err
	}
	return nil
}

// ポインターレシーバーとして
func (ur *userRepository) CreateUser(user *model.User) error {
	// ポインタを渡している
	if err := ur.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}
