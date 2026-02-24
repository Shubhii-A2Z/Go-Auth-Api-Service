package controllers

import (
	"AuthInGo/services"
	"fmt"
	"net/http"
)

// // UserController depending on UserService Interface: (Controllerlayer->ServiceLayer)
type UserController struct {
	UserService services.UserService
}

func NewUserController(userService services.UserService) *UserController{
	return &UserController{
		UserService: userService,
	}
}

func (uc *UserController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("RegisterUser Controller Called")
	uc.UserService.CreateUser()
}

func (uc *UserController) GetUserById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetUserById Controller Called")
	uc.UserService.GetUserById()
}