package helpers

import (
	"github.com/aliykh/reddit-feed/pkg/customErrors"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// AddTranslation - register a new translation for a specific validator identified by its tag
func AddTranslation(tag string, errMessage string) {
	engine := binding.Validator.Engine().(*validator.Validate)
	registerFn := func(ut ut.Translator) error {
		return ut.Add(tag, errMessage, false)
	}

	transFn := func(ut ut.Translator, fe validator.FieldError) string {
		param := fe.Param()
		tag := fe.Tag()

		t, err := ut.T(tag, fe.Field(), param)
		if err != nil {
			return fe.(error).Error()
		}
		return t
	}

	_ = engine.RegisterTranslation(tag, customErrors.Trans, registerFn, transFn)
}
