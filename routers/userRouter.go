package router

import (
	"AuthInGo/controllers"
	"AuthInGo/middlewares"

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
	r.With(middlewares.UserCreateRequestValidator).Post("/signup",ur.UserController.RegisterUser)
	r.Get("/profile",ur.UserController.GetUserById)
	r.Get("/profiles",ur.UserController.GetUsers)
	r.With(middlewares.UserLoginRequestValidator).Post("/login",ur.UserController.LoginUser)
}