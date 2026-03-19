package admin

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
)

// 定義済みの役職claimリスト
var adminClaims = []string{"account", "infra"}

// GetAdminClaims は役職claimリストのコピーを返す（外部からのミューテートを防ぐ）
func GetAdminClaims() []string {
	result := make([]string, len(adminClaims))
	copy(result, adminClaims)
	return result
}

type SelectClient interface {
	Select(dest interface{}, queryPath string, params interface{}) error
}

// いずれかの役職claimを持っているか確認する（幹部判定）
func CheckUserIsAdmin(dbClient SelectClient, userId string) (bool, *response.Error) {
	params := struct {
		UserId      string   `twowaysql:"userId"`
		AdminClaims []string `twowaysql:"adminClaims"`
	}{
		UserId:      userId,
		AdminClaims: adminClaims,
	}

	result := []struct {
		IsAdmin bool `db:"is_admin"`
	}{}

	err := dbClient.Select(&result, "sql/admin/select_user_is_admin.sql", &params)
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

// 特定のclaimを持っているか確認する
func CheckUserHasClaim(dbClient SelectClient, userId string, claim string) (bool, *response.Error) {
	params := struct {
		UserId string `twowaysql:"userId"`
		Claim  string `twowaysql:"claim"`
	}{
		UserId: userId,
		Claim:  claim,
	}

	result := []struct {
		HasClaim bool `db:"has_claim"`
	}{}

	err := dbClient.Select(&result, "sql/admin/select_user_has_claim.sql", &params)
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

	return result[0].HasClaim, nil
}
