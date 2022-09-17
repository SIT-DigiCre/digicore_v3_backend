package users

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func PostUserMe(ctx echo.Context, db db.DBClient, requestBody api.ReqPostUserMe) (api.ResPostUserMe, *response.Error) {

	return api.ResPostUserMe{}, nil
}
