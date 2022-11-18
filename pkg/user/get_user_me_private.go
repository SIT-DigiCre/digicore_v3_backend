package user

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

func GetUserMePrivate(ctx echo.Context, dbClient db.Client) (api.ResGetUserMePrivate, *response.Error) {
	res := api.ResGetUserMePrivate{}
	userId := ctx.Get("user_id").(string)
	private, err := getUserPrivateFromUserId(dbClient, userId)
	if err != nil {
		return api.ResGetUserMePrivate{}, err
	}
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

func getUserPrivateFromUserId(dbClient db.Client, userId string) (private, *response.Error) {
	params := struct {
		UserId string `twowaysql:"userId"`
	}{
		UserId: userId,
	}
	privates := []private{}
	err := dbClient.Select(&privates, "sql/user/select_user_private_from_user_id.sql", &params)
	if err != nil {
		return private{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	if len(privates) == 0 {
		return private{}, &response.Error{Code: http.StatusNotFound, Level: "Info", Message: "個人情報が有りません", Log: "no rows in result"}
	}
	return privates[0], nil
}
