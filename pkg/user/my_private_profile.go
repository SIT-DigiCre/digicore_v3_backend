package user

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type PrivateProfile struct {
	FirstName             string `json:"first_name"`
	LastName              string `json:"last_name"`
	FirstNameKana         string `json:"first_name_kana"`
	LastNameKana          string `json:"last_name_kana"`
	PhoneNumber           string `json:"phone_number"`
	Address               string `json:"address"`
	ParentName            string `json:"parent_name"`
	ParentCellphoneNumber string `json:"parent_cellphone_number"`
	ParentHomephoneNumber string `json:"parent_homephone_number"`
	ParentAddress         string `json:"parent_address"`
}

func (p PrivateProfile) validate() error {
	return nil
}

type ResponseGetMyPrivateProfile struct {
	PrivateProfile PrivateProfile `json:"private_profile"`
	Error          string         `json:"error"`
}

type ResponseSetMyPrivateProfile struct {
	Error string `json:"error"`
}

// Get my private prodile
// @Accept json
// @Router /user/my/private [get]
// @Success 200 {object} ResponseGetMyPrivateProfile
func (c Context) GetMyPrivateProfile(e echo.Context) error {
	userId, err := GetUserId(&e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseGetMyPrivateProfile{Error: ""})
	}
	privateProfile := PrivateProfile{}
	err = c.DB.QueryRow("SELECT first_name, last_name, first_name_kana, last_name_kana, phone_number, address, parent_name, parent_cellphone_number, parent_homephone_number, parent_address FROM UserPrivateProfile WHERE user_id = UUID_TO_BIN(?)", userId).
		Scan(&privateProfile.FirstName, &privateProfile.LastName, &privateProfile.FirstNameKana, &privateProfile.LastNameKana, &privateProfile.PhoneNumber, &privateProfile.Address, &privateProfile.ParentName, &privateProfile.ParentCellphoneNumber, &privateProfile.ParentHomephoneNumber, &privateProfile.ParentAddress)
	if err == sql.ErrNoRows {
		return e.JSON(http.StatusBadRequest, ResponseGetMyPrivateProfile{Error: err.Error()})
	} else if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseGetMyPrivateProfile{Error: err.Error()})
	}
	return e.JSON(http.StatusOK, ResponseGetMyPrivateProfile{PrivateProfile: privateProfile})
}

// Set my private prodile
// @Accept json
// @Param PrivateProfile body PrivateProfile true "my private profile"
// @Router /user/my/private [post]
// @Success 200 {object} ResponseSetMyPrivateProfile
func (c Context) SetMyPrivateProfile(e echo.Context) error {
	userId, err := GetUserId(&e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseGetMyPrivateProfile{})
	}
	fmt.Println(userId)
	privateProfile := PrivateProfile{}
	if err := e.Bind(&privateProfile); err != nil {
		return e.JSON(http.StatusBadRequest, ResponseSetMyPrivateProfile{Error: err.Error()})
	}
	if err := privateProfile.validate(); err != nil {
		return e.JSON(http.StatusBadRequest, ResponseSetMyPrivateProfile{Error: err.Error()})
	}
	_, err = c.DB.Exec(`INSERT INTO UserPrivateProfile (user_id, first_name, last_name, first_name_kana, last_name_kana, phone_number, address, parent_name, parent_cellphone_number, parent_homephone_number, parent_address) VALUES (UUID_TO_BIN(?), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
				ON DUPLICATE KEY UPDATE first_name = ?, last_name = ?, first_name_kana = ?, last_name_kana = ?, phone_number = ?, address = ?, parent_name = ?, parent_cellphone_number = ?, parent_homephone_number = ?, parent_address = ?`,
		userId, privateProfile.FirstName, privateProfile.LastName, privateProfile.FirstNameKana, privateProfile.LastNameKana, privateProfile.PhoneNumber, privateProfile.Address, privateProfile.ParentName, privateProfile.ParentCellphoneNumber, privateProfile.ParentHomephoneNumber, privateProfile.ParentAddress,
		privateProfile.FirstName, privateProfile.LastName, privateProfile.FirstNameKana, privateProfile.LastNameKana, privateProfile.PhoneNumber, privateProfile.Address, privateProfile.ParentName, privateProfile.ParentCellphoneNumber, privateProfile.ParentHomephoneNumber, privateProfile.ParentAddress)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseSetMyPrivateProfile{Error: err.Error()})
	}
	return e.JSON(http.StatusOK, ResponseSetMyPrivateProfile{})
}
