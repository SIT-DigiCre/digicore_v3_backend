package user

import (
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"unicode/utf8"

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

type RequestUpdateMyPrivateProfile struct {
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

var phoneNumberRegex = regexp.MustCompile(`^[0-9]{1,32}$`)

func (p RequestUpdateMyPrivateProfile) validate() error {
	errorMsg := []string{}
	if 255 < utf8.RuneCountInString(p.FirstName) {
		errorMsg = append(errorMsg, "名前は255文字未満である必要があります")
	}
	if 255 < utf8.RuneCountInString(p.LastName) {
		errorMsg = append(errorMsg, "名字は255文字未満である必要があります")
	}
	if 255 < utf8.RuneCountInString(p.FirstNameKana) {
		errorMsg = append(errorMsg, "名前のカナは255文字未満である必要があります")
	}
	if 255 < utf8.RuneCountInString(p.LastNameKana) {
		errorMsg = append(errorMsg, "名字のカナは255文字未満である必要があります")
	}
	if !phoneNumberRegex.MatchString(p.PhoneNumber) {
		errorMsg = append(errorMsg, "電話番号はハイフン無し数字のみの32文字未満である必要があります")
	}
	if 255 < utf8.RuneCountInString(p.Address) {
		errorMsg = append(errorMsg, "住所は255文字未満である必要があります")
	}
	if 255 < utf8.RuneCountInString(p.ParentName) {
		errorMsg = append(errorMsg, "保護者名は255文字未満である必要があります")
	}
	if !phoneNumberRegex.MatchString(p.ParentCellphoneNumber) {
		errorMsg = append(errorMsg, "保護者の携帯電話番号はハイフン無し数字のみの32文字未満である必要があります")
	}
	if p.ParentHomephoneNumber != "" && !phoneNumberRegex.MatchString(p.ParentHomephoneNumber) {
		errorMsg = append(errorMsg, "保護者の電話番号はハイフン無し数字のみの32文字未満である必要があります")
	}
	if 255 < utf8.RuneCountInString(p.ParentAddress) {
		errorMsg = append(errorMsg, "保護者の住所は255文字未満である必要があります")
	}
	if len(errorMsg) != 0 {
		return fmt.Errorf(strings.Join(errorMsg, ","))
	}
	return nil
}

type ResponseUpdateMyPrivateProfile struct {
	Error string `json:"error"`
}

// Get my private prodile
// @Router /user/my/private [get]
// @Security Authorization
// @Success 200 {object} ResponseGetMyPrivateProfile
func (c Context) GetMyPrivateProfile(e echo.Context) error {
	userId, err := GetUserId(&e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseGetMyPrivateProfile{Error: err.Error()})
	}
	privateProfile := PrivateProfile{}
	err = c.DB.QueryRow("SELECT first_name, last_name, first_name_kana, last_name_kana, phone_number, address, parent_name, parent_cellphone_number, parent_homephone_number, parent_address FROM user_private_profiles WHERE user_id = UUID_TO_BIN(?)", userId).
		Scan(&privateProfile.FirstName, &privateProfile.LastName, &privateProfile.FirstNameKana, &privateProfile.LastNameKana, &privateProfile.PhoneNumber, &privateProfile.Address, &privateProfile.ParentName, &privateProfile.ParentCellphoneNumber, &privateProfile.ParentHomephoneNumber, &privateProfile.ParentAddress)
	if err == sql.ErrNoRows {
		return e.JSON(http.StatusNotFound, ResponseGetMyPrivateProfile{Error: "データが登録されていません"})
	} else if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseGetMyPrivateProfile{Error: "取得に失敗しました"})
	}
	return e.JSON(http.StatusOK, ResponseGetMyPrivateProfile{PrivateProfile: privateProfile})
}

// Update my private profile
// @Accept json
// @Param RequestUpdateMyPrivateProfile body RequestUpdateMyPrivateProfile true "my private profile"
// @Security Authorization
// @Router /user/my/private [put]
// @Success 200 {object} ResponseUpdateMyPrivateProfile
func (c Context) UpdateMyPrivateProfile(e echo.Context) error {
	userId, err := GetUserId(&e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseUpdateMyPrivateProfile{Error: err.Error()})
	}
	privateProfile := RequestUpdateMyPrivateProfile{}
	if err := e.Bind(&privateProfile); err != nil {
		return e.JSON(http.StatusBadRequest, ResponseUpdateMyPrivateProfile{Error: "データの読み込みに失敗しました"})
	}
	if err := privateProfile.validate(); err != nil {
		return e.JSON(http.StatusBadRequest, ResponseUpdateMyPrivateProfile{Error: err.Error()})
	}
	_, err = c.DB.Exec(`INSERT INTO user_private_profiles (user_id, first_name, last_name, first_name_kana, last_name_kana, phone_number, address, parent_name, parent_cellphone_number, parent_homephone_number, parent_address) VALUES (UUID_TO_BIN(?), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
				ON DUPLICATE KEY UPDATE first_name = ?, last_name = ?, first_name_kana = ?, last_name_kana = ?, phone_number = ?, address = ?, parent_name = ?, parent_cellphone_number = ?, parent_homephone_number = ?, parent_address = ?`,
		userId, privateProfile.FirstName, privateProfile.LastName, privateProfile.FirstNameKana, privateProfile.LastNameKana, privateProfile.PhoneNumber, privateProfile.Address, privateProfile.ParentName, privateProfile.ParentCellphoneNumber, privateProfile.ParentHomephoneNumber, privateProfile.ParentAddress,
		privateProfile.FirstName, privateProfile.LastName, privateProfile.FirstNameKana, privateProfile.LastNameKana, privateProfile.PhoneNumber, privateProfile.Address, privateProfile.ParentName, privateProfile.ParentCellphoneNumber, privateProfile.ParentHomephoneNumber, privateProfile.ParentAddress)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseUpdateMyPrivateProfile{Error: "更新に失敗しました"})
	}
	return e.JSON(http.StatusOK, ResponseUpdateMyPrivateProfile{})
}
