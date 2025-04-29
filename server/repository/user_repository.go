package repository

import (
	"tablelink_project/server/model"
	"time"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(model.User) error
	UpdateUser(*model.User) error
	DeleteUser(userID int) error
	GetUserByEmail(email string) (*model.User, error)
	GetUserByID(userID int) (*model.User, error)
	GetAllUsers() ([]model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (ur *userRepository) CreateUser(user model.User) error {
	user.UpdatedAt = time.Now()
	user.CreatedAt = time.Now()
	err := ur.db.Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) UpdateUser(user *model.User) error {
	user.UpdatedAt = time.Now()
	err := ur.db.Model(&model.User{}).Where("id = ?", user.ID).Updates(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) DeleteUser(userID int) error {
	err := ur.db.Delete(&model.User{}, userID).Error
	if err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) GetUserByEmail(email string) (*model.User, error) {
	user := &model.User{}
	err := ur.db.Model(model.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *userRepository) GetUserByID(userID int) (*model.User, error) {
	user := &model.User{}
	err := ur.db.First(&user, userID).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *userRepository) GetAllUsers() ([]model.User, error) {
	user := []model.User{}
	err := ur.db.Model(model.User{}).Take(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}
