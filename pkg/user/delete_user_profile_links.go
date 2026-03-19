package user

import(
    "github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"net/http"
)


func deleteUserProfileLinks(dbClient db.TransactionClient, userId string, linkUrl string) *response.Error {
	params := struct {
		UserId string `twowaysql:"user_id"`
		LinkUrl string `twowaysql:"link_url"`
	}{
		UserId: userId,
		LinkUrl: linkUrl,
	}
	_, err := dbClient.Exec("sql/user/delete_user_profile_links.sql", &params,nil)
	if err != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	return nil
}
