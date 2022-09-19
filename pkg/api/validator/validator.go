package validator

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/go-playground/locales/ja_JP"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	ja_translations "github.com/go-playground/validator/v10/translations/ja"
	"github.com/sirupsen/logrus"
)

var (
	trans    ut.Translator
	validate *validator.Validate
)

func init() {
	japanese := ja_JP.New()
	uni := ut.New(japanese, japanese)
	trans, _ = uni.GetTranslator("ja")

	validate = validator.New()

	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		fieldName := field.Tag.Get("ja")
		if fieldName == "-" {
			return ""
		}
		return fieldName
	})

	err := ja_translations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		logrus.Fatal("Failed to create validation translator")
	}
}

func Validate(s interface{}) *response.Error {
	err := validate.Struct(s)

	if err == nil {
		return nil
	}
	errs := err.(validator.ValidationErrors)
	message := ""
	for _, ve := range errs.Translate(trans) {
		if message != "" {
			message += ", "
		}
		message += ve
	}
	return &response.Error{Code: http.StatusBadRequest, Level: "Info", Message: fmt.Sprintf("送信内容が制約に違反しています(%s)", message), Log: fmt.Sprintf("Validation error: %s", message)}
}
