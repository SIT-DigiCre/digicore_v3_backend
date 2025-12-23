package activity

import (
	"net/http"
	"time"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func PostActivityCheckin(ctx echo.Context, dbClient db.TransactionClient, requestBody api.ReqPostActivityCheckin) (api.BlankSuccess, *response.Error) {
	userId := ctx.Get("user_id").(string)

	checkInAt := time.Now()

	// 既存の在室レコードがあればチェックアウトしてから新規レコードを作成
	_, err := executeCheckout(dbClient, userId, requestBody.Place, checkInAt)
	if err != nil {
		return api.BlankSuccess{}, err
	}

	if err := executeCheckin(dbClient, userId, requestBody.Place, checkInAt); err != nil {
		return api.BlankSuccess{}, err
	}

	return api.BlankSuccess{Success: true}, nil
}

func executeCheckin(dbClient db.TransactionClient, userId string, place string, checkedInAt time.Time) *response.Error {
	params := struct {
		UserId             string    `twowaysql:"userId"`
		Place              string    `twowaysql:"place"`
		InitialCheckedInAt time.Time `twowaysql:"initialCheckedInAt"`
		CheckedInAt        time.Time `twowaysql:"checkedInAt"`
	}{
		UserId:             userId,
		Place:              place,
		InitialCheckedInAt: checkedInAt,
		CheckedInAt:        checkedInAt,
	}

	_, err := dbClient.Exec("sql/activity/insert_activity.sql", &params, true)
	if err != nil {
		return &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Info",
			Message: "DBエラーが発生しました",
			Log:     err.Error(),
		}
	}
	return nil
}
