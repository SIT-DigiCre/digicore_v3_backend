package user

import (
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
	return e.JSON(http.StatusOK, ResponseGetMyPrivateProfile{})
}

// Set my private prodile
// @Accept json
// @Param PrivateProfile body PrivateProfile true "my private profile"
// @Router /user/my/private [post]
// @Success 200 {object} ResponseSetMyPrivateProfile
func (c Context) SetMyPrivateProfile(e echo.Context) error {
	return e.JSON(http.StatusOK, ResponseSetMyPrivateProfile{})
}
