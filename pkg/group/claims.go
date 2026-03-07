package group

import (
	"net/http"
	"sort"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
)

type claim struct {
	Claim string `db:"claim"`
}

func GetClaimsFromUserId(dbClient db.Client, userId string) ([]string, *response.Error) {
	params := struct {
		UserId string `twowaysql:"userId"`
	}{
		UserId: userId,
	}
	claims := []claim{}
	err := dbClient.Select(&claims, "sql/group/select_claim_group_from_user_id.sql", &params)
	if err != nil {
		return nil, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "権限一覧の取得に失敗しました", Log: err.Error()}
	}

	claimStrings := make([]string, 0, len(claims))
	for _, groupClaim := range claims {
		claimStrings = append(claimStrings, groupClaim.Claim)
	}
	sort.Strings(claimStrings)

	return claimStrings, nil
}
