package user

import (
	"database/sql"
	"net/http"
	"time"

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
	err = c.DB.QueryRow("SELECT year, transfer_name, checked, created_at, updated_at FROM UserPayment WHERE year = ? AND user_id = UUID_TO_BIN(?)", util.NowFiscalYear(), userId).Scan(&payment.Year, &payment.TransferName, &payment.Checked, &payment.CreatedAt, &payment.UpdatedAt)
	if err == sql.ErrNoRows {
		return e.JSON(http.StatusOK, ResponseGetMyPayment{Error: err.Error()})
	} else if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseGetMyPayment{Error: err.Error()})
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
	rows, err := c.DB.Query("SELECT year, transfer_name, checked, created_at, updated_at FROM UserPayment WHERE user_id = UUID_TO_BIN(?)", userId)
	defer rows.Close()
	for rows.Next() {
		payment := Payment{}
		if err := rows.Scan(&payment.Year, &payment.TransferName, &payment.Checked, &payment.CreatedAt, &payment.UpdatedAt); err != nil {
			return e.JSON(http.StatusBadRequest, ResponseGetMyPayment{Error: err.Error()})
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
		return e.JSON(http.StatusBadRequest, ResponseUpdateMyPayment{Error: err.Error()})
	}
	if err := payment.validate(); err != nil {
		return e.JSON(http.StatusBadRequest, ResponseUpdateMyPayment{Error: err.Error()})
	}
	_, err = c.DB.Exec(`INSERT INTO UserPayment (user_id, year, transfer_name) VALUES (UUID_TO_BIN(?), ?, ?) ON DUPLICATE KEY UPDATE transfer_name = ?`,
		userId, util.NowFiscalYear(), payment.TransferName, payment.TransferName)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseUpdateMyPayment{Error: err.Error()})
	}
	return e.JSON(http.StatusOK, ResponseUpdateMyPayment{})
}
