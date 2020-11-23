package validate

import (
	"reflect"
	"strings"

	zhcn "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhcn_translations "github.com/go-playground/validator/v10/translations/zh"
)

var Validate *validator.Validate
var trans ut.Translator

func init() {
	zh := zhcn.New()
	uni := ut.New(zh, zh)
	trans, _ = uni.GetTranslator("zh")

	Validate = validator.New()
	Validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		label := field.Tag.Get("label")
		if label == "" {
			return field.Name
		}
		return label
	})
	_ = zhcn_translations.RegisterDefaultTranslations(Validate, trans)
}
func Translate(errs validator.ValidationErrors) string {
	var errList []string
	for _, e := range errs {
		// can translate each error one at a time.
		errList = append(errList, e.Translate(trans))
	}
	if len(errList) > 0 {
		return errList[0]
	} else {
		return strings.Join(errList, "|")
	}
}
