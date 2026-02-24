package config

import (
	env "AuthInGo/config/env"
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

func SetupDB() (*sql.DB,error) {
	cfg := mysql.NewConfig()

	cfg.User=env.GetString("DB_USER","root")
	cfg.Passwd=env.GetString("DB_PASS","root")
	cfg.Net=env.GetString("DB_NET","tcp")
	cfg.Addr=env.GetString("DB_ADDR","127.0.0.1:3306")
	cfg.DBName=env.GetString("DB_NAME","auth_dev")
	cfg.ParseTime=true

	fmt.Println("Connecting to database:",cfg.FormatDSN())

	// validating the DSN: opens a database specified by its database driver name and a driver-specific data source name
	db,err:=sql.Open("mysql",cfg.FormatDSN())

	if(err!=nil){
		fmt.Println("Error connecting database",err)
		return nil,err
	}

	// verifying the connection to the database is still alive
	pingErr:=db.Ping()
	
	if(pingErr!=nil){
		fmt.Println("Error Pinging database",pingErr)
		return nil,pingErr
	}

	return db,nil
}