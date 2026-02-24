package app

import (
	config "AuthInGo/config/env"
	dbConfig "AuthInGo/config/db"
	"AuthInGo/controllers"
	repo "AuthInGo/db/repositories"
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
	}
}

func (app *Application) Run() error{

	db,err:=dbConfig.SetupDB()
	if err!=nil {
		fmt.Println("Error Connecting Database",err)
		return err
	}

	userRepo:=repo.NewUserRepository(db)
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