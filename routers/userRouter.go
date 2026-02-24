package router

import (
	"AuthInGo/controllers"


	"github.com/go-chi/chi/v5"
)

type UserRouter struct {
	UserController controllers.UserController
}

func NewUserRouter(userController controllers.UserController) Router{
	return &UserRouter{
		UserController: userController,
	}
}

func (ur *UserRouter) Register(r chi.Router){
	r.Post("/signup",ur.UserController.RegisterUser)
	r.Get("/profile",ur.UserController.GetUserById)
}