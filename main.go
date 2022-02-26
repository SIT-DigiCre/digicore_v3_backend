package main

import "github.com/SIT-DigiCre/digicore_v3_backend/pkg/server"

func main() {
	s := server.CreateEchoServer()
	s.Logger.Fatal(s.Start(":8000"))
}
