package db

import "database/sql"

type UserRepository interface {
	Create() error
}

type UserRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository() *UserRepositoryImpl{
	return &UserRepositoryImpl{}
}

func (u *UserRepositoryImpl) Create() error {
	return nil
}