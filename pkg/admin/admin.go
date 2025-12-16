package admin

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
)

type SelectClient interface {
	Select(dest interface{}, queryPath string, params interface{}) error
}

func CheckUserIsAdmin(dbClient SelectClient, userId string) (bool, *response.Error) {
	params := struct {
		UserId string `twowaysql:"userId"`
	}{
		UserId: userId,
	}

	result := []struct {
		IsAdmin bool `db:"is_admin"`
	}{}

	err := dbClient.Select(&result, "sql/admin/select_user_is_admin.sql", &params)
	if err != nil {
		return false, &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Info",
			Message: "DBエラーが発生しました",
			Log:     err.Error(),
		}
	}

	if len(result) == 0 {
		return false, nil
	}

	return result[0].IsAdmin, nil
}
