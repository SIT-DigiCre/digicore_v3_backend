package group

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
)

// CheckUserIsAdmin は、ユーザーが管理者かどうかを確認する（エクスポート関数）
func CheckUserIsAdmin(dbClient db.Client, userId string) (bool, *response.Error) {
	return checkUserIsAdmin(dbClient, userId)
}

// checkUserIsAdmin は、ユーザーが管理者かどうかを確認する
func checkUserIsAdmin(dbClient db.Client, userId string) (bool, *response.Error) {
	params := struct {
		UserId string `twowaysql:"userId"`
	}{
		UserId: userId,
	}

	result := []struct {
		IsAdmin bool `db:"is_admin"`
	}{}

	err := dbClient.Select(&result, "sql/group/select_user_is_admin.sql", &params)
	if err != nil {
		return false, &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Info",
			Message: "DBエラーが発生しました",
			Log:     err.Error(),
		}
	}

	if len(result) == 0 {
		return false, nil
	}

	return result[0].IsAdmin, nil
}

// insertGroup は、グループをデータベースに挿入する
func insertGroup(dbClient db.TransactionClient, name, description string, joinable bool) *response.Error {
	params := struct {
		Name        string `twowaysql:"name"`
		Description string `twowaysql:"description"`
		Joinable    bool   `twowaysql:"joinable"`
	}{
		Name:        name,
		Description: description,
		Joinable:    joinable,
	}

	_, err := dbClient.Exec("sql/group/insert_group.sql", &params, true)
	if err != nil {
		return &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Info",
			Message: "DBエラーが発生しました",
			Log:     err.Error(),
		}
	}

	return nil
}

// insertGroupUser は、ユーザーをグループに追加する
func insertGroupUser(dbClient db.TransactionClient, userId, groupId string) *response.Error {
	params := struct {
		UserId  string `twowaysql:"userId"`
		GroupId string `twowaysql:"groupId"`
	}{
		UserId:  userId,
		GroupId: groupId,
	}

	_, err := dbClient.Exec("sql/group/insert_groups_users.sql", &params, false)
	if err != nil {
		return &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Info",
			Message: "DBエラーが発生しました",
			Log:     err.Error(),
		}
	}

	return nil
}

// insertGroupClaim は、グループにclaimを追加する
func insertGroupClaim(dbClient db.TransactionClient, groupId, claim string) *response.Error {
	params := struct {
		GroupId string `twowaysql:"groupId"`
		Claim   string `twowaysql:"claim"`
	}{
		GroupId: groupId,
		Claim:   claim,
	}

	_, err := dbClient.Exec("sql/group/insert_group_claims.sql", &params, false)
	if err != nil {
		return &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Info",
			Message: "DBエラーが発生しました",
			Log:     err.Error(),
		}
	}

	return nil
}

// checkGroupExists は、グループが存在するかを確認する
func checkGroupExists(dbClient db.TransactionClient, groupId string) (bool, *response.Error) {
	params := struct {
		GroupId string `twowaysql:"groupId"`
	}{
		GroupId: groupId,
	}

	result := []struct {
		GroupExists bool `db:"group_exists"`
	}{}

	err := dbClient.Select(&result, "sql/group/select_group_exists.sql", &params)
	if err != nil {
		return false, &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Info",
			Message: "DBエラーが発生しました",
			Log:     err.Error(),
		}
	}

	if len(result) == 0 {
		return false, nil
	}

	return result[0].GroupExists, nil
}

// checkGroupIsAdminGroup は、グループがadminグループかを確認する
func checkGroupIsAdminGroup(dbClient db.TransactionClient, groupId string) (bool, *response.Error) {
	params := struct {
		GroupId string `twowaysql:"groupId"`
	}{
		GroupId: groupId,
	}

	result := []struct {
		IsAdminGroup bool `db:"is_admin_group"`
	}{}

	err := dbClient.Select(&result, "sql/group/select_group_is_admin_group.sql", &params)
	if err != nil {
		return false, &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Info",
			Message: "DBエラーが発生しました",
			Log:     err.Error(),
		}
	}

	if len(result) == 0 {
		return false, nil
	}

	return result[0].IsAdminGroup, nil
}

// checkGroupJoinable は、グループのjoinableを確認する
func checkGroupJoinable(dbClient db.TransactionClient, groupId string) (bool, *response.Error) {
	params := struct {
		GroupId string `twowaysql:"groupId"`
	}{
		GroupId: groupId,
	}

	result := []struct {
		Joinable bool `db:"joinable"`
	}{}

	err := dbClient.Select(&result, "sql/group/select_group_joinable.sql", &params)
	if err != nil {
		return false, &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Info",
			Message: "DBエラーが発生しました",
			Log:     err.Error(),
		}
	}

	if len(result) == 0 {
		return false, &response.Error{
			Code:    http.StatusNotFound,
			Level:   "Info",
			Message: "指定されたグループが存在しません",
			Log:     "group not found",
		}
	}

	return result[0].Joinable, nil
}

// checkUserIsGroupMember は、ユーザーがグループに所属しているかを確認する
func checkUserIsGroupMember(dbClient db.TransactionClient, userId, groupId string) (bool, *response.Error) {
	params := struct {
		UserId  string `twowaysql:"userId"`
		GroupId string `twowaysql:"groupId"`
	}{
		UserId:  userId,
		GroupId: groupId,
	}

	result := []struct {
		IsMember bool `db:"is_member"`
	}{}

	err := dbClient.Select(&result, "sql/group/select_is_group_member.sql", &params)
	if err != nil {
		return false, &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Info",
			Message: "DBエラーが発生しました",
			Log:     err.Error(),
		}
	}

	if len(result) == 0 {
		return false, nil
	}

	return result[0].IsMember, nil
}

// checkUserExists は、ユーザーが存在するかを確認する
func checkUserExists(dbClient db.TransactionClient, userId string) (bool, *response.Error) {
	params := struct {
		UserId string `twowaysql:"userId"`
	}{
		UserId: userId,
	}

	result := []struct {
		UserExists bool `db:"user_exists"`
	}{}

	err := dbClient.Select(&result, "sql/group/select_user_exists.sql", &params)
	if err != nil {
		return false, &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Info",
			Message: "DBエラーが発生しました",
			Log:     err.Error(),
		}
	}

	if len(result) == 0 {
		return false, nil
	}

	return result[0].UserExists, nil
}

// incrementGroupUserCount は、グループのuser_countをインクリメントする
func incrementGroupUserCount(dbClient db.TransactionClient, groupId string) *response.Error {
	params := struct {
		GroupId string `twowaysql:"groupId"`
	}{
		GroupId: groupId,
	}

	_, err := dbClient.Exec("sql/group/update_group_user_count_increment.sql", &params, false)
	if err != nil {
		return &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Info",
			Message: "DBエラーが発生しました",
			Log:     err.Error(),
		}
	}

	return nil
}
