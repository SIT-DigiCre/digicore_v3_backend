package user

import (
	"net/http"
	"strings"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func PutUserMePrivate(ctx echo.Context, dbClient db.TransactionClient, requestBody api.ReqPutUserMePrivate) (api.ResGetUserMePrivate, *response.Error) {
	userId := ctx.Get("user_id").(string)
	if err := validateParentNameFields(requestBody); err != nil {
		return api.ResGetUserMePrivate{}, err
	}
	err := updateUserPrivate(dbClient, userId, requestBody)
	if err != nil {
		return api.ResGetUserMePrivate{}, err
	}

	return GetUserMePrivate(ctx, dbClient)
}

// 姓と名は必ずセットで送信する必要がある。未送信（nil）は対象外
func validateParentNameFields(req api.ReqPutUserMePrivate) *response.Error {
	if req.ParentLastName != nil && strings.TrimSpace(*req.ParentLastName) == "" {
		return &response.Error{Code: http.StatusBadRequest, Level: "Info", Message: "緊急連絡先の名字に空文字は指定できません", Log: "parentLastName is empty"}
	}
	if req.ParentFirstName != nil && strings.TrimSpace(*req.ParentFirstName) == "" {
		return &response.Error{Code: http.StatusBadRequest, Level: "Info", Message: "緊急連絡先の名前に空文字は指定できません", Log: "parentFirstName is empty"}
	}
	if (req.ParentLastName != nil) != (req.ParentFirstName != nil) {
		return &response.Error{Code: http.StatusBadRequest, Level: "Info", Message: "緊急連絡先の名字と名前は両方指定してください", Log: "parentLastName and parentFirstName must be provided together"}
	}
	return nil
}

func updateUserPrivate(dbClient db.TransactionClient, userId string, requestBody api.ReqPutUserMePrivate) *response.Error {
	parentLastName, parentFirstName, resolveErr := resolveParentNameFields(dbClient, userId, requestBody)
	if resolveErr != nil {
		return resolveErr
	}

	params := struct {
		UserId                string  `twowaysql:"userId"`
		FirstName             string  `twowaysql:"firstName"`
		LastName              string  `twowaysql:"lastName"`
		FirstNameKana         string  `twowaysql:"firstNameKana"`
		LastNameKana          string  `twowaysql:"lastNameKana"`
		IsMale                bool    `twowaysql:"isMale"`
		PhoneNumber           string  `twowaysql:"phoneNumber"`
		Address               string  `twowaysql:"address"`
		ParentLastName        string  `twowaysql:"parentLastName"`
		ParentFirstName       string  `twowaysql:"parentFirstName"`
		ParentCellphoneNumber string  `twowaysql:"parentCellphoneNumber"`
		ParentHomephoneNumber *string `twowaysql:"parentHomephoneNumber"`
		ParentAddress         string  `twowaysql:"parentAddress"`
	}{
		UserId:                userId,
		FirstName:             requestBody.FirstName,
		LastName:              requestBody.LastName,
		FirstNameKana:         requestBody.FirstNameKana,
		LastNameKana:          requestBody.LastNameKana,
		IsMale:                requestBody.IsMale,
		PhoneNumber:           requestBody.PhoneNumber,
		Address:               requestBody.Address,
		ParentLastName:        parentLastName,
		ParentFirstName:       parentFirstName,
		ParentCellphoneNumber: requestBody.ParentCellphoneNumber,
		ParentHomephoneNumber: requestBody.ParentHomephoneNumber,
		ParentAddress:         requestBody.ParentAddress,
	}
	_, _, err := dbClient.DuplicateUpdate("sql/user/insert_user_private.sql", "sql/user/update_user_private.sql", &params)
	if err != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	return nil
}

// resolveParentNameFields はリクエストで送信された値があればそれを使い、未送信なら既存値を返す。既存レコードが無い場合は未送信を空文字とする
func resolveParentNameFields(dbClient db.Client, userId string, req api.ReqPutUserMePrivate) (parentLastName, parentFirstName string, resErr *response.Error) {
	// 2フィールドすべてが送信済みの場合は既存取得を省略する
	if req.ParentLastName != nil && req.ParentFirstName != nil {
		return *req.ParentLastName, *req.ParentFirstName, nil
	}

	existing, err := getUserPrivateFromUserId(dbClient, userId)
	if err != nil {
		if err.Code != http.StatusNotFound {
			return "", "", err
		}
		// 新規登録の場合は未送信を空文字で返す
		if req.ParentLastName != nil {
			parentLastName = *req.ParentLastName
		}
		if req.ParentFirstName != nil {
			parentFirstName = *req.ParentFirstName
		}
		return parentLastName, parentFirstName, nil
	}
	parentLastName = existing.ParentLastName
	parentFirstName = existing.ParentFirstName
	if req.ParentLastName != nil {
		parentLastName = *req.ParentLastName
	}
	if req.ParentFirstName != nil {
		parentFirstName = *req.ParentFirstName
	}
	return parentLastName, parentFirstName, nil
}
