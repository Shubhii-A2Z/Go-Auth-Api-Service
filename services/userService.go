package services

import (
	db "AuthInGo/db/repositories"
)

type UserService interface {
	CreateUser() error
	GetUserById() error
	GetAllUsers() error
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
	u.userRepository.Create()
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
