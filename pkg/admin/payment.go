package admin

import (
	"net/http"
	"strconv"
	"time"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/util"
	"github.com/labstack/echo/v4"
)

type Payment struct {
	Id            string    `json:"id"`
	StudentNumber string    `json:"student_number"`
	Year          int       `json:"year"`
	TransferName  string    `json:"transfer_name"`
	Checked       bool      `json:"checked"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ResponseGetAllPayments struct {
	Payments []Payment `json:"payments"`
	Error    string    `json:"error"`
}

// Get all payments
// @Router /admin/payments [get]
// @Param year query int false "year"
// @Security Authorization
// @Success 200 {object} ResponseGetAllPayments
func (c Context) GetAllPayments(e echo.Context) error {
	year, err := strconv.Atoi(e.QueryParam("year"))
	if err != nil {
		year = util.NowFiscalYear()
	}
	payments := []Payment{}
	rows, err := c.DB.Query("SELECT BIN_TO_UUID(UserPayment.id), User.student_number, transfer_name, year, checked, created_at, updated_at FROM UserPayment LEFT JOIN User ON UserPayment.user_id = User.id WHERE year = ?", year)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseGetAllPayments{Error: "DBの読み込みに失敗しました"})
	}
	defer rows.Close()
	for rows.Next() {
		payment := Payment{}
		if err := rows.Scan(&payment.Id, &payment.StudentNumber, &payment.TransferName, &payment.Year, &payment.Checked, &payment.CreatedAt, &payment.UpdatedAt); err != nil {
			return e.JSON(http.StatusBadRequest, ResponseGetAllPayments{Error: "DBの読み込みに失敗しました"})
		}
		payments = append(payments, payment)
	}
	return e.JSON(http.StatusOK, ResponseGetAllPayments{Payments: payments})
}
