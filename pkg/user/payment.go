package user

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/util"
	"github.com/labstack/echo/v4"
)

type Payment struct {
	Year         int       `json:"year"`
	TransferName string    `json:"transfer_name"`
	Checked      bool      `json:"checked"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type RequestUpdateMyPayment struct {
	TransferName string `json:"transfer_name"`
}

func (p RequestUpdateMyPayment) validate() error {
	errorMsg := []string{}
	if 255 < utf8.RuneCountInString(p.TransferName) {
		errorMsg = append(errorMsg, "振込依頼人名は255文字未満である必要があります")
	}
	if utf8.RuneCountInString(p.TransferName) <= 0 {
		errorMsg = append(errorMsg, "振込依頼人名は1文字以上ある必要があります")
	}
	if len(errorMsg) != 0 {
		return fmt.Errorf(strings.Join(errorMsg, ","))
	}
	return nil
}

type ResponseGetMyPayment struct {
	Payment Payment `json:"payment"`
	Error   string  `json:"error"`
}

type ResponseGetMyPaymentHistory struct {
	Payments []Payment `json:"payments"`
	Error    string    `json:"error"`
}

type ResponseUpdateMyPayment struct {
	Error string `json:"error"`
}

// Get my payment
// @Router /user/my/payment [get]
// @Security Authorization
// @Success 200 {object} ResponseGetMyPayment
func (c Context) GetMyPayment(e echo.Context) error {
	userId, err := GetUserId(&e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseGetMyProfile{Error: err.Error()})
	}
	payment := Payment{}
	err = c.DB.QueryRow("SELECT year, transfer_name, checked, created_at, updated_at FROM user_payments WHERE year = ? AND user_id = UUID_TO_BIN(?)", util.NowFiscalYear(), userId).Scan(&payment.Year, &payment.TransferName, &payment.Checked, &payment.CreatedAt, &payment.UpdatedAt)
	if err == sql.ErrNoRows {
		return e.JSON(http.StatusNotFound, ResponseGetMyPayment{Error: "振り込みデータが存在しません"})
	} else if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseGetMyPayment{Error: "DBの読み込みに失敗しました"})
	}
	return e.JSON(http.StatusOK, ResponseGetMyPayment{Payment: payment})
}

// Get my payment history
// @Router /user/my/payment/history [get]
// @Security Authorization
// @Success 200 {object} ResponseGetMyPaymentHistory
func (c Context) GetMyPaymentHistory(e echo.Context) error {
	userId, err := GetUserId(&e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseGetMyProfile{Error: err.Error()})
	}
	payments := []Payment{}
	rows, err := c.DB.Query("SELECT year, transfer_name, checked, created_at, updated_at FROM user_payments WHERE user_id = UUID_TO_BIN(?)", userId)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseUpdateMyPayment{Error: "DBの読み込みに失敗しました"})
	}
	defer rows.Close()
	for rows.Next() {
		payment := Payment{}
		if err := rows.Scan(&payment.Year, &payment.TransferName, &payment.Checked, &payment.CreatedAt, &payment.UpdatedAt); err != nil {
			return e.JSON(http.StatusInternalServerError, ResponseGetMyPayment{Error: "DBの読み込みに失敗しました"})
		}
		payments = append(payments, payment)
	}
	return e.JSON(http.StatusOK, ResponseGetMyPaymentHistory{Payments: payments})
}

// Update my payment
// @Accept json
// @Param RequestUpdateMyPayment body RequestUpdateMyPayment true "my payment"
// @Router /user/my/payment [put]
// @Security Authorization
// @Success 200 {object} ResponseUpdateMyPayment
func (c Context) UpdateMyPayment(e echo.Context) error {
	userId, err := GetUserId(&e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseUpdateMyPayment{Error: err.Error()})
	}
	payment := RequestUpdateMyPayment{}
	if err := e.Bind(&payment); err != nil {
		return e.JSON(http.StatusBadRequest, ResponseUpdateMyPayment{Error: "データの読み込みに失敗しました"})
	}
	if err := payment.validate(); err != nil {
		return e.JSON(http.StatusBadRequest, ResponseUpdateMyPayment{Error: err.Error()})
	}
	_, err = c.DB.Exec(`INSERT INTO user_payments (user_id, year, transfer_name) VALUES (UUID_TO_BIN(?), ?, ?) ON DUPLICATE KEY UPDATE transfer_name = ?`,
		userId, util.NowFiscalYear(), payment.TransferName, payment.TransferName)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseUpdateMyPayment{Error: "データの登録に失敗しました"})
	}
	return e.JSON(http.StatusOK, ResponseUpdateMyPayment{})
}
