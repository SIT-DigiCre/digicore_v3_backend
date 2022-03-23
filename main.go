package main

import (
	"fmt"
	"os"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/server"
)

// @title Digicore
// @version 3.0
// @description This is digicore backend api

// @host localhost:8000
func main() {
	db, err := server.CreateDbConnection(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE")))
	if err != nil {
		panic(err)
	}

	s := server.CreateEchoServer(db)
	s.Logger.Fatal(s.Start(":8000"))
}
