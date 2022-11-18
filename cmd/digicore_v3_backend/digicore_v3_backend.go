package main

import "github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/server"

func main() {
	e := server.CreateEchoServer()
	e.Logger.Fatal(e.Start(":8000"))
}
