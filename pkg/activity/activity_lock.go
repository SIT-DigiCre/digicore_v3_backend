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

	// activity の checkin / checkout は「そのユーザーの現在在室中レコード」を更新するため、
	// 同一ユーザーに対する並行リクエストだけを直列化できれば十分です。
	// user_profiles は user_id ごとに一意なので、この行を FOR UPDATE で掴んで
	// activities 更新の前にユーザー単位の排他制御をかけています。
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
