package main

import (
	"strings"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/server"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/batch"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/env"
	"github.com/sirupsen/logrus"
)

func main() {
	if missing := env.Missing(); len(missing) > 0 {
		logrus.Warnf("未設定の環境変数: %s", strings.Join(missing, ", "))
	}
	batch.Run()
	e := server.CreateEchoServer()
	e.Logger.Fatal(e.Start(":8000"))
}
