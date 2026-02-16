package services

import db "AuthInGo/db/repositories"

type UserService interface {
	CreateUser() error
}

// UserService depending on UserRepo Interface
type UserServiceImpl struct {
	userRepository db.UserRepository
}

func NewUserService(userRepo db.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		userRepository: userRepo,
	}
}

func (u *UserServiceImpl) CreateUser() error {
	return nil
}
