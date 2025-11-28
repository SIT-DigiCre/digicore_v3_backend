package budget

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

func GetBudget(ctx echo.Context, dbClient db.Client, params api.GetBudgetParams) (api.ResGetBudget, *response.Error) {
	res := api.ResGetBudget{}
	budget, err := getBudgetList(dbClient, params.Offset, params.ProposerId)
	if err != nil {
		return api.ResGetBudget{}, err
	}
	rerr := copier.Copy(&res.Budgets, &budget)
	if rerr != nil {
		return api.ResGetBudget{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "稟議一覧の取得に失敗しました", Log: rerr.Error()}
	}
	if res.Budgets == nil {
		res.Budgets = []api.ResGetBudgetObjectBudget{}
	}
	return res, nil
}

type budget struct {
	UserId     string `db:"user_id"`
	Username   string `db:"user_name"`
	IconUrl    string `db:"icon_url"`
	BudgetId   string `db:"budget_id"`
	Name       string `db:"name"`
	Class      string `db:"class"`
	Status     string `db:"status"`
	Settlement int    `db:"settlement"`
	Budget     int    `db:"budget"`
	UpdatedAt  string `db:"updated_at"`
	Proposer   budgetObjectProposer
}

type budgetObjectProposer struct {
	UserId   string
	Username string
	IconUrl  string
}

func getBudgetList(dbClient db.Client, offset *int, proposerId *string) ([]budget, *response.Error) {
	params := struct {
		Offset     *int    `twowaysql:"offset"`
		ProposerId *string `twowaysql:"proposerId"`
	}{
		Offset:     offset,
		ProposerId: proposerId,
	}
	budgets := []budget{}
	err := dbClient.Select(&budgets, "sql/budget/select_budget.sql", &params)
	if err != nil {
		return []budget{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "稟議一覧の取得に失敗しました", Log: err.Error()}
	}
	for i := range budgets {
		budgets[i].Proposer.IconUrl = budgets[i].IconUrl
		budgets[i].Proposer.UserId = budgets[i].UserId
		budgets[i].Proposer.Username = budgets[i].Username
	}
	return budgets, nil
}
