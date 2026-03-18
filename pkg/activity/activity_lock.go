package activity

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
)

func lockActivityUser(dbClient db.TransactionClient, userId string) *response.Error {
	params := struct {
		UserId string `twowaysql:"userId"`
	}{
		UserId: userId,
	}

	lockedUser := []struct {
		UserId string `db:"user_id"`
	}{}
	if err := dbClient.Select(&lockedUser, "sql/activity/select_user_profile_for_update_by_user_id.sql", &params); err != nil {
		return &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Info",
			Message: "DBエラーが発生しました",
			Log:     err.Error(),
		}
	}

	return nil
}
