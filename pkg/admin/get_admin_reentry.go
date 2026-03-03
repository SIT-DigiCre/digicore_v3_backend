package admin

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

func GetAdminReentry(_ echo.Context, dbClient db.Client) (api.ResGetAdminReentry, *response.Error) {
	reentries := []adminReentry{}
	params := struct{}{}
	err := dbClient.Select(&reentries, "sql/reentry/select_pending_reentries.sql", &params)
	if err != nil {
		return api.ResGetAdminReentry{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "承認待ち再入部申請一覧の取得に失敗しました", Log: err.Error()}
	}

	res := api.ResGetAdminReentry{}
	rerr := copier.Copy(&res.Reentries, &reentries)
	if rerr != nil {
		return api.ResGetAdminReentry{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "承認待ち再入部申請一覧の取得に失敗しました", Log: rerr.Error()}
	}
	if res.Reentries == nil {
		res.Reentries = []api.ResGetAdminReentryObjectReentry{}
	}
	return res, nil
}
