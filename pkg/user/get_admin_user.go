package user

import (
	"net/http"

	"github.com/jinzhu/copier"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/admin"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/utils"
	"github.com/labstack/echo/v4"
)

func GetAdminUser(ctx echo.Context, dbClient db.Client, params api.GetAdminUserParams) (api.ResGetAdminUser, *response.Error) {
	userId := ctx.Get("user_id").(string)

	isAdmin, authErr := admin.CheckUserIsAdmin(dbClient, userId)
	if authErr != nil {
		return api.ResGetAdminUser{}, authErr
	}
	if !isAdmin {
		return api.ResGetAdminUser{}, &response.Error{
			Code:    http.StatusForbidden,
			Level:   "Info",
			Message: "管理者ユーザー一覧を取得する権限がありません",
			Log:     "user is not admin",
		}
	}

	res := api.ResGetAdminUser{}

	limit := utils.ClampInt(params.Limit, 1, 500, 100)
	offset := utils.ClampInt(params.Offset, 0, int(^uint(0)>>1), 0)

	adminUsers, total, err := getAdminUserList(dbClient, offset, limit, params.Query, params.SchoolGrade, params.IsAdmin)
	if err != nil {
		return api.ResGetAdminUser{}, err
	}

	copyErr := copier.Copy(&res.Users, &adminUsers)
	if copyErr != nil {
		return api.ResGetAdminUser{}, &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Error",
			Message: "不明なエラーが発生しました",
			Log:     copyErr.Error(),
		}
	}
	if res.Users == nil {
		res.Users = []api.ResGetAdminUserObjectUser{}
	}
	res.Total = total

	return res, nil
}

type adminUser struct {
	UserId         string                  `db:"user_id"`
	StudentNumber  string                  `db:"student_number"`
	IsAdmin        bool                    `db:"is_admin"`
	Profile        adminUserProfile        `db:"-"`
	PrivateProfile adminUserPrivateProfile `db:"-"`
}

type adminUserProfile struct {
	Username          string `db:"username"`
	IconUrl           string `db:"icon_url"`
	SchoolGrade       int    `db:"school_grade"`
	DiscordUserId     string `db:"discord_userid"`
	ShortIntroduction string `db:"short_introduction"`
	Introduction      string `db:"introduction"`
	ActiveLimit       string `db:"active_limit"`
}

type adminUserPrivateProfile struct {
	FirstName             string `db:"first_name"`
	LastName              string `db:"last_name"`
	FirstNameKana         string `db:"first_name_kana"`
	LastNameKana          string `db:"last_name_kana"`
	IsMale                bool   `db:"is_male"`
	PhoneNumber           string `db:"phone_number"`
	Address               string `db:"address"`
	ParentName            string `db:"parent_name"`
	ParentLastName        string `db:"parent_last_name"`
	ParentFirstName       string `db:"parent_first_name"`
	ParentCellphoneNumber string `db:"parent_cellphone_number"`
	ParentHomephoneNumber string `db:"parent_homephone_number"`
	ParentAddress         string `db:"parent_address"`
}

type adminUserRow struct {
	UserId            string `db:"user_id"`
	StudentNumber     string `db:"student_number"`
	Username          string `db:"username"`
	SchoolGrade       int    `db:"school_grade"`
	IconUrl           string `db:"icon_url"`
	DiscordUserId     string `db:"discord_userid"`
	ActiveLimit       string `db:"active_limit"`
	ShortIntroduction string `db:"short_introduction"`
	IsAdmin           bool   `db:"is_admin"`
	Introduction      string `db:"introduction"`
	FirstName         string `db:"first_name"`
	LastName          string `db:"last_name"`
	FirstNameKana     string `db:"first_name_kana"`
	LastNameKana      string `db:"last_name_kana"`
	IsMale            bool   `db:"is_male"`
	PhoneNumber       string `db:"phone_number"`
	Address           string `db:"address"`
	ParentName        string `db:"parent_name"`
	ParentLastName    string `db:"parent_last_name"`
	ParentFirstName   string `db:"parent_first_name"`
	ParentCellphone   string `db:"parent_cellphone_number"`
	ParentHomephone   string `db:"parent_homephone_number"`
	ParentAddress     string `db:"parent_address"`
}

func (r adminUserRow) toAdminUser() adminUser {
	return adminUser{
		UserId:        r.UserId,
		StudentNumber: r.StudentNumber,
		IsAdmin:       r.IsAdmin,
		Profile: adminUserProfile{
			Username:          r.Username,
			IconUrl:           r.IconUrl,
			SchoolGrade:       r.SchoolGrade,
			DiscordUserId:     r.DiscordUserId,
			ShortIntroduction: r.ShortIntroduction,
			Introduction:      r.Introduction,
			ActiveLimit:       r.ActiveLimit,
		},
		PrivateProfile: adminUserPrivateProfile{
			FirstName:             r.FirstName,
			LastName:              r.LastName,
			FirstNameKana:         r.FirstNameKana,
			LastNameKana:          r.LastNameKana,
			IsMale:                r.IsMale,
			PhoneNumber:           r.PhoneNumber,
			Address:               r.Address,
			ParentName:            r.ParentName,
			ParentLastName:        r.ParentLastName,
			ParentFirstName:       r.ParentFirstName,
			ParentCellphoneNumber: r.ParentCellphone,
			ParentHomephoneNumber: r.ParentHomephone,
			ParentAddress:         r.ParentAddress,
		},
	}
}

func getAdminUserList(dbClient db.Client, offset, limit int, query *string, schoolGrade *int, isAdmin *bool) ([]adminUser, int, *response.Error) {
	params := struct {
		Offset      *int    `twowaysql:"offset"`
		Limit       *int    `twowaysql:"limit"`
		Query       *string `twowaysql:"query"`
		SchoolGrade *int    `twowaysql:"schoolGrade"`
		IsAdmin     *bool   `twowaysql:"isAdmin"`
	}{
		Offset:      &offset,
		Limit:       &limit,
		Query:       query,
		SchoolGrade: schoolGrade,
		IsAdmin:     isAdmin,
	}

	total, countErr := countAdminUserList(dbClient, &params)
	if countErr != nil {
		return []adminUser{}, 0, countErr
	}

	rows := []adminUserRow{}
	if err := dbClient.Select(&rows, "sql/user/select_admin_user_list.sql", &params); err != nil {
		return []adminUser{}, 0, &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Error",
			Message: "不明なエラーが発生しました",
			Log:     err.Error(),
		}
	}

	adminUsers := make([]adminUser, 0, len(rows))
	for _, row := range rows {
		adminUsers = append(adminUsers, row.toAdminUser())
	}

	return adminUsers, total, nil
}

func countAdminUserList(dbClient db.Client, params interface{}) (int, *response.Error) {
	result := []struct {
		Total int `db:"total"`
	}{}
	if err := dbClient.Select(&result, "sql/user/count_admin_user_list.sql", params); err != nil {
		return 0, &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Error",
			Message: "不明なエラーが発生しました",
			Log:     err.Error(),
		}
	}
	if len(result) == 0 {
		return 0, nil
	}
	return result[0].Total, nil
}
