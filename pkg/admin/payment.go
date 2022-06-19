package admin

import (
	"database/sql"
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

type ResponseGetPayment struct {
	Payment Payment `json:"payment"`
	Error   string  `json:"error"`
}

type ResponseUpdatePayment struct {
	Error string `json:"error"`
}

type RequestUpdatePayment struct {
	Checked bool `json:"checked"`
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
	rows, err := c.DB.Query("SELECT BIN_TO_UUID(user_payments.id), users.student_number, transfer_name, year, checked, created_at, updated_at FROM user_payments LEFT JOIN users ON user_payments.user_id = users.id WHERE year = ?", year)
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

// Get payment
// @Router /admin/payments/{id} [get]
// @Param id path string true "payment id"
// @Security Authorization
// @Success 200 {object} ResponseGetPayment
func (c Context) GetPayment(e echo.Context) error {
	id := e.Param("id")
	if id == "" {
		return e.JSON(http.StatusBadRequest, ResponseGetPayment{Error: "データの読み込みに失敗しました"})
	}
	payment := Payment{}
	err := c.DB.QueryRow("SELECT BIN_TO_UUID(user_payments.id), users.student_number, transfer_name, year, checked, created_at, updated_at FROM user_payments LEFT JOIN users ON user_payments.user_id = users.id WHERE user_payments.id = UUID_TO_BIN(?)", id).Scan(&payment.Id, &payment.StudentNumber, &payment.TransferName, &payment.Year, &payment.Checked, &payment.CreatedAt, &payment.UpdatedAt)
	if err == sql.ErrNoRows {
		return e.JSON(http.StatusNotFound, ResponseGetPayment{Error: "データが見つかりませんでした"})
	} else if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseGetPayment{Error: "DBの読み込みに失敗しました"})
	}
	return e.JSON(http.StatusOK, ResponseGetPayment{Payment: payment})
}

// Update payment
// @Router /admin/payments/{id} [put]
// @Param id path string true "payment id"
// @Param RequestUpdatePayment body RequestUpdatePayment true "payment data"
// @Security Authorization
// @Success 200 {object} ResponseUpdatePayment
func (c Context) UpdatePayment(e echo.Context) error {
	id := e.Param("id")
	if id == "" {
		return e.JSON(http.StatusBadRequest, ResponseGetPayment{Error: "データの読み込みに失敗しました"})
	}
	payment := RequestUpdatePayment{}
	if err := e.Bind(&payment); err != nil {
		return e.JSON(http.StatusBadRequest, ResponseUpdatePayment{Error: "データの読み込みに失敗しました"})
	}
	_, err := c.DB.Exec("UPDATE user_payments SET checked = ? WHERE id = UUID_TO_BIN(?)", payment.Checked, id)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseUpdatePayment{Error: "DBの読み込みに失敗しました"})
	}
	return e.JSON(http.StatusOK, ResponseUpdatePayment{})
}
