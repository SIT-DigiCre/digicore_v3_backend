package user

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/grant"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/group"
	"github.com/labstack/echo/v4"
)

func GetUserMeGrants(ctx echo.Context, dbClient db.Client) (api.ResGetUserMeGrants, *response.Error) {
	userId := ctx.Get("user_id").(string)

	claims, err := group.GetClaimsFromUserId(dbClient, userId)
	if err != nil {
		return api.ResGetUserMeGrants{}, err
	}

	return api.ResGetUserMeGrants{
		Grants: grant.ResolveFromClaims(claims),
	}, nil
}
