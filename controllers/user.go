package controllers

import (
	"AuthInGo/services"
	"net/http"
)

type UserController struct {
	UserService services.UserService
}

func NewUserController(userService services.UserService) *UserController{
	return &UserController{
		UserService: userService,
	}
}

func (uc *UserController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("User Registered"))
}