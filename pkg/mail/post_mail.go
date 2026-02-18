package mail

import (
	"fmt"
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/env"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/user"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func PostMail(ctx echo.Context, dbClient db.Client, requestBody api.ReqPostMail) (api.ResPostMail, *response.Error) {
	// addressesとuserIdsの両方が空の場合はエラー
	addressCount := 0
	if requestBody.Addresses != nil {
		addressCount = len(*requestBody.Addresses)
	}
	userIdCount := 0
	if requestBody.UserIds != nil {
		userIdCount = len(*requestBody.UserIds)
	}
	if addressCount == 0 && userIdCount == 0 {
		return api.ResPostMail{}, &response.Error{
			Code:    http.StatusBadRequest,
			Level:   "Info",
			Message: "送信先アドレスまたは送信先ユーザーIDのいずれかを指定してください",
			Log:     "both addresses and userIds are empty",
		}
	}

	addresses := make([]string, 0, addressCount+userIdCount+1)
	failures := []struct {
		Address string `json:"address"`
		Error   string `json:"error"`
	}{}

	// メールアドレスを追加
	if requestBody.Addresses != nil {
		for _, addr := range *requestBody.Addresses {
			addresses = append(addresses, string(addr))
		}
	}

	// ユーザーIDから学籍番号を一括取得してメールアドレスを生成
	if requestBody.UserIds != nil && len(*requestBody.UserIds) > 0 {
		userIds := make([]string, len(*requestBody.UserIds))
		for i, uid := range *requestBody.UserIds {
			userIds[i] = uid.String()
		}
		studentNumbers, respErr := user.GetStudentNumbersFromUserIds(dbClient, userIds)
		if respErr != nil {
			return api.ResPostMail{}, respErr
		}
		for _, uid := range userIds {
			sn, ok := studentNumbers[uid]
			if !ok || sn == "" {
				logrus.Warnf("ユーザーID %s の学籍番号が見つかりません", uid)
				failures = append(failures, struct {
					Address string `json:"address"`
					Error   string `json:"error"`
				}{
					Address: fmt.Sprintf("user_id:%s", uid),
					Error:   "学籍番号が見つかりません",
				})
				continue
			}
			addresses = append(addresses, fmt.Sprintf("%s@shibaura-it.ac.jp", sn))
		}
	}

	if requestBody.SendToAdmin != nil && *requestBody.SendToAdmin {
		if env.AdminEmail != "" {
			addresses = append(addresses, env.AdminEmail)
		}
	}

	successCount := 0

	for _, address := range addresses {
		err := sendEmail(address, requestBody.Subject, requestBody.Body)
		if err != nil {
			failures = append(failures, struct {
				Address string `json:"address"`
				Error   string `json:"error"`
			}{
				Address: address,
				Error:   err.Error(),
			})
			logrus.Errorf("メール送信失敗 [%s]: %v", address, err)
		} else {
			successCount++
		}
	}

	res := api.ResPostMail{
		SuccessCount: successCount,
		Failures:     failures,
	}

	return res, nil
}
