package user

import (
	"net/http"

	"github.com/jinzhu/copier"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func GetAdminUser(ctx echo.Context, dbClient db.Client, params api.GetAdminUserParams) (api.ResGetAdminUser, *response.Error) {
	res := api.ResGetAdminUser{}

	limit := params.Limit
	if limit == nil {
		defaultLimit := 100
		limit = &defaultLimit
	}
	if *limit > 500 {
		maxLimit := 500
		limit = &maxLimit
	}
	if *limit <= 0 {
		minLimit := 1
		limit = &minLimit
	}

	offset := params.Offset
	if offset == nil {
		defaultOffset := 0
		offset = &defaultOffset
	}
	if *offset < 0 {
		minOffset := 0
		offset = &minOffset
	}

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
	Total             int    `db:"total"`
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

func getAdminUserList(dbClient db.Client, offset, limit *int, query *string, schoolGrade *int, isAdmin *bool) ([]adminUser, int, *response.Error) {
	params := struct {
		Offset      *int    `twowaysql:"offset"`
		Limit       *int    `twowaysql:"limit"`
		Query       *string `twowaysql:"query"`
		SchoolGrade *int    `twowaysql:"schoolGrade"`
		IsAdmin     *bool   `twowaysql:"isAdmin"`
	}{
		Offset:      offset,
		Limit:       limit,
		Query:       query,
		SchoolGrade: schoolGrade,
		IsAdmin:     isAdmin,
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

	total := 0
	adminUsers := make([]adminUser, 0, len(rows))
	for _, row := range rows {
		adminUsers = append(adminUsers, row.toAdminUser())
		total = row.Total
	}

	return adminUsers, total, nil
}
