package admin

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func DeleteExpiredUserPrivateProfiles(_ echo.Context, dbClient db.TransactionClient) (api.BlankSuccess, *response.Error) {
	_, err := dbClient.Exec("sql/admin/delete_expired_user_private_profiles.sql", nil, false)
	if err != nil {
		return api.BlankSuccess{}, &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Error",
			Message: "個人情報の削除に失敗しました",
			Log:     err.Error(),
		}
	}

	return api.BlankSuccess{Success: true}, nil
}
