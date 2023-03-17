package user

import (
	"strconv"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/util"
	"github.com/labstack/echo/v4"
)

func PutUserMeRenewal(ctx echo.Context, dbClient db.TransactionClient) (api.ResGetUserMe, *response.Error) {
	userId := ctx.Get("user_id").(string)
	err := util.RenewalActiveLimit(dbClient, userId, strconv.Itoa(util.GetYear())+"-06-01")
	if err != nil {
		return api.ResGetUserMe{}, err
	}
	return GetUserMe(ctx, dbClient)
}
