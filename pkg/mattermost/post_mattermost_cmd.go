package mattermost

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func PostMattermostCmd(ctx echo.Context, dbClient db.TransactionClient, requestBody api.ReqPostMattermostCmd) (api.ResPostMattermostCmd, *response.Error) {
	decode, rerr := url.QueryUnescape(requestBody.Text)
	if rerr != nil {
		return api.ResPostMattermostCmd{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "デコードに失敗しました", Log: rerr.Error()}
	}
	lines := strings.Split(decode, "\n")
	infos := strings.Split(lines[0], " ")
	if requestBody.Command == "/remind" {
		jst, _ := time.LoadLocation("Asia/Tokyo")
		remindDate, rerr := time.ParseInLocation("2006-01-02T15:04", infos[0], jst)
		channelName := requestBody.ChannelName
		if 3 <= len(infos) {
			channelName = infos[2][1:]
		}
		if rerr != nil {
			return api.ResPostMattermostCmd{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "時刻のパースに失敗しました", Log: rerr.Error()}
		}
		err := insertMattermostRemind(dbClient, requestBody.UserName, channelName, strings.Join(lines[1:], "\n"), remindDate)
		if err != nil {
			return api.ResPostMattermostCmd{}, err
		}
		return api.ResPostMattermostCmd{Text: fmt.Sprintf("%sに%sへの投稿を受け付けました", remindDate.Format(time.DateTime), channelName), IconEmoji: "alarm_clock", Username: "remind"}, nil
	}
	return api.ResPostMattermostCmd{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "コマンドが存在しません", Log: "No foud command"}
}

func insertMattermostRemind(dbClient db.TransactionClient, userName string, channelName string, body string, remindDate time.Time) *response.Error {
	params := struct {
		UserName    string    `twowaysql:"user_name"`
		ChannelName string    `twowaysql:"channel_name"`
		Body        string    `twowaysql:"body"`
		RemindDate  time.Time `twowaysql:"remind_date"`
	}{
		UserName:    userName,
		ChannelName: channelName,
		Body:        body,
		RemindDate:  remindDate,
	}
	_, err := dbClient.Exec("sql/mattermost/insert_mattermost_remind_post.sql", &params, false)
	if err != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	return nil
}
