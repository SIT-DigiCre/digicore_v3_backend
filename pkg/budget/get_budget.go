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
	budget, err := getBudgetList(dbClient, params.Offset)
	if err != nil {
		return api.ResGetBudget{}, err
	}
	rerr := copier.Copy(&res, &budget)
	if rerr != nil {
		return api.ResGetBudget{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "稟議一覧の取得に失敗しました", Log: rerr.Error()}
	}
	return res, nil
}

type budget struct {
	UserId     string `db:"user_id"`
	UserName   string `db:"user_name"`
	IconUrl    string `db:"icon_url"`
	BudgetId   string `db:"budget_id"`
	Title      string `db:"title"`
	Class      string `db:"class"`
	Status     string `db:"status"`
	Settlement string `db:"settlement"`
	Budget     string `db:"budget"`
	UpdatedAt  string `db:"updated_at"`
}

func getBudgetList(dbClient db.Client, offset *int) ([]budget, *response.Error) {
	params := struct {
		Offset *int `twowaysql:"offset"`
	}{
		Offset: offset,
	}
	budget := []budget{}
	err := dbClient.Select(&budget, "sql/budget/select_budget.sql", &params)
	if err != nil {
		return nil, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "稟議一覧の取得に失敗しました", Log: err.Error()}
	}
	if len(budget) == 0 {
		return nil, &response.Error{Code: http.StatusNotFound, Level: "Info", Message: "稟議がありません。", Log: "no rows in result"}
	}
	return budget, nil
}
