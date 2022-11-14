package work

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func GetWorkWork(ctx echo.Context, dbClient db.Client, params api.GetWorkWorkParams) (api.ResGetWorkWork, *response.Error) {

	return api.ResGetWorkWork{}, nil
}
