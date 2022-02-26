package main

import "github.com/SIT-DigiCre/digicore_v3_backend/pkg/server"

// @title Digicore
// @version 3.0
// @description This is digicore backend api

// @host http://localhost:8000
func main() {
	s := server.CreateEchoServer()
	s.Logger.Fatal(s.Start(":8000"))
}
