package main

import (
	"AuthInGo/app"
	config "AuthInGo/config/env"
	dbConfig "AuthInGo/config/db"
	"fmt"
)

func main() {
	config.Load()
	cfg := app.NewConfig()
	app:=app.NewApplication(cfg)

	db,err:=dbConfig.SetupDB()
	if err!=nil {
		fmt.Println("Error connecting database")
		return
	}
	defer db.Close()

	app.Run()
}