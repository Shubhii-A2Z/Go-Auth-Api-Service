package services

import (
	db "AuthInGo/db/repositories"
	"AuthInGo/utils"
	"fmt"
)

type UserService interface {
	CreateUser() error
	GetUserById() error
	GetAllUsers() error
	LoginUser() error
}

// UserService depending on UserRepo Interface: (ServiceLayer->RepositoryLayer)
type UserServiceImpl struct {
	userRepository db.UserRepository
}

func NewUserService(userRepo db.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		userRepository: userRepo,
	}
}

func (u *UserServiceImpl) CreateUser() error {
	plainPassword:="hashed_password_example"
	hashPass,_:=utils.HashPassword(plainPassword)
	u.userRepository.Create("username_example","user@example.com",hashPass)
	return nil
}

func (u *UserServiceImpl) GetUserById() error {
	u.userRepository.GetById()
	return nil
}

func (u *UserServiceImpl) GetAllUsers() error {
	u.userRepository.GetAll()
	return nil
}

func (u *UserServiceImpl) LoginUser() error {
	resp:=utils.CheckPasswordHash("hashed_password_example1","$2a$10$0qpHsLF4oLm9wtj6OiJcT.eUykowubF2ZW.Hx06ZWWq6Iw7yB/i.2")
	fmt.Println("Login Response:",resp)
	return nil
}
