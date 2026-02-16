package app

import (
	config "AuthInGo/config/env"
	"AuthInGo/controllers"
	db "AuthInGo/db/repositories"
	router "AuthInGo/routers"
	"AuthInGo/services"
	"fmt"
	"net/http"
	"time"
)

type Config struct {
	Addr string
}

type Application struct {
	Config Config
	Store db.Storage
}

// Constructor for Config
func NewConfig() Config {
	port := config.GetString("PORT",":8080")
	return Config{
		Addr: port,
	}
}

// Constructor for Application
func NewApplication(cfg Config) *Application{
	return &Application{
		Config: cfg,
		Store: *db.NewStorage(),
	}
}

func (app *Application) Run() error{

	userRepo:=db.NewUserRepository()
	userService:=services.NewUserService(userRepo)
	userController:=controllers.NewUserController(userService)
	userRouter:=router.NewUserRouter(*userController)

	server:=&http.Server{
		Addr: app.Config.Addr,
		Handler: router.SetUpRouter(userRouter),
		ReadTimeout: 10*time.Second,
		WriteTimeout: 10*time.Second,
	}
	
	fmt.Println("Server started at PORT",app.Config.Addr)
	return server.ListenAndServe()
}