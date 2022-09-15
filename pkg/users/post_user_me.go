package users

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func PostUserMe(ctx echo.Context, db db.DBClient) (api.ResPostUserMe, *response.Error) {
	var req api.ReqPostUserMe
	ctx.Bind(&req)
	return api.ResPostUserMe{}, nil
}
