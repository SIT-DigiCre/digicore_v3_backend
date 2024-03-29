package budget

import (
	"database/sql"
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/utils"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/maps"
)

func GetBudgetBudgetId(ctx echo.Context, dbClient db.Client, budgetId string) (api.ResGetBudgetBudgetId, *response.Error) {
	res := api.ResGetBudgetBudgetId{}
	budget, err := getBudgetFromBudgetId(dbClient, budgetId)
	if err != nil {
		return api.ResGetBudgetBudgetId{}, err
	}
	rerr := copier.Copy(&res, &budget)
	if rerr != nil {
		return api.ResGetBudgetBudgetId{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "稟議の取得に失敗しました", Log: rerr.Error()}
	}
	return res, nil
}

type budgetDetail struct {
	BudgetId      string `db:"budget_id"`
	Name          string `db:"name"`
	Class         string `db:"class"`
	Status        string `db:"status"`
	Purpose       string `db:"purpose"`
	Budget        int    `db:"budget"`
	Settlement    int    `db:"settlement"`
	MattermostUrl string `db:"mattermost_url"`
	Remark        string `db:"remark"`

	Proposer UserInfo

	ProposerUserId   string `db:"proposer_user_id"`
	ProposerIconUrl  string `db:"proposer_icon_url"`
	ProposerUsername string `db:"proposer_username"`

	Approver UserInfo

	ApproverUserId   sql.NullString `db:"approver_user_id"`
	ApproverIconUrl  sql.NullString `db:"approver_icon_url"`
	ApproverUsername sql.NullString `db:"approver_username"`

	Files []utils.FileInfo

	ApprovedAt sql.NullString `db:"approved_at"`
	CreatedAt  string         `db:"created_at"`
	UpdatedAt  string         `db:"updated_at"`
}

type UserInfo struct {
	IconUrl  string `json:"iconUrl"`
	UserId   string `json:"userId"`
	Username string `json:"username"`
}

func getBudgetFromBudgetId(dbClient db.Client, budgetId string) (budgetDetail, *response.Error) {
	params := struct {
		BudgetId string `twowaysql:"budgetId"`
	}{
		BudgetId: budgetId,
	}
	budgetDetails := []budgetDetail{}
	rerr := dbClient.Select(&budgetDetails, "sql/budget/select_budget_from_budget_id.sql", &params)
	if rerr != nil {
		return budgetDetail{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "稟議の取得に失敗しました", Log: rerr.Error()}
	}
	if len(budgetDetails) == 0 {
		return budgetDetail{}, &response.Error{Code: http.StatusNotFound, Level: "Info", Message: "稟議がありません。", Log: "no rows in result"}
	}
	budgetDetails[0].Proposer = UserInfo{UserId: budgetDetails[0].ProposerUserId, IconUrl: budgetDetails[0].ProposerIconUrl, Username: budgetDetails[0].ProposerUsername}
	budgetDetails[0].Approver = UserInfo{UserId: budgetDetails[0].ApproverUserId.String, IconUrl: budgetDetails[0].ApproverIconUrl.String, Username: budgetDetails[0].ApproverUsername.String}
	files, err := getBudgetFileInfo(dbClient, budgetId)
	if err != nil {
		return budgetDetail{}, err
	}
	budgetDetails[0].Files = files
	return budgetDetails[0], nil
}

func getBudgetFileInfo(dbClient db.Client, budgetId string) ([]utils.FileInfo, *response.Error) {
	params := struct {
		BudgetId string `twowaysql:"budgetId"`
	}{
		BudgetId: budgetId,
	}
	rowFileIds := []utils.FileId{}
	rerr := dbClient.Select(&rowFileIds, "sql/budget/select_budget_file_from_budget_id.sql", &params)
	if rerr != nil {
		return []utils.FileInfo{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "稟議の取得に失敗しました", Log: rerr.Error()}
	}
	fileIds := []string{}
	for _, v := range rowFileIds {
		fileIds = append(fileIds, v.FileId)
	}
	res, err := utils.GetFileInfo(dbClient, fileIds)
	if err != nil {
		return []utils.FileInfo{}, err
	}
	return maps.Values(res), nil
}
