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
	if requestBody.Command == "/remind" {
		return mattermostCmdRemind(dbClient, decode, requestBody)
	} else if requestBody.Command == "/remind_list" {
		return mattermostCmdRemindList(dbClient, requestBody)

	} else if requestBody.Command == "/remind_delete" {
		return mattermostCmdRemindDelete(dbClient, decode, requestBody)
	}
	return api.ResPostMattermostCmd{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "コマンドが存在しません", Log: "No foud command"}
}

func mattermostCmdRemind(dbClient db.TransactionClient, decode string, requestBody api.ReqPostMattermostCmd) (api.ResPostMattermostCmd, *response.Error) {
	lines := strings.Split(decode, "\n")
	infos := strings.Split(lines[0], " ")
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

func mattermostCmdRemindList(dbClient db.Client, requestBody api.ReqPostMattermostCmd) (api.ResPostMattermostCmd, *response.Error) {
	reminds, err := selectMattermostRemind(dbClient, requestBody.UserName, requestBody.ChannelName)
	if err != nil {
		return api.ResPostMattermostCmd{}, err
	}
	text := ""
	for _, remind := range reminds {
		text += fmt.Sprintf("- %s(%s)\n%s\n", remind.RemindDate.Format("2006-01-02T15:04"), remind.Id, remind.Body)
	}
	if text == "" {
		text = "予約済み投稿はありません"
	}
	return api.ResPostMattermostCmd{Text: text, IconEmoji: "alarm_clock", Username: "remind"}, nil
}

func mattermostCmdRemindDelete(dbClient db.TransactionClient, decode string, requestBody api.ReqPostMattermostCmd) (api.ResPostMattermostCmd, *response.Error) {
	lines := strings.Split(decode, "\n")
	infos := strings.Split(lines[0], " ")
	err := deleteMattermostRemind(dbClient, requestBody.UserName, requestBody.ChannelName, infos[0])
	if err != nil {
		return api.ResPostMattermostCmd{}, err
	}
	return api.ResPostMattermostCmd{Text: fmt.Sprintf("%sの投稿を削除しました", infos[0]), IconEmoji: "alarm_clock", Username: "remind"}, nil
}

func insertMattermostRemind(dbClient db.TransactionClient, userName string, channelName string, body string, remindDate time.Time) *response.Error {
	params := struct {
		UserName    string    `twowaysql:"userName"`
		ChannelName string    `twowaysql:"channelName"`
		Body        string    `twowaysql:"body"`
		RemindDate  time.Time `twowaysql:"remindDate"`
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

type remind struct {
	Id         string    `db:"id"`
	Body       string    `db:"body"`
	RemindDate time.Time `db:"remind_date"`
}

func selectMattermostRemind(dbClient db.Client, userName string, channelName string) ([]remind, *response.Error) {
	params := struct {
		UserName    string `twowaysql:"userName"`
		ChannelName string `twowaysql:"channelName"`
	}{
		UserName:    userName,
		ChannelName: channelName,
	}
	reminds := []remind{}
	err := dbClient.Select(&reminds, "sql/mattermost/select_mattermost_remind_post.sql", &params)
	if err != nil {
		return []remind{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "予約済み投稿の取得に失敗しました", Log: err.Error()}
	}
	return reminds, nil
}

func deleteMattermostRemind(dbClient db.TransactionClient, userName string, channelName string, id string) *response.Error {
	params := struct {
		Id          string `twowaysql:"id"`
		UserName    string `twowaysql:"userName"`
		ChannelName string `twowaysql:"channelName"`
	}{
		Id:          id,
		UserName:    userName,
		ChannelName: channelName,
	}
	_, err := dbClient.Exec("sql/mattermost/delete_mattermost_remind_post.sql", &params, false)
	if err != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "予約済み投稿の削除に失敗しました", Log: err.Error()}
	}
	return nil
}
