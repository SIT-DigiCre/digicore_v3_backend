package mail

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/env"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/group"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func PostMail(ctx echo.Context, dbClient db.Client, requestBody api.ReqPostMail) (api.ResPostMail, *response.Error) {
	userIdRaw := ctx.Get("user_id")
	if userIdRaw == nil {
		return api.ResPostMail{}, &response.Error{
			Code:    http.StatusUnauthorized,
			Level:   "Info",
			Message: "認証が必要です",
			Log:     "user_id is not set in context",
		}
	}
	userId, ok := userIdRaw.(string)
	if !ok || userId == "" {
		return api.ResPostMail{}, &response.Error{
			Code:    http.StatusUnauthorized,
			Level:   "Info",
			Message: "認証が必要です",
			Log:     "user_id is invalid",
		}
	}

	isAdmin, err := group.CheckUserIsAdmin(dbClient, userId)
	if err != nil {
		return api.ResPostMail{}, err
	}
	if !isAdmin {
		return api.ResPostMail{}, &response.Error{
			Code:    http.StatusForbidden,
			Level:   "Info",
			Message: "メール送信の権限がありません",
			Log:     "user is not admin",
		}
	}

	addresses := requestBody.Addresses

	if requestBody.SendToAdmin != nil && *requestBody.SendToAdmin {
		if env.AdminEmail != "" {
			addresses = append(addresses, env.AdminEmail)
		}
	}

	successCount := 0
	failures := []struct {
		Address string `json:"address"`
		Error   string `json:"error"`
	}{}

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
