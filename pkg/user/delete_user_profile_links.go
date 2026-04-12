package user

import(
	"net/http"
	
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/labstack/echo/v4"
    "github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
)

func DeleteUserProfileLinks(ctx echo.Context, dbClient db.TransactionClient, id string) (api.ResGetUserUserId, *response.Error) {
	err := DeleteUserProfileLinksFromId(dbClient, id)
	if err != nil {
		return api.ResGetUserUserId{}, err
	}

	return GetUserUserId(ctx, dbClient, ctx.Get("user_id").(string))
}

func DeleteUserProfileLinksFromId(dbClient db.TransactionClient, id string) *response.Error {
	params := struct {
		Id string `twowaysql:"id"`
	}{
		Id: id,
	}
	_, err := dbClient.Exec("sql/user/delete_user_profile_links.sql", &params, false)
	if err != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "リンクの削除に失敗しました", Log: err.Error()}
	}

	return nil
}
