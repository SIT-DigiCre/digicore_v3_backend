package activity

import (
	"net/http"
	"time"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

func GetActivityPlacePlaceHistory(ctx echo.Context, dbClient db.Client, place string, startAt time.Time, endAt time.Time) (api.ResGetActivityPlacePlaceHistory, *response.Error) {
	res := api.ResGetActivityPlacePlaceHistory{}

	if startAt.After(endAt) {
		errResp := &response.Error{
			Code:    http.StatusBadRequest,
			Level:   "Info",
			Message: "開始日時は終了日時以前である必要があります",
			Log:     "startAt is after endAt",
		}
		return api.ResGetActivityPlacePlaceHistory{}, errResp
	}

	users, errResp := selectPlaceHistory(dbClient, place, startAt, endAt)
	if errResp != nil {
		return api.ResGetActivityPlacePlaceHistory{}, errResp
	}
	rerr := copier.Copy(&res.Users, &users)
	if rerr != nil {
		return api.ResGetActivityPlacePlaceHistory{}, &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Error",
			Message: "訪問履歴のデータベース取得に失敗しました",
			Log:     rerr.Error(),
		}
	}
	if res.Users == nil {
		res.Users = []api.ResGetActivityPlacePlaceHistoryObjectUser{}
	}
	return res, nil
}

type placeHistoryUser struct {
	UserId            string `db:"user_id"`
	Username          string `db:"username"`
	ShortIntroduction string `db:"short_introduction"`
	IconUrl           string `db:"icon_url"`
	CheckInCount      int    `db:"check_in_count"`
}

func selectPlaceHistory(dbClient db.Client, place string, startDate time.Time, endDate time.Time) ([]placeHistoryUser, *response.Error) {
	sameDay := startDate.Year() == endDate.Year() &&
		startDate.Month() == endDate.Month() &&
		startDate.Day() == endDate.Day()
	params := struct {
		SameDay   bool      `twowaysql:"sameDay"`
		Place     string    `twowaysql:"place"`
		StartDate time.Time `twowaysql:"startDate"`
		EndDate   time.Time `twowaysql:"endDate"`
	}{
		SameDay:   sameDay,
		Place:     place,
		StartDate: startDate,
		EndDate:   endDate,
	}
	users := []placeHistoryUser{}
	err := dbClient.Select(&users, "sql/activity/select_place_history.sql", &params)
	if err != nil {
		return []placeHistoryUser{}, &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Error",
			Message: "訪問履歴の取得に失敗しました",
			Log:     err.Error(),
		}
	}
	return users, nil
}
