package batch

import (
	"fmt"
	"time"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/utils"
	"github.com/sirupsen/logrus"
)

func batch_mattermost_post() {
	ticker := time.NewTicker(time.Minute * 1)
	defer ticker.Stop()
	for range ticker.C {
		dbTranisactionClient, err := db.OpenTransaction()
		if err != nil {
			continue
		}
		logrus.Info("Run mattermost post batch")
		run_mattermost_post(&dbTranisactionClient)
		err = dbTranisactionClient.Commit()
		if err != nil {
			continue
		}
	}
}

type post struct {
	Id          string `db:"id"`
	UserName    string `db:"user_name"`
	ChannelName string `db:"channel_name"`
	Body        string `db:"body"`
}

func run_mattermost_post(dbClient db.TransactionClient) {
	params := struct {
		RemindDate time.Time `twowaysql:"remindDate"`
	}{
		RemindDate: time.Now(),
	}
	posts := []post{}
	err := dbClient.Select(&posts, "sql/batch/select_mattermost_remind_post.sql", &params)
	if err != nil {
		logrus.Error("Failed to get mattermost post remind")
		return
	}
	for _, post := range posts {
		remind_mattermost_post(dbClient, post)
	}
}

func remind_mattermost_post(dbClient db.TransactionClient, post post) {
	logrus.Info(fmt.Sprintf("Post mattermost remind(%s)", post.Id))
	utils.NoticeMattermost(post.Body, post.ChannelName, post.UserName, "bell")
	params := struct {
		Id string `twowaysql:"id"`
	}{
		Id: post.Id,
	}
	_, err := dbClient.Exec("sql/batch/update_mattermost_remind_post.sql", &params, false)
	if err != nil {
		logrus.Error(fmt.Sprintf("Failed to update mattermost remind post(%s): %v", post.Id, err))
	}
}
