package validator

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"

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

	err := validate.RegisterValidation("phonenumber", func(fl validator.FieldLevel) bool {
		fmt.Print(fl.Field().String())
		return regexp.MustCompile(`^\d{0,15}$`).MatchString(fl.Field().String())
	})
	if err != nil {
		logrus.Fatal("Failed to create validation translator")
	}

	err = ja_translations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		logrus.Fatal("Failed to create validation translator")
	}
}

func Validate(s interface{}) *response.Error {
	err := validate.Struct(s)

	if err == nil {
		return nil
	}
	message := ""
	for _, err := range err.(validator.ValidationErrors) {
		mes := err.Translate(trans)
		if message != "" {
			message += ", "
		}
		if err.ActualTag() == "phonenumber" {
			message += fmt.Sprintf("%sは電話番号でなければなりません", err.Field())
		} else {
			message += mes
		}
	}
	return &response.Error{Code: http.StatusBadRequest, Level: "Info", Message: fmt.Sprintf("送信内容が制約に違反しています(%s)", message), Log: fmt.Sprintf("Validation error: %s", message)}
}
