package helpers

import (
	"reflect"
	"strings"
	"sync"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Custom Validator take from: https://seb-nyberg.medium.com/better-validation-errors-in-go-gin-88f983564a3d

type DefaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

var _ binding.StructValidator = &DefaultValidator{}

func (v *DefaultValidator) ValidateStruct(obj interface{}) error {
	if kindOfData(obj) == reflect.Struct {
		v.lazyInit()

		if err := v.validate.Struct(obj); err != nil {
			return error(err)
		}
	}

	return nil
}

func (v *DefaultValidator) Engine() interface{} {
	v.lazyInit()
	return v.validate
}

func (v *DefaultValidator) lazyInit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName("binding")
		// Print JSON name on validator.FieldError.Field()
		v.validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			return name
		})
	})
}

func kindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}
