package user

	import (
		"net/http"

		"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
		"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
		"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	)


func GetPublicMemberCount(dbClient db.Client) (api.ResGetPublicMemberCount, *response.Error) {
	memberCounts := []api.ResGetPublicMemberCount{}
	err := dbClient.Select(&memberCounts, "sql/user/select_member_count.sql", nil)
	if err != nil {
		return api.ResGetPublicMemberCount{}, &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Error",
			Message: "メンバー数の取得に失敗しました",
			Log:     err.Error(),
		}
	}

	if len(memberCounts) == 0 {
		return api.ResGetPublicMemberCount{}, &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Error",
			Message: "メンバー数の取得に失敗しました",
			Log:     "no rows in result",
		}
	}

	res := memberCounts[0]
	return res, nil
}
