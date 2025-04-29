package service

import (
	"html"
	"strings"
	"tablelink_project/server/model"
	"tablelink_project/server/repository"
	"tablelink_project/server/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(model.User) error
	UpdateUser(*model.User) error
	DeleteUser(userID int) error
	LoginCheck(username, password string) (string, error)
	GetUserByID(userID int) (*model.User, error)
	GetAllUsers() ([]model.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (us *userService) CreateUser(user model.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	user.Email = html.EscapeString(strings.TrimSpace(user.Email))

	return us.userRepo.CreateUser(user)
}

func (us *userService) UpdateUser(user *model.User) error {
	return us.userRepo.UpdateUser(user)
}

func (us *userService) DeleteUser(userID int) error {
	return us.userRepo.DeleteUser(userID)
}

func (us *userService) LoginCheck(email, password string) (string, error) {
	var err error

	user, err := us.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", err
	}

	err = verifyPassword(password, user.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	token, err := utils.GenerateToken(user.ID)

	if err != nil {
		return "", err
	}

	user.LastAccess = time.Now()
	err = us.userRepo.UpdateUser(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (us *userService) GetUserByID(userID int) (*model.User, error) {
	return us.userRepo.GetUserByID(userID)
}

func (us *userService) GetAllUsers() ([]model.User, error) {
	return us.userRepo.GetAllUsers()
}

func verifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
