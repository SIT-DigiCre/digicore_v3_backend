package validator

import (
	"fmt"
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func Check(s interface{}) *response.Error {
	err := validate.Struct(s)
	if err != nil {
		return &response.Error{Code: http.StatusBadRequest, Level: "Info", Message: fmt.Sprintf("送信内容が制約に違反しています(%s)", err.Error()), Log: fmt.Sprintf("Validation error: %s", err.Error())}
	}
	return nil
}
