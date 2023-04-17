package status

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/labstack/echo/v4"
)

func GetStatus(ctx echo.Context) (api.ResGetStatus, *response.Error) {
	// return api.Status{}, &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "error", Log: "error"}
	return api.ResGetStatus{Status: true}, nil
}
