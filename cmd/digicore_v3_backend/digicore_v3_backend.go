package main

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/server"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/batch"
)

func main() {
	batch.Run()
	e := server.CreateEchoServer()
	e.Logger.Fatal(e.Start(":8000"))
}
