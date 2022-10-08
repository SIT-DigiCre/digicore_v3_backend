package users

import (
	"database/sql"
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

func GetUserMePrivate(ctx echo.Context, dbClient db.Client) (api.ResGetUserMePrivate, *response.Error) {
	userID := ctx.Get("user_id").(string)
	private, err := getUserPrivateFromUserID(userID, dbClient)
	if err != nil {
		return api.ResGetUserMePrivate{}, err
	}
	res := api.ResGetUserMePrivate{}
	rerr := copier.Copy(&res, &private)
	if rerr != nil {
		return api.ResGetUserMePrivate{}, &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "プロフィールの読み込みに失敗しました", Log: rerr.Error()}
	}
	return res, nil
}

type private struct {
	FirstName             string `db:"first_name"`
	LastName              string `db:"last_name"`
	FirstNameKana         string `db:"first_name_kana"`
	LastNameKana          string `db:"last_name_kana"`
	IsMale                bool   `db:"is_male"`
	PhoneNumber           string `db:"phone_number"`
	Address               string `db:"address"`
	ParentName            string `db:"parent_name"`
	ParentCellphoneNumber string `db:"parent_cellphone_number"`
	ParentHomephoneNumber string `db:"parent_homephone_number"`
	ParentAddress         string `db:"parent_address"`
}

func getUserPrivateFromUserID(userID string, dbClient db.Client) (private, *response.Error) {
	params := struct {
		UserID string `twowaysql:"userID"`
	}{
		UserID: userID,
	}
	privates := []private{}
	err := dbClient.Select(&privates, "sql/users/select_user_private_from_user_id.sql", &params)
	if err == sql.ErrNoRows {
		return private{}, &response.Error{Code: http.StatusNotFound, Level: "Info", Message: "プロフィールが有りません", Log: sql.ErrNoRows.Error()}
	}
	if err != nil {
		return private{}, &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: err.Error()}
	}
	return privates[0], nil
}
