package activity

import (
	"net/http"
	"time"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

func GetActivityPlacePlaceHistory(ctx echo.Context, dbClient db.Client, place string, period api.GetActivityPlacePlaceHistoryParamsPeriod, date openapi_types.Date) (api.ResGetActivityPlacePlaceHistory, *response.Error) {
	res := api.ResGetActivityPlacePlaceHistory{}

	// 日付をパース
	parsedDate, err := time.Parse("2006-01-02", date.String())
	if err != nil {
		return api.ResGetActivityPlacePlaceHistory{}, &response.Error{
			Code:    http.StatusBadRequest,
			Level:   "Info",
			Message: "日付の形式が不正です",
			Log:     err.Error(),
		}
	}

	// 日付範囲を計算
	startDate, endDate, errResp := calculateDateRange(string(period), parsedDate)
	if errResp != nil {
		return api.ResGetActivityPlacePlaceHistory{}, errResp
	}

	users, errResp := selectPlaceHistory(dbClient, place, startDate, endDate)
	if errResp != nil {
		return api.ResGetActivityPlacePlaceHistory{}, errResp
	}
	rerr := copier.Copy(&res.Users, &users)
	if rerr != nil {
		return api.ResGetActivityPlacePlaceHistory{}, &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Error",
			Message: "訪問履歴の取得に失敗しました",
			Log:     rerr.Error(),
		}
	}
	if res.Users == nil {
		res.Users = []api.ResGetActivityPlacePlaceHistoryObjectUser{}
	}
	return res, nil
}

func calculateDateRange(period string, date time.Time) (time.Time, time.Time, *response.Error) {
	var startDate, endDate time.Time

	switch period {
	case "day":
		// 指定日の00:00:00 ～ 23:59:59
		startDate = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		endDate = time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, date.Location())
	case "week":
		// 指定日を含む週の1週間前の月曜00:00:00 ～ 日曜23:59:59
		// 指定日の曜日を取得（0=日曜、1=月曜、...、6=土曜）
		weekday := int(date.Weekday())
		if weekday == 0 {
			weekday = 7 // 日曜を7に変換
		}
		// 指定日から1週間前の月曜を計算
		daysToMonday := weekday - 1 + 7 // 1週間前の月曜までの日数
		oneWeekAgoMonday := date.AddDate(0, 0, -daysToMonday)
		startDate = time.Date(oneWeekAgoMonday.Year(), oneWeekAgoMonday.Month(), oneWeekAgoMonday.Day(), 0, 0, 0, 0, date.Location())
		// 1週間前の日曜を計算（開始日の月曜から6日後）
		oneWeekAgoSunday := oneWeekAgoMonday.AddDate(0, 0, 6)
		endDate = time.Date(oneWeekAgoSunday.Year(), oneWeekAgoSunday.Month(), oneWeekAgoSunday.Day(), 23, 59, 59, 999999999, date.Location())
	case "month":
		// 指定日を含む月の1か月前の1日00:00:00 ～ 月末23:59:59
		oneMonthAgo := date.AddDate(0, -1, 0)
		startDate = time.Date(oneMonthAgo.Year(), oneMonthAgo.Month(), 1, 0, 0, 0, 0, date.Location())
		// 1か月前の月末を計算
		oneMonthAgoNextMonth := oneMonthAgo.AddDate(0, 1, 0)
		lastDayOfOneMonthAgo := time.Date(oneMonthAgoNextMonth.Year(), oneMonthAgoNextMonth.Month(), 1, 0, 0, 0, 0, date.Location()).AddDate(0, 0, -1)
		endDate = time.Date(lastDayOfOneMonthAgo.Year(), lastDayOfOneMonthAgo.Month(), lastDayOfOneMonthAgo.Day(), 23, 59, 59, 999999999, date.Location())
	default:
		return time.Time{}, time.Time{}, &response.Error{
			Code:    http.StatusBadRequest,
			Level:   "Info",
			Message: "periodはday、week、monthのいずれかである必要があります",
			Log:     "不正なperiod値: " + period,
		}
	}

	return startDate, endDate, nil
}

type placeHistoryUser struct {
	UserId            string `db:"user_id"`
	Username          string `db:"username"`
	ShortIntroduction string `db:"short_introduction"`
	IconUrl           string `db:"icon_url"`
	CheckInCount      int    `db:"check_in_count"`
}

func selectPlaceHistory(dbClient db.Client, place string, startDate time.Time, endDate time.Time) ([]placeHistoryUser, *response.Error) {
	params := struct {
		Place     string    `twowaysql:"place"`
		StartDate time.Time `twowaysql:"startDate"`
		EndDate   time.Time `twowaysql:"endDate"`
	}{
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
