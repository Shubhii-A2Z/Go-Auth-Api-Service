package app

import (
	config "AuthInGo/config/env"
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

	server:=&http.Server{
		Addr: app.Config.Addr,
		Handler: nil,
		ReadTimeout: 10*time.Second,
		WriteTimeout: 10*time.Second,
	}
	
	fmt.Println("Server started at PORT",app.Config.Addr)
	return server.ListenAndServe()
}