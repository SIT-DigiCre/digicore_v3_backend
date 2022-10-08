package response

import (
	"fmt"
	"runtime/debug"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type Error struct {
	Code    int
	Level   string
	Message string
	Log     string
}

func ErrorResponse(ctx echo.Context, err *Error) error {
	switch err.Level {
	case "Info":
		logrus.Info(err.Log)
	case "Warn":
		logrus.Warn(err.Log)
	case "Error":
		logrus.Error(err.Log)
		fmt.Print("!!!!!!!!!!!!!!!!!!!!!!!")
		debug.PrintStack()
	default:
		logrus.Debug(err.Log)
	}
	return ctx.JSON(err.Code, api.Error{Level: err.Level, Message: err.Message})
}
