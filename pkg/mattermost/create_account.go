package mattermost

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/env"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/user"
	"github.com/labstack/echo/v4"
	"github.com/mattermost/mattermost-server/v5/model"
)

var inviteID string
var inviteIDGeneratedAt time.Time

type ResponseCreatedUser struct {
	Username string `json:"username"`
}
type ResponseError struct {
	Message string `json:"message"`
}

type RequestCreateUserInfo struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}
func (cu *RequestCreateUserInfo) validate() error {
	if len(cu.Username) == 0 {
		return errors.New("ユーザ名が空です")
	}
	if len(cu.Password) < 8 {
		return errors.New("パスワードが短すぎます")
	}

	return nil
}

// Create Mattermost user and add to team
// @Router /mattermost/create_user
// @Security Authorization
// @Success 200 {object}  ResponseCreatedUser
// @Failure 500 {object} ResponseError
func (c Context) CreateUser(e echo.Context) error {
	userId, err := user.GetUserId(&e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseError{ Message: err.Error() })
	}
	createUserInfo := RequestCreateUserInfo{}
	if err := e.Bind(&createUserInfo); err != nil {
		return e.JSON(http.StatusBadRequest, ResponseError{ Message: "データの読み込みに失敗しました" })
	}
	if err := createUserInfo.validate(); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}
	studentNumber, err := user.GetStudentNumber(c.DB, userId)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseError{ Message: err.Error() })
	}
	email := fmt.Sprintf("%s@shibaura-it.ac.jp", studentNumber)

	client := model.NewAPIv4Client(env.MattermostURL)
	client.Login(env.MattermostAdminAccount, env.MattermostAdminPassword)
	team, _ := client.RegenerateTeamInviteId(env.MattermostDigicreTeamID)
	if team == nil {
		return e.JSON(http.StatusInternalServerError, ResponseError{ Message: "招待IDの生成に失敗しました" })
	}
	inviteID = team.InviteId

	user := &model.User{
		Username: createUserInfo.Username,
		Email: email,
		Password: createUserInfo.Password,
		Nickname: createUserInfo.Nickname,
	}
	createdUser, _ := client.CreateUserWithInviteId(user, inviteID)
	if createdUser == nil {
		return e.JSON(http.StatusInternalServerError, ResponseError{ Message: "ユーザの追加に失敗しました。同じユーザIDが既に使われている可能性があります。" })
	}
	// ユーザ作成には成功しているので、以降の処理は失敗してもOKを返す
	responseOK :=  e.JSON(http.StatusOK, ResponseCreatedUser{ Username: createdUser.Username })

	var iconURL string
	err = c.DB.QueryRow("SELECT icon_url FROM user_profiles WHERE user_id = UUID_TO_BIN(?)", userId).
		Scan(&iconURL)
	if err == sql.ErrNoRows {
		return responseOK
	} else if err != nil {
		return responseOK
	}
	if (0 < len(iconURL)) {
		res, err := http.Get(iconURL)
		if err != nil {
			return responseOK
		}
		defer res.Body.Close()
		iconData, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return responseOK
		}
		client.SetProfileImage(createdUser.Id, iconData)
	}
	return responseOK
}
